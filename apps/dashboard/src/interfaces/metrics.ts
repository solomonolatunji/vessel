export interface MetricResult {
  metric: Record<string, string>;
  values: [number, string][];
}

export interface MetricData {
  resultType: string;
  result: MetricResult[];
}

export interface GetHistoricalMetricsResponse {
  cpu: MetricData;
  memory: MetricData;
  network_rx: MetricData;
  network_tx: MetricData;
}
