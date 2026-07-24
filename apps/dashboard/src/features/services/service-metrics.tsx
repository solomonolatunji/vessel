import { Activity, Cpu, HardDrive, Network, RefreshCw } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Card, CardContent } from '#/components/ui/card';
import { useGetMetrics } from '#/hooks/useDeployments';
import { useHistoricalMetrics } from '#/hooks/useMetrics';
import { MetricChartCard } from './metric-chart-card';
import { formatBytes, formatBytesPerSec, getLatestValue, parseMetricSeries } from './metrics-utils';

const RANGES = [
  { label: '1h', value: '1h' },
  { label: '24h', value: '24h' },
  { label: '7d', value: '7d' },
];

interface ServiceMetricsPageProps {
  serviceId: string;
}

export function ServiceMetricsPage({ serviceId }: ServiceMetricsPageProps) {
  const [range, setRange] = useState('1h');

  const { data: liveRes, isLoading: liveLoading, refetch, isFetching } = useGetMetrics(serviceId);

  const { data: historicalRes, isLoading: historicalLoading } = useHistoricalMetrics(
    serviceId,
    range
  );

  const live = liveRes?.data;
  const historical = historicalRes?.data;

  const cpuSeries = parseMetricSeries(historical?.cpu);
  const memSeries = parseMetricSeries(historical?.memory);
  const netRxSeries = parseMetricSeries(historical?.network_rx);
  const netTxSeries = parseMetricSeries(historical?.network_tx);

  const statusColor =
    live?.status === 'running'
      ? 'bg-emerald-500/15 text-emerald-600 border-emerald-500/30'
      : 'bg-rose-500/15 text-rose-600 border-rose-500/30';

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-bold text-2xl">Metrics</h1>
          <p className="text-muted-foreground text-sm">Real-time and historical resource usage</p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => refetch()}
          disabled={isFetching}
          className="gap-2"
        >
          <RefreshCw className={`h-3.5 w-3.5 ${isFetching ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <LiveStatCard
          label="Status"
          icon={<Activity className="h-4 w-4" />}
          loading={liveLoading}
          value={
            <span className={`rounded-full border px-2 py-0.5 font-medium text-xs ${statusColor}`}>
              {live?.status ?? 'unknown'}
            </span>
          }
        />
        <LiveStatCard
          label="CPU Usage"
          icon={<Cpu className="h-4 w-4" />}
          loading={liveLoading}
          value={`${(live?.cpuUsagePercentage ?? live?.cpuPercent ?? 0).toFixed(2)}%`}
        />
        <LiveStatCard
          label="Memory Used"
          icon={<HardDrive className="h-4 w-4" />}
          loading={liveLoading}
          value={formatBytes(
            live?.memoryUsageBytes ?? (live?.memoryMB ? live.memoryMB * 1024 * 1024 : 0)
          )}
          sub={live?.memoryLimitBytes ? `/ ${formatBytes(live.memoryLimitBytes)} limit` : undefined}
        />
        <LiveStatCard
          label="Network I/O"
          icon={<Network className="h-4 w-4" />}
          loading={liveLoading}
          value={`↓ ${formatBytesPerSec(getLatestValue(netRxSeries))}`}
          sub={`↑ ${formatBytesPerSec(getLatestValue(netTxSeries))}`}
        />
      </div>

      <div className="flex items-center gap-2">
        <span className="text-muted-foreground text-sm">Time range:</span>
        <div className="flex rounded-lg border bg-muted/40 p-0.5">
          {RANGES.map((r) => (
            <button
              key={r.value}
              type="button"
              onClick={() => setRange(r.value)}
              className={`rounded-md px-3 py-1 font-medium text-sm transition-all ${
                range === r.value
                  ? 'bg-background text-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground'
              }`}
            >
              {r.label}
            </button>
          ))}
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        <MetricChartCard
          title="CPU Usage"
          icon={<Cpu className="h-4 w-4" />}
          badge={`${(getLatestValue(cpuSeries) * 100).toFixed(2)}%`}
          data={cpuSeries.map((p) => ({ ...p, value: p.value * 100 }))}
          isLoading={historicalLoading}
          color="hsl(var(--primary))"
          formatY={(v) => `${v.toFixed(0)}%`}
          formatTooltip={(v) => `${v.toFixed(2)}%`}
        />
        <MetricChartCard
          title="Memory Usage"
          icon={<HardDrive className="h-4 w-4" />}
          badge={formatBytes(getLatestValue(memSeries))}
          data={memSeries}
          isLoading={historicalLoading}
          color="hsl(142 71% 45%)"
          formatY={(v) => formatBytes(v, 0)}
          formatTooltip={(v) => formatBytes(v)}
        />
        <MetricChartCard
          title="Network Inbound (Rx)"
          icon={<Network className="h-4 w-4" />}
          badge={formatBytesPerSec(getLatestValue(netRxSeries))}
          data={netRxSeries}
          isLoading={historicalLoading}
          color="hsl(217 91% 60%)"
          formatY={(v) => formatBytes(v, 0)}
          formatTooltip={(v) => formatBytesPerSec(v)}
        />
        <MetricChartCard
          title="Network Outbound (Tx)"
          icon={<Network className="h-4 w-4" />}
          badge={formatBytesPerSec(getLatestValue(netTxSeries))}
          data={netTxSeries}
          isLoading={historicalLoading}
          color="hsl(38 92% 50%)"
          formatY={(v) => formatBytes(v, 0)}
          formatTooltip={(v) => formatBytesPerSec(v)}
        />
      </div>
    </div>
  );
}

interface LiveStatCardProps {
  label: string;
  icon: React.ReactNode;
  loading: boolean;
  value: React.ReactNode;
  sub?: string;
}

function LiveStatCard({ label, icon, loading, value, sub }: LiveStatCardProps) {
  return (
    <Card>
      <CardContent className="pt-5">
        <div className="mb-2 flex items-center gap-2 text-muted-foreground text-xs uppercase tracking-wide">
          {icon}
          {label}
        </div>
        {loading ? (
          <div className="h-7 w-24 animate-pulse rounded bg-muted" />
        ) : (
          <>
            <div className="font-bold text-xl">{value}</div>
            {sub && <p className="mt-0.5 text-muted-foreground text-xs">{sub}</p>}
          </>
        )}
      </CardContent>
    </Card>
  );
}
