import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface UpdownQuery extends DataQuery {
  // Query type: 'checks', 'metrics', 'performance', 'downtimes', 'uptime'
  queryType: string;
  // Check token for specific queries
  checkToken?: string;
  // Time range for metrics (from/to)
  timeRange?: {
    from: string;
    to: string;
  };
  // Grouping for metrics ('time' or 'host')
  groupBy?: string;
}

export const DEFAULT_QUERY: Partial<UpdownQuery> = {
  queryType: 'checks',
};

// Interface for an UpDown.io check
export interface UpdownCheck {
  token: string;
  url: string;
  type: string;
  alias: string;
  last_status: number;
  uptime: number;
  down: boolean;
  down_since: string | null;
  up_since: string | null;
  error: string | null;
  period: number;
  apdex_t: number;
  string_match: string;
  enabled: boolean;
  published: boolean;
  disabled_locations: string[];
  recipients: string[];
  last_check_at: string;
  next_check_at: string;
  created_at: string;
  mute_until: string | null;
  favicon_url: string | null;
  custom_headers: Record<string, string>;
  http_verb: string;
  http_body: string;
  ssl?: {
    tested_at: string;
    expires_at: string;
    valid: boolean;
    error: string | null;
  };
}

// Interface for downtimes
export interface UpdownDowntime {
  id: string;
  details_url: string;
  error: string;
  started_at: string;
  ended_at: string | null;
  duration: number | null;
  partial: boolean;
}

// Interface for metrics
export interface UpdownMetrics {
  uptime?: number;
  apdex?: number;
  requests?: {
    samples: number;
    failures: number;
    satisfied: number;
    tolerated: number;
    by_response_time: {
      under125: number;
      under250: number;
      under500: number;
      under1000: number;
      under2000: number;
      under4000: number;
    };
  };
  timings?: {
    redirect: number;
    namelookup: number;
    connection: number;
    handshake: number;
    response: number;
    total: number;
  };
}

export interface DataPoint {
  Time: number;
  Value: number;
}

export interface DataSourceResponse {
  datapoints: DataPoint[];
}

/**
 * These are options configured for each DataSource instance
 */
export interface UpdownDataSourceOptions extends DataSourceJsonData {
  // Base API URL (default: https://updown.io/api)
  apiUrl?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface UpdownSecureJsonData {
  apiKey?: string;
}
