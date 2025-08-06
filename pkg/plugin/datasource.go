package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/nextmap-io/grafana-updownio-datasource/pkg/models"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ backend.CallResourceHandler   = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	config, err := models.LoadPluginSettings(settings)
	if err != nil {
		return nil, fmt.Errorf("unable to load settings: %w", err)
	}

	ds := &Datasource{
		settings: config,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	// Set up resource handler
	mux := http.NewServeMux()
	mux.HandleFunc("/checks", ds.handleGetChecks)
	mux.HandleFunc("/health", ds.handleHealth)
	ds.resourceHandler = httpadapter.New(mux)

	return ds, nil
}

// Datasource is an UpDown.io datasource which can respond to data queries, reports
// its health and provides API resources.
type Datasource struct {
	settings        *models.PluginSettings
	httpClient      *http.Client
	resourceHandler backend.CallResourceHandler
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *Datasource) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	// Unmarshal the JSON into our queryModel.
	var qm models.QueryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	switch qm.QueryType {
	case "checks":
		return d.queryChecks(ctx, query)
	case "metrics":
		return d.queryMetrics(ctx, query, qm)
	case "downtimes":
		return d.queryDowntimes(ctx, query, qm)
	case "uptime":
		return d.queryUptime(ctx, query, qm)
	default:
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unknown query type: %s", qm.QueryType))
	}
}

func (d *Datasource) queryChecks(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	checks, err := d.fetchChecks(ctx)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to fetch checks: %v", err))
	}

	// Create a frame with check information
	frame := data.NewFrame("checks")

	// Add fields
	tokens := make([]string, len(checks))
	urls := make([]string, len(checks))
	aliases := make([]string, len(checks))
	uptimes := make([]float64, len(checks))
	statuses := make([]bool, len(checks))

	for i, check := range checks {
		tokens[i] = check.Token
		urls[i] = check.URL
		aliases[i] = check.Alias
		uptimes[i] = check.Uptime
		statuses[i] = !check.Down
	}

	frame.Fields = append(frame.Fields,
		data.NewField("token", nil, tokens),
		data.NewField("url", nil, urls),
		data.NewField("alias", nil, aliases),
		data.NewField("uptime", nil, uptimes),
		data.NewField("status", nil, statuses),
	)

	return backend.DataResponse{Frames: []*data.Frame{frame}}
}

func (d *Datasource) queryMetrics(ctx context.Context, query backend.DataQuery, qm models.QueryModel) backend.DataResponse {
	if qm.CheckToken == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "check token is required for metrics query")
	}

	metrics, err := d.fetchMetrics(ctx, qm.CheckToken, query.TimeRange.From, query.TimeRange.To, qm.GroupBy)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to fetch metrics: %v", err))
	}

	// Create a frame with metrics
	frame := data.NewFrame("metrics")

	// For simplicity, we return uptime and apdex as scalar values
	// In a more advanced implementation, we could parse time-grouped data
	if metrics.Uptime != nil {
		frame.Fields = append(frame.Fields,
			data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
			data.NewField("uptime", nil, []float64{*metrics.Uptime, *metrics.Uptime}),
		)
	}

	if metrics.Apdex != nil {
		if len(frame.Fields) == 0 {
			frame.Fields = append(frame.Fields,
				data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
			)
		}
		frame.Fields = append(frame.Fields,
			data.NewField("apdex", nil, []float64{*metrics.Apdex, *metrics.Apdex}),
		)
	}

	return backend.DataResponse{Frames: []*data.Frame{frame}}
}

