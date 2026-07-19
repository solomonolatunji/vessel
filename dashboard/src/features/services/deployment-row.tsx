import { formatDistanceToNow } from 'date-fns';
import { Badge } from '#/components/ui/badge';
import { Card, CardContent } from '#/components/ui/card';

export interface Deployment {
  id: string;
  status: 'created' | 'building' | 'deploying' | 'running' | 'stopped' | 'error';
  createdAt: string;
  activeContainerId?: string;
  startingContainerId?: string;
  healthCheckStatus?: 'probing' | 'healthy' | 'unhealthy';
}

interface DeploymentRowProps {
  deployment: Deployment;
  onViewLogs?: (id: string) => void;
}

export function DeploymentRow({ deployment, onViewLogs }: DeploymentRowProps) {
  const getStatusColor = (status: Deployment['status']) => {
    switch (status) {
      case 'running':
        return 'bg-green-500/10 text-green-500 hover:bg-green-500/20';
      case 'deploying':
        return 'bg-blue-500/10 text-blue-500 hover:bg-blue-500/20';
      case 'building':
        return 'bg-yellow-500/10 text-yellow-500 hover:bg-yellow-500/20';
      case 'error':
        return 'bg-red-500/10 text-red-500 hover:bg-red-500/20';
      default:
        return 'bg-zinc-500/10 text-zinc-500 hover:bg-zinc-500/20';
    }
  };

  return (
    <Card className="mb-2">
      <CardContent className="flex items-center justify-between p-4">
        <div className="flex flex-col gap-2">
          <div className="flex items-center gap-3">
            <span className="font-medium font-mono text-sm">{deployment.id.slice(0, 8)}</span>
            <Badge variant="secondary" className={getStatusColor(deployment.status)}>
              {deployment.status.toUpperCase()}
            </Badge>
            <span className="text-xs text-zinc-500">
              {formatDistanceToNow(new Date(deployment.createdAt), { addSuffix: true })}
            </span>
          </div>

          {deployment.status === 'deploying' && (
            <div className="mt-2 rounded-md border border-blue-200 bg-blue-50/50 p-3 dark:border-blue-900 dark:bg-blue-950/20">
              <p className="mb-2 font-medium text-blue-600 text-xs dark:text-blue-400">
                Zero-Downtime Hot Swap in Progress
              </p>
              <div className="flex items-center gap-4 text-sm">
                {deployment.activeContainerId && (
                  <div className="flex items-center gap-2">
                    <span className="h-2 w-2 rounded-full bg-green-500" />
                    <span className="text-zinc-600 dark:text-zinc-400">
                      Active:{' '}
                      <span className="font-mono">{deployment.activeContainerId.slice(0, 8)}</span>
                    </span>
                  </div>
                )}
                {deployment.startingContainerId && (
                  <div className="flex items-center gap-2 border-l pl-4 dark:border-zinc-700">
                    <span className="h-2 w-2 animate-pulse rounded-full bg-blue-500" />
                    <span className="text-zinc-600 dark:text-zinc-400">
                      Probing:{' '}
                      <span className="font-mono">
                        {deployment.startingContainerId.slice(0, 8)}
                      </span>
                    </span>
                    {deployment.healthCheckStatus && (
                      <Badge variant="outline" className="ml-2 text-xs">
                        Health: {deployment.healthCheckStatus}
                      </Badge>
                    )}
                  </div>
                )}
              </div>
            </div>
          )}
        </div>

        <div>
          <button
            type="button"
            className="font-medium text-sm text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
            onClick={() => onViewLogs && onViewLogs(deployment.id)}
          >
            View Logs
          </button>
        </div>
      </CardContent>
    </Card>
  );
}
