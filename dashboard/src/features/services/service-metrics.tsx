import { Activity, Cpu, HardDrive } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { Skeleton } from '#/components/ui/skeleton';
import { useGetMetrics } from '#/hooks/useDeployments';

interface ServiceMetricsPanelProps {
  serviceId: string;
}

export function ServiceMetricsPanel({ serviceId }: ServiceMetricsPanelProps) {
  const { data: response, isLoading, isError } = useGetMetrics(serviceId);

  if (isError) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-muted-foreground">
          Failed to load metrics. The container might not be running.
        </CardContent>
      </Card>
    );
  }

  const metrics = response?.data;

  return (
    <div className="grid gap-4 md:grid-cols-3">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="font-medium text-sm">Status</CardTitle>
          <Activity className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <Skeleton className="h-8 w-1/2" />
          ) : (
            <div className="font-bold text-2xl capitalize">{metrics?.status || 'running'}</div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="font-medium text-sm">CPU Usage</CardTitle>
          <Cpu className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <Skeleton className="h-8 w-1/2" />
          ) : (
            <div className="font-bold text-2xl">
              {(metrics?.cpuUsagePercentage ?? metrics?.cpuPercent ?? 0).toFixed(2)}%
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="font-medium text-sm">Memory Usage</CardTitle>
          <HardDrive className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <Skeleton className="h-8 w-1/2" />
          ) : (
            <div className="font-bold text-2xl">
              {metrics
                ? formatBytes(
                    metrics.memoryUsageBytes ??
                      (metrics.memoryMB ? metrics.memoryMB * 1024 * 1024 : 0)
                  )
                : '0 B'}
            </div>
          )}
          <p className="mt-1 text-muted-foreground text-xs">
            {metrics?.memoryLimitBytes ? `Limit: ${formatBytes(metrics.memoryLimitBytes)}` : ''}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}

function formatBytes(bytes: number, decimals = 2) {
  if (!+bytes) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / k ** i).toFixed(dm))} ${sizes[i]}`;
}
