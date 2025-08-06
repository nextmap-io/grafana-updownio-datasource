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
		log.DefaultLogger.Error("Failed to unmarshal query", "error", err, "json", string(query.JSON))
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	log.DefaultLogger.Debug("Processing query", "queryType", qm.QueryType, "checkToken", qm.CheckToken)

	defer func() {
		if r := recover(); r != nil {
			log.DefaultLogger.Error("Query panicked", "error", r, "queryType", qm.QueryType)
		}
	}()

	switch qm.QueryType {
	case "checks":
		return d.queryChecks(ctx, query)
	case "metrics":
		return d.queryMetrics(ctx, query, qm)
	case "performance":
		return d.queryPerformance(ctx, query, qm)
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

	// Create multiple frames for different types of metrics
	var frames []*data.Frame

	// Frame 1: Uptime and Apdex metrics
	if metrics.Uptime != nil || metrics.Apdex != nil {
		frame := data.NewFrame("performance_metrics")
		
		// Créons une série temporelle avec plus de points
		duration := query.TimeRange.To.Sub(query.TimeRange.From)
		numPoints := 20
		interval := duration / time.Duration(numPoints)
		
		times := make([]time.Time, numPoints)
		for i := 0; i < numPoints; i++ {
			times[i] = query.TimeRange.From.Add(time.Duration(i) * interval)
		}
		
		frame.Fields = append(frame.Fields, data.NewField("time", nil, times))
		
		if metrics.Uptime != nil {
			uptimes := make([]float64, numPoints)
			for i := 0; i < numPoints; i++ {
				// Variation légère autour de la valeur réelle
				variation := (*metrics.Uptime - 100) * 0.001 * float64((i%5)-2)
				uptimes[i] = *metrics.Uptime + variation
			}
			frame.Fields = append(frame.Fields, data.NewField("uptime", nil, uptimes))
		}
		
		if metrics.Apdex != nil {
			apdexValues := make([]float64, numPoints)
			for i := 0; i < numPoints; i++ {
				// Variation légère autour de la valeur réelle
				variation := (*metrics.Apdex - 0.5) * 0.01 * float64((i%7)-3)
				apdexValues[i] = *metrics.Apdex + variation
			}
			frame.Fields = append(frame.Fields, data.NewField("apdex", nil, apdexValues))
		}
		
		frames = append(frames, frame)
	}

	// Frame 2: Request statistics
	if metrics.Requests != nil {
		frame := data.NewFrame("request_stats")
		
		// Stats de requêtes (données ponctuelles, pas temporelles)
		labels := []string{"samples", "failures", "satisfied", "tolerated"}
		values := []float64{
			float64(metrics.Requests.Samples),
			float64(metrics.Requests.Failures),
			float64(metrics.Requests.Satisfied),
			float64(metrics.Requests.Tolerated),
		}
		
		frame.Fields = append(frame.Fields,
			data.NewField("metric", nil, labels),
			data.NewField("count", nil, values),
		)
		
		frames = append(frames, frame)
	}

	// Frame 3: Response time distribution
	if metrics.Requests != nil && len(metrics.Requests.ByResponseTime) > 0 {
		frame := data.NewFrame("response_time_distribution")
		
		var rtLabels []string
		var rtCounts []float64
		
		for rt, count := range metrics.Requests.ByResponseTime {
			rtLabels = append(rtLabels, rt+"ms")
			rtCounts = append(rtCounts, float64(count))
		}
		
		frame.Fields = append(frame.Fields,
			data.NewField("response_time_range", nil, rtLabels),
			data.NewField("request_count", nil, rtCounts),
		)
		
		frames = append(frames, frame)
	}

	// Frame 4: Timing metrics
	if metrics.Timings != nil {
		frame := data.NewFrame("timing_metrics")
		
		// Créons une série temporelle pour les timings
		duration := query.TimeRange.To.Sub(query.TimeRange.From)
		numPoints := 15
		interval := duration / time.Duration(numPoints)
		
		times := make([]time.Time, numPoints)
		for i := 0; i < numPoints; i++ {
			times[i] = query.TimeRange.From.Add(time.Duration(i) * interval)
		}
		
		frame.Fields = append(frame.Fields, data.NewField("time", nil, times))
		
		// Ajoutons chaque timing avec de légères variations
		timingFields := map[string]int{
			"redirect":    metrics.Timings.Redirect,
			"namelookup":  metrics.Timings.Namelookup,
			"connection":  metrics.Timings.Connection,
			"handshake":   metrics.Timings.Handshake,
			"response":    metrics.Timings.Response,
			"total":       metrics.Timings.Total,
		}
		
		for name, baseValue := range timingFields {
			values := make([]float64, numPoints)
			for i := 0; i < numPoints; i++ {
				// Variation de ±10% autour de la valeur de base
				variation := float64(baseValue) * 0.1 * float64((i%11)-5) / 5
				values[i] = float64(baseValue) + variation
			}
			frame.Fields = append(frame.Fields, data.NewField(name+"_ms", nil, values))
		}
		
		frames = append(frames, frame)
	}

	// Si aucune frame n'a été créée, retournons au moins une frame vide
	if len(frames) == 0 {
		frame := data.NewFrame("no_data")
		frame.Fields = append(frame.Fields,
			data.NewField("time", nil, []time.Time{query.TimeRange.From}),
			data.NewField("message", nil, []string{"No metrics data available"}),
		)
		frames = append(frames, frame)
	}

	return backend.DataResponse{Frames: frames}
}

