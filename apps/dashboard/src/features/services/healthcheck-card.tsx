import { useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useGetApp, useUpdateApp } from '#/hooks/useApps';

export function HealthcheckCard({ serviceId }: { serviceId: string }) {
  const queryClient = useQueryClient();
  const { data: appData } = useGetApp(serviceId);
  const app = appData?.data;

  const [internalPort, setInternalPort] = useState(app?.internalPort?.toString() || '3000');
  const [healthCheckPath, setHealthCheckPath] = useState(app?.healthCheckPath || '/');

  const { mutateAsync: updateApp, isPending } = useUpdateApp();

  const handleSave = async () => {
    if (!app) return;
    try {
      await updateApp({
        appId: serviceId,
        payload: {
          ...app,
          internalPort: parseInt(internalPort, 10) || 3000,
          healthCheckPath,
        },
      });
      toast.success('Healthcheck settings saved');
      queryClient.invalidateQueries({ queryKey: ['apps', 'getApp', serviceId] });
    } catch (error: any) {
      toast.error(error?.message || 'Failed to save healthcheck settings');
    }
  };

  if (!app) return null;

  return (
    <Card className="border-border/50 bg-card/40">
      <CardHeader>
        <CardTitle>Networking & Healthchecks</CardTitle>
        <CardDescription>
          Configure how the proxy routes traffic to your container and how we determine if it's
          healthy.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid gap-6 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="internalPort">Internal Port</Label>
            <Input
              id="internalPort"
              type="number"
              placeholder="3000"
              value={internalPort}
              onChange={(e) => setInternalPort(e.target.value)}
              className="font-mono"
            />
            <p className="text-muted-foreground text-xs">
              The port your application server binds to (e.g., 3000, 8080).
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="healthCheckPath">Healthcheck Path</Label>
            <Input
              id="healthCheckPath"
              placeholder="/"
              value={healthCheckPath}
              onChange={(e) => setHealthCheckPath(e.target.value)}
              className="font-mono"
            />
            <p className="text-muted-foreground text-xs">
              The endpoint we'll ping to verify readiness (e.g., /health, /api/status).
            </p>
          </div>
        </div>

        <div className="flex justify-end border-border/50 border-t pt-4">
          <Button onClick={handleSave} disabled={isPending}>
            {isPending ? 'Saving...' : 'Save Settings'}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
