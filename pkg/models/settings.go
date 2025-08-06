package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSettings struct {
	ApiUrl  string                `json:"apiUrl"`
	Secrets *SecretPluginSettings `json:"-"`
}

type SecretPluginSettings struct {
	ApiKey string `json:"apiKey"`
}

// UpdownCheck représente un check UpDown.io
type UpdownCheck struct {
	Token               string            `json:"token"`
	URL                 string            `json:"url"`
	Type                string            `json:"type"`
	Alias               string            `json:"alias"`
	LastStatus          int               `json:"last_status"`
	Uptime              float64           `json:"uptime"`
	Down                bool              `json:"down"`
	DownSince           *string           `json:"down_since"`
	UpSince             *string           `json:"up_since"`
	Error               *string           `json:"error"`
	Period              int               `json:"period"`
	ApdexT              float64           `json:"apdex_t"`
	StringMatch         string            `json:"string_match"`
	Enabled             bool              `json:"enabled"`
	Published           bool              `json:"published"`
	DisabledLocations   []string          `json:"disabled_locations"`
	Recipients          []string          `json:"recipients"`
	LastCheckAt         string            `json:"last_check_at"`
	NextCheckAt         string            `json:"next_check_at"`
	CreatedAt           string            `json:"created_at"`
	MuteUntil           *string           `json:"mute_until"`
	FaviconURL          *string           `json:"favicon_url"`
	CustomHeaders       map[string]string `json:"custom_headers"`
	HTTPVerb            string            `json:"http_verb"`
	HTTPBody            string            `json:"http_body"`
	SSL                 *SSLInfo          `json:"ssl,omitempty"`
}

// SSLInfo contient les informations SSL d'un check
type SSLInfo struct {
	TestedAt  string  `json:"tested_at"`
	ExpiresAt string  `json:"expires_at"`
	Valid     bool    `json:"valid"`
	Error     *string `json:"error"`
}

// UpdownMetrics représente les métriques d'un check
type UpdownMetrics struct {
	Uptime   *float64 `json:"uptime,omitempty"`
	Apdex    *float64 `json:"apdex,omitempty"`
	Requests *struct {
		Samples             int `json:"samples"`
		Failures            int `json:"failures"`
		Satisfied           int `json:"satisfied"`
		Tolerated           int `json:"tolerated"`
		ByResponseTime      map[string]int `json:"by_response_time"`
	} `json:"requests,omitempty"`
	Timings *struct {
		Redirect    int `json:"redirect"`
		Namelookup  int `json:"namelookup"`
		Connection  int `json:"connection"`
		Handshake   int `json:"handshake"`
		Response    int `json:"response"`
		Total       int `json:"total"`
	} `json:"timings,omitempty"`
}

// UpdownDowntime représente un temps d'arrêt
type UpdownDowntime struct {
	ID         string  `json:"id"`
	DetailsURL string  `json:"details_url"`
	Error      string  `json:"error"`
	StartedAt  string  `json:"started_at"`
	EndedAt    *string `json:"ended_at"`
	Duration   *int    `json:"duration"`
	Partial    bool    `json:"partial"`
}

// QueryModel représente une requête UpDown.io
type QueryModel struct {
	QueryType   string `json:"queryType"`
	CheckToken  string `json:"checkToken,omitempty"`
	GroupBy     string `json:"groupBy,omitempty"`
}

func LoadPluginSettings(source backend.DataSourceInstanceSettings) (*PluginSettings, error) {
	settings := PluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	// Valeur par défaut pour l'URL de l'API
	if settings.ApiUrl == "" {
		settings.ApiUrl = "https://updown.io/api"
	}

	settings.Secrets = loadSecretPluginSettings(source.DecryptedSecureJSONData)

	return &settings, nil
}

func loadSecretPluginSettings(source map[string]string) *SecretPluginSettings {
	return &SecretPluginSettings{
		ApiKey: source["apiKey"],
	}
}
