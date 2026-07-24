import { formatDistanceToNow } from 'date-fns';
import { ExternalLink, GitPullRequest, Loader2 } from 'lucide-react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { useGetApp, useUpdateApp } from '#/hooks/useApps';
import { useListPRPreviews } from '#/hooks/usePRPreviews';
import type { PRPreviewStatus } from '#/interfaces/deployment';

export function PRPreviews({ serviceId }: { serviceId: string }) {
  const { data: appData, isLoading: isAppLoading } = useGetApp(serviceId);
  const { data: previewsData, isLoading: isPreviewsLoading } = useListPRPreviews(serviceId);
  const updateApp = useUpdateApp();

  if (isAppLoading || isPreviewsLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const app = appData?.data;
  if (!app) return <div>Service not found.</div>;

  const previews = previewsData?.data || [];

  const handleToggle = (checked: boolean) => {
    updateApp.mutate({
      appId: serviceId,
      payload: {
        ...app,
        enablePRPreviews: checked,
      },
    });
  };

  const getStatusBadge = (status: PRPreviewStatus) => {
    switch (status) {
      case 'active':
        return (
          <Badge variant="default" className="bg-green-500/10 text-green-500 hover:bg-green-500/20">
            Ready
          </Badge>
        );
      case 'building':
        return (
          <Badge variant="secondary" className="bg-blue-500/10 text-blue-500 hover:bg-blue-500/20">
            Building
          </Badge>
        );
      case 'failed':
        return <Badge variant="destructive">Failed</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <CardTitle>PR Previews</CardTitle>
              <CardDescription>
                Automatically deploy isolated environments for every pull request.
              </CardDescription>
            </div>
            <div className="flex items-center space-x-2">
              <Switch
                id="pr-previews"
                checked={app.enablePRPreviews}
                onCheckedChange={handleToggle}
                disabled={updateApp.isPending}
              />
              <Label htmlFor="pr-previews">Enable</Label>
            </div>
          </div>
        </CardHeader>
      </Card>

      {app.enablePRPreviews && (
        <Card>
          <CardHeader>
            <CardTitle>Active Previews</CardTitle>
            <CardDescription>Currently deployed pull request environments.</CardDescription>
          </CardHeader>
          <CardContent>
            {previews.length === 0 ? (
              <div className="rounded-lg border border-dashed py-12 text-center text-muted-foreground">
                <GitPullRequest className="mx-auto mb-3 h-8 w-8 opacity-20" />
                <p>No active PR previews found.</p>
                <p className="text-sm">Open a pull request to see it deployed here.</p>
              </div>
            ) : (
              <div className="space-y-4">
                {previews.map((preview) => (
                  <div
                    key={preview.id}
                    className="flex items-center justify-between rounded-lg border bg-card/50 p-4"
                  >
                    <div className="flex items-center space-x-4">
                      <div className="rounded-full bg-primary/10 p-2">
                        <GitPullRequest className="h-5 w-5 text-primary" />
                      </div>
                      <div>
                        <div className="flex items-center space-x-2">
                          <h4 className="font-medium text-sm">PR #{preview.prNumber}</h4>
                          <span className="font-mono text-muted-foreground text-xs">
                            {preview.commitHash?.substring(0, 7)}
                          </span>
                        </div>
                        <p className="mt-1 text-muted-foreground text-xs">
                          Updated{' '}
                          {formatDistanceToNow(new Date(preview.updatedAt), { addSuffix: true })}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center space-x-4">
                      {getStatusBadge(preview.status)}

                      {preview.previewDomain && preview.status === 'active' && (
                        <Button variant="outline" size="sm" asChild>
                          <a
                            href={`https://${preview.previewDomain}`}
                            target="_blank"
                            rel="noreferrer"
                            className="flex items-center"
                          >
                            <ExternalLink className="mr-2 h-4 w-4" />
                            Visit
                          </a>
                        </Button>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
