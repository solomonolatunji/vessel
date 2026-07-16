export interface MetricResult {
  metric: Record<string, string>;
  values: [number, string][]; // TSDB usually returns [timestamp_in_seconds_or_ms, value_as_string]
}

export interface MetricData {
  resultType: string;
  result: MetricResult[];
}

export interface GetHistoricalMetricsResponse {
  cpu: MetricData;
  memory: MetricData;
}
