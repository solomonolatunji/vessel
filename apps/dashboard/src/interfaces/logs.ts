export interface LogStream {
  stream: Record<string, string>;
  values: [string, string][];
}

export interface GetHistoricalLogsResponse {
  result: LogStream[];
}
