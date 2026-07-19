import { GitBranch, Loader2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { useDisconnect, useGetStatus } from '#/hooks/useGit';

export function GitAppsManager() {
  const { data: statusResponse, isLoading } = useGetStatus();
  const disconnect = useDisconnect();

  const statuses = statusResponse?.data || [];

  // For now, we only support github, but this can be expanded
  const providers = [{ id: 'github', name: 'GitHub', icon: GitBranch }];

  const handleConnect = (providerId: string) => {
    // In a real implementation, this would redirect to an OAuth URL for the provider
    // e.g., window.location.href = `/api/v1/git/${providerId}/auth`;
    console.log(`Connect to ${providerId}`);
  };

  const handleDisconnect = (providerId: string) => {
    disconnect.mutate({ provider: providerId });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Git Applications</CardTitle>
        <CardDescription>
          Connect your Vessl instance to Git providers to deploy applications from your
          repositories.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center p-6">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : (
          <div className="space-y-4">
            {providers.map((provider) => {
              const status = statuses.find((s) => s.provider === provider.id);
              const isConnected = status?.connected || false;
              const Icon = provider.icon;

              return (
                <div
                  key={provider.id}
                  className="flex items-center justify-between rounded-lg border p-4"
                >
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-muted">
                      <Icon className="h-5 w-5" />
                    </div>
                    <div>
                      <p className="font-medium">{provider.name}</p>
                      <p className="text-muted-foreground text-sm">
                        {isConnected ? 'Connected' : 'Not connected'}
                      </p>
                    </div>
                  </div>
                  <div>
                    {isConnected ? (
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDisconnect(provider.id)}
                        disabled={disconnect.isPending}
                      >
                        Disconnect
                      </Button>
                    ) : (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleConnect(provider.id)}
                      >
                        Connect
                      </Button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