func (d *Datasource) queryDowntimes(ctx context.Context, query backend.DataQuery, qm models.QueryModel) backend.DataResponse {
	if qm.CheckToken == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "check token is required for downtimes query")
	}

	downtimes, err := d.fetchDowntimes(ctx, qm.CheckToken)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to fetch downtimes: %v", err))
	}

	// Create a frame with downtimes
	frame := data.NewFrame("downtimes")

	if len(downtimes) > 0 {
		ids := make([]string, len(downtimes))
		errors := make([]string, len(downtimes))
		startedAts := make([]time.Time, len(downtimes))
		durations := make([]*int, len(downtimes))

		for i, downtime := range downtimes {
			ids[i] = downtime.ID
			errors[i] = downtime.Error
			if startedAt, err := time.Parse(time.RFC3339, downtime.StartedAt); err == nil {
				startedAts[i] = startedAt
			}
			durations[i] = downtime.Duration
		}

		frame.Fields = append(frame.Fields,
			data.NewField("id", nil, ids),
			data.NewField("error", nil, errors),
			data.NewField("started_at", nil, startedAts),
			data.NewField("duration", nil, durations),
		)
	}

	return backend.DataResponse{Frames: []*data.Frame{frame}}
}

func (d *Datasource) queryUptime(ctx context.Context, query backend.DataQuery, qm models.QueryModel) backend.DataResponse {
	if qm.CheckToken == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "check token is required for uptime query")
	}

	check, err := d.fetchCheck(ctx, qm.CheckToken)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to fetch check: %v", err))
	}

	// Create a frame with uptime
	frame := data.NewFrame("uptime")
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
		data.NewField("uptime", nil, []float64{check.Uptime, check.Uptime}),
	)

	return backend.DataResponse{Frames: []*data.Frame{frame}}
}

// CallResource HTTP style resource call handler
func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return d.resourceHandler.CallResource(ctx, req, sender)
}

func (d *Datasource) handleGetChecks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checks, err := d.fetchChecks(r.Context())
	if err != nil {
		log.DefaultLogger.Error("Failed to fetch checks", "error", err)
		http.Error(w, "Failed to fetch checks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checks)
}

func (d *Datasource) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// API methods
func (d *Datasource) fetchChecks(ctx context.Context) ([]models.UpdownCheck, error) {
	url := fmt.Sprintf("%s/checks", d.settings.ApiUrl)
	var checks []models.UpdownCheck
	
	err := d.makeAPIRequest(ctx, "GET", url, nil, &checks)
	return checks, err
}

func (d *Datasource) fetchCheck(ctx context.Context, token string) (*models.UpdownCheck, error) {
	url := fmt.Sprintf("%s/checks/%s", d.settings.ApiUrl, token)
	var check models.UpdownCheck
	
	err := d.makeAPIRequest(ctx, "GET", url, nil, &check)
	return &check, err
}

func (d *Datasource) fetchMetrics(ctx context.Context, token string, from, to time.Time, groupBy string) (*models.UpdownMetrics, error) {
	url := fmt.Sprintf("%s/checks/%s/metrics", d.settings.ApiUrl, token)
	
	// Ajouter les paramètres de temps et groupement
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	q := req.URL.Query()
	q.Add("from", from.Format(time.RFC3339))
	q.Add("to", to.Format(time.RFC3339))
	if groupBy != "" {
		q.Add("group", groupBy)
	}
	req.URL.RawQuery = q.Encode()
	
	var metrics models.UpdownMetrics
	err = d.makeAPIRequestWithRequest(req, &metrics)
	return &metrics, err
}

func (d *Datasource) fetchDowntimes(ctx context.Context, token string) ([]models.UpdownDowntime, error) {
	url := fmt.Sprintf("%s/checks/%s/downtimes", d.settings.ApiUrl, token)
	var downtimes []models.UpdownDowntime
	
	err := d.makeAPIRequest(ctx, "GET", url, nil, &downtimes)
	return downtimes, err
}

func (d *Datasource) makeAPIRequest(ctx context.Context, method, url string, body io.Reader, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	
	return d.makeAPIRequestWithRequest(req, result)
}

func (d *Datasource) makeAPIRequestWithRequest(req *http.Request, result interface{}) error {
	// Ajouter la clé API en header
	req.Header.Set("X-API-KEY", d.settings.Secrets.ApiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Grafana-UpdownIO-Datasource/1.0")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if d.settings.Secrets.ApiKey == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "API key missing",
		}, nil
	}

	// Test the API by fetching checks
	_, err := d.fetchChecks(ctx)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Unable to connect to UpDown.io API: %v", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "UpDown.io API connection successful",
	}, nil
}
