import type { MetricData } from '#/interfaces/metrics';

export function parseMetricSeries(data: MetricData | undefined): { time: number; value: number }[] {
  if (!data?.result?.length) return [];
  const series = data.result[0];
  return series.values.map(([ts, val]) => ({
    time: ts * 1000,
    value: parseFloat(val) || 0,
  }));
}

export function formatBytes(bytes: number, decimals = 1): string {
  if (!+bytes) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / k ** i).toFixed(decimals))} ${sizes[i]}`;
}

export function formatBytesPerSec(bytes: number): string {
  return `${formatBytes(bytes)}/s`;
}

export function formatTimestamp(ms: number): string {
  return new Date(ms).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

export function getLatestValue(series: { time: number; value: number }[]): number {
  if (!series.length) return 0;
  return series[series.length - 1].value;
}