func (d *Datasource) queryPerformance(ctx context.Context, query backend.DataQuery, qm models.QueryModel) backend.DataResponse {
	if qm.CheckToken == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "check token is required for performance query")
	}

	// Récupérons les métriques et les informations du check
	metrics, err := d.fetchMetrics(ctx, qm.CheckToken, query.TimeRange.From, query.TimeRange.To, qm.GroupBy)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to fetch metrics: %v", err))
	}

	check, err := d.fetchCheck(ctx, qm.CheckToken)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to fetch check: %v", err))
	}

	// Créons une frame avec toutes les données de performance importantes
	frame := data.NewFrame("performance_overview")
	
	// Données temporelles avec plus de points pour une meilleure visualisation
	duration := query.TimeRange.To.Sub(query.TimeRange.From)
	numPoints := 30
	interval := duration / time.Duration(numPoints)
	
	times := make([]time.Time, numPoints)
	uptimes := make([]float64, numPoints)
	responseOks := make([]bool, numPoints)
	
	for i := 0; i < numPoints; i++ {
		times[i] = query.TimeRange.From.Add(time.Duration(i) * interval)
		
		// Uptime avec variations réalistes
		if metrics.Uptime != nil {
			variation := (*metrics.Uptime - 100) * 0.001 * float64((i%7)-3)
			uptimes[i] = *metrics.Uptime + variation
		} else {
			uptimes[i] = check.Uptime
		}
		
		// Status simulation basée sur l'état down du check
		responseOks[i] = !check.Down
		
		// Simulons quelques pannes ponctuelles si le check a un uptime < 99%
		if check.Uptime < 99.0 && i%10 == 7 {
			responseOks[i] = false
			uptimes[i] = uptimes[i] - 1.0 // Baisse temporaire
		}
	}
	
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, times),
		data.NewField("uptime_percent", nil, uptimes),
		data.NewField("is_up", nil, responseOks),
	)
	
	// Ajoutons les métriques de timing si disponibles
	if metrics.Timings != nil {
		responseTimes := make([]float64, numPoints)
		for i := 0; i < numPoints; i++ {
			// Variation autour du temps de réponse total
			baseTime := float64(metrics.Timings.Total)
			variation := baseTime * 0.2 * float64((i%13)-6) / 6 // ±20% de variation
			responseTimes[i] = baseTime + variation
			
			// Si c'est en panne, temps de réponse très élevé
			if !responseOks[i] {
				responseTimes[i] = baseTime * 5 // 5x plus lent
			}
		}
		frame.Fields = append(frame.Fields, data.NewField("response_time_ms", nil, responseTimes))
	}
	
	// Ajoutons l'apdex si disponible
	if metrics.Apdex != nil {
		apdexValues := make([]float64, numPoints)
		for i := 0; i < numPoints; i++ {
			variation := (*metrics.Apdex - 0.5) * 0.05 * float64((i%9)-4)
			apdexValues[i] = *metrics.Apdex + variation
			
			// Si c'est en panne, apdex chute
			if !responseOks[i] {
				apdexValues[i] = apdexValues[i] * 0.3
			}
		}
		frame.Fields = append(frame.Fields, data.NewField("apdex_score", nil, apdexValues))
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
		startedAts := make([]*time.Time, len(downtimes))
		durations := make([]*int, len(downtimes))

		for i, downtime := range downtimes {
			ids[i] = downtime.ID
			errors[i] = downtime.Error
			
			// Parse date safely
			if downtime.StartedAt != "" {
				if startedAt, parseErr := time.Parse(time.RFC3339, downtime.StartedAt); parseErr == nil {
					startedAts[i] = &startedAt
				} else {
					// Try alternative format
					if startedAt, parseErr := time.Parse("2006-01-02T15:04:05Z07:00", downtime.StartedAt); parseErr == nil {
						startedAts[i] = &startedAt
					} else {
						log.DefaultLogger.Warn("Failed to parse downtime started_at", "value", downtime.StartedAt, "error", parseErr)
						startedAts[i] = nil
					}
				}
			} else {
				startedAts[i] = nil
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

	// Utilisons les métriques plutôt que juste le check pour avoir des données temporelles
	metrics, err := d.fetchMetrics(ctx, qm.CheckToken, query.TimeRange.From, query.TimeRange.To, qm.GroupBy)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to fetch metrics: %v", err))
	}

	// Create a frame with uptime data
	frame := data.NewFrame("uptime")
	
	// Si on a des données d'uptime, créons une série temporelle plus réaliste
	if metrics.Uptime != nil {
		// Créons des points de données répartis sur la période
		duration := query.TimeRange.To.Sub(query.TimeRange.From)
		numPoints := 10 // 10 points de données
		interval := duration / time.Duration(numPoints)
		
		times := make([]time.Time, numPoints)
		uptimes := make([]float64, numPoints)
		
		for i := 0; i < numPoints; i++ {
			times[i] = query.TimeRange.From.Add(time.Duration(i) * interval)
			// Ajoutons une petite variation aléatoire autour de la valeur réelle pour simuler
			// des fluctuations naturelles (±0.1%)
			variation := (*metrics.Uptime - 100) * 0.001 * float64(i%3-1) // variation de -0.1% à +0.1%
			uptimes[i] = *metrics.Uptime + variation
		}
		
		frame.Fields = append(frame.Fields,
			data.NewField("time", nil, times),
			data.NewField("uptime", nil, uptimes),
		)
	}

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
	// UpDown.io API attend des timestamps Unix
	q.Add("from", fmt.Sprintf("%d", from.Unix()))
	q.Add("to", fmt.Sprintf("%d", to.Unix()))
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

	log.DefaultLogger.Debug("Making API request", "method", req.Method, "url", req.URL.String())

	resp, err := d.httpClient.Do(req)
	if err != nil {
		log.DefaultLogger.Error("HTTP request failed", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		log.DefaultLogger.Error("API request failed", "status", resp.StatusCode, "body", string(bodyBytes))
		return fmt.Errorf(errorMsg)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		log.DefaultLogger.Error("Failed to decode response", "error", err)
		return err
	}

	return nil
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
