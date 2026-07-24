import { useQueryClient } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router';
import { Copy, RefreshCw, Trash2, Webhook } from 'lucide-react';
import { toast } from 'sonner';
import { v4 as uuidv4 } from 'uuid';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { useGetApp, useUpdateApp } from '#/hooks/useApps';

export const Route = createFileRoute('/_dashboard/services/$serviceId/webhooks')({
  component: WebhooksRoute,
});

function WebhooksRoute() {
  const { serviceId } = Route.useParams();
  const queryClient = useQueryClient();
  const { data: appData } = useGetApp(serviceId);
  const { mutateAsync: updateApp, isPending } = useUpdateApp();

  const app = appData?.data;

  if (!app) return null;

  const handleGenerateToken = async () => {
    try {
      const newToken = uuidv4();
      await updateApp({
        appId: serviceId,
        payload: {
          ...app,
          deployToken: newToken,
        },
      });
      toast.success('Generated new deploy token');
      queryClient.invalidateQueries({ queryKey: ['apps', 'getApp', serviceId] });
    } catch (error: any) {
      toast.error(error?.message || 'Failed to generate token');
    }
  };

  const handleRevokeToken = async () => {
    try {
      await updateApp({
        appId: serviceId,
        payload: {
          ...app,
          deployToken: '',
        },
      });
      toast.success('Revoked deploy token');
      queryClient.invalidateQueries({ queryKey: ['apps', 'getApp', serviceId] });
    } catch (error: any) {
      toast.error(error?.message || 'Failed to revoke token');
    }
  };

  const webhookUrl = `${window.location.protocol}//${window.location.host}/api/webhooks/git/services/${serviceId}?token=${app.deployToken || '<YOUR_TOKEN>'}`;

  return (
    <div className="mx-auto max-w-5xl space-y-6 p-6">
      <div>
        <h2 className="font-bold text-2xl">Webhooks</h2>
        <p className="text-muted-foreground">
          Trigger deployments from external systems automatically.
        </p>
      </div>

      <Card className="border-border/50 bg-card/40">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
              <Webhook className="h-5 w-5" />
            </div>
            <div>
              <CardTitle>Manual Deploy Webhook</CardTitle>
              <CardDescription>
                Trigger a deployment by sending a POST request to this URL.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-2">
            <div className="font-medium text-sm">Deploy Webhook URL</div>
            <div className="flex gap-2">
              <Input
                readOnly
                value={app.deployToken ? webhookUrl : 'Generate a token to view your webhook URL'}
                className="font-mono text-sm"
              />
              <Button
                variant="outline"
                disabled={!app.deployToken}
                onClick={() => {
                  navigator.clipboard.writeText(webhookUrl);
                  toast.success('Copied to clipboard');
                }}
              >
                <Copy className="h-4 w-4" />
              </Button>
            </div>
          </div>

          <div className="space-y-2">
            <div className="font-medium text-sm">cURL Example</div>
            <div className="rounded-lg border border-border/50 bg-black/50 p-4 font-mono text-muted-foreground text-sm">
              {app.deployToken ? (
                <span>
                  <span className="text-primary">curl</span> -X POST {webhookUrl}
                </span>
              ) : (
                <span>Generate a token to see example.</span>
              )}
            </div>
          </div>

          <div className="flex items-center gap-4 border-border/50 border-t pt-6">
            <Button
              onClick={handleGenerateToken}
              disabled={isPending}
              variant={app.deployToken ? 'outline' : 'default'}
              className="gap-2"
            >
              <RefreshCw className="h-4 w-4" />
              {app.deployToken ? 'Regenerate Token' : 'Generate Token'}
            </Button>

            {app.deployToken && (
              <Button
                variant="destructive"
                onClick={handleRevokeToken}
                disabled={isPending}
                className="gap-2"
              >
                <Trash2 className="h-4 w-4" />
                Revoke Token
              </Button>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
