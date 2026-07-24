import { Loader2, Wrench } from 'lucide-react';
import { toast } from 'sonner';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Switch } from '#/components/ui/switch';
import { useGetApp, useRedeployApp, useUpdateApp } from '#/hooks/useApps';

export function MaintenanceModeCard({ serviceId }: { serviceId: string }) {
  const { data: appData, isLoading } = useGetApp(serviceId);
  const updateApp = useUpdateApp();
  const redeployApp = useRedeployApp();

  if (isLoading) {
    return (
      <Card>
        <CardContent className="flex justify-center p-6">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    );
  }

  const app = appData?.data;
  if (!app) return null;

  const handleToggle = async (checked: boolean) => {
    try {
      await updateApp.mutateAsync({
        appId: serviceId,
        payload: { ...app, maintenanceMode: checked },
      });
      toast.success('Maintenance mode updated. Redeploying service to apply changes...');
      await redeployApp.mutateAsync({ appId: serviceId });
      toast.success('Service redeployment started successfully');
    } catch (_error) {
      toast.error('Failed to update maintenance mode');
    }
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <Wrench className="h-5 w-5 text-amber-500" />
          <CardTitle>Maintenance Mode</CardTitle>
        </div>
        <CardDescription>
          Temporarily replace your application with a static maintenance page. Useful for manual
          migrations or heavy updates.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between rounded-lg border p-4">
          <div className="space-y-0.5">
            <h4 className="font-medium text-sm">Enable Maintenance Mode</h4>
            <p className="text-muted-foreground text-sm">
              All web traffic will instantly see a "Under Maintenance" page.
            </p>
          </div>
          <Switch
            checked={!!app.maintenanceMode}
            onCheckedChange={handleToggle}
            disabled={updateApp.isPending || redeployApp.isPending}
          />
        </div>
      </CardContent>
    </Card>
  );
}
