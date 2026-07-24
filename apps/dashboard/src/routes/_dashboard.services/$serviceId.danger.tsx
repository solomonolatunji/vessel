import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { PowerOff, RefreshCw, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { useDeleteApp, useGetApp, useRestartApp, useStopApp } from '#/hooks/useApps';
import {
  useDeleteDatabase,
  useGetDatabase,
  useRestartDatabase,
  useStopDatabase,
} from '#/hooks/useDatabases';

export const Route = createFileRoute('/_dashboard/services/$serviceId/danger')({
  component: DangerZoneRoute,
});

function DangerZoneRoute() {
  const { serviceId } = Route.useParams();
  const navigate = useNavigate();

  const { data: appData } = useGetApp(serviceId);
  const { data: dbData } = useGetDatabase(serviceId);
  const deleteApp = useDeleteApp();
  const deleteDb = useDeleteDatabase();
  const restartApp = useRestartApp();
  const stopApp = useStopApp();
  const restartDb = useRestartDatabase();
  const stopDb = useStopDatabase();

  const app = appData?.data;
  const db = dbData?.data;
  const isDatabase = !!db;
  const resourceName = app?.name || db?.name || 'this service';
  const projectId = app?.projectId || db?.projectId;

  const handleDelete = async () => {
    if (
      !confirm(
        `Are you absolutely sure you want to delete ${resourceName}? This action cannot be undone and will permanently delete all data, volumes, and configurations.`
      )
    ) {
      return;
    }

    try {
      if (isDatabase) {
        await deleteDb.mutateAsync({ id: serviceId });
      } else {
        await deleteApp.mutateAsync({ appId: serviceId });
      }
      toast.success('Service deleted successfully');
      navigate({ to: `/projects/${projectId}` });
    } catch (_error) {
      toast.error('Failed to delete service');
    }
  };

  const handleRestart = async () => {
    try {
      if (isDatabase) {
        await restartDb.mutateAsync({ id: serviceId });
      } else {
        await restartApp.mutateAsync({ appId: serviceId });
      }
      toast.success('Service restart initiated');
    } catch (_error) {
      toast.error('Failed to restart service');
    }
  };

  const handleStop = async () => {
    try {
      if (isDatabase) {
        await stopDb.mutateAsync({ id: serviceId });
      } else {
        await stopApp.mutateAsync({ appId: serviceId });
      }
      toast.success('Service stop initiated');
    } catch (_error) {
      toast.error('Failed to stop service');
    }
  };

  return (
    <div className="max-w-4xl space-y-6">
      <div>
        <h1 className="font-bold text-2xl text-red-500">Danger Zone</h1>
        <p className="mt-1 text-muted-foreground">
          Destructive operations for {resourceName}. Please proceed with caution.
        </p>
      </div>

      <div className="overflow-hidden rounded-xl border border-red-500/20 bg-red-500/5">
        <div className="border-red-500/10 border-b p-6">
          <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
            <div>
              <h3 className="font-semibold text-lg">Restart Service</h3>
              <p className="mt-1 text-muted-foreground text-sm">
                Gracefully restart the container. This may cause a brief period of downtime.
              </p>
            </div>
            <Button
              variant="outline"
              onClick={handleRestart}
              disabled={restartApp.isPending || restartDb.isPending}
            >
              <RefreshCw className="mr-2 h-4 w-4" />
              Restart Container
            </Button>
          </div>
        </div>

        <div className="border-red-500/10 border-b p-6">
          <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
            <div>
              <h3 className="font-semibold text-lg">Stop Service</h3>
              <p className="mt-1 text-muted-foreground text-sm">
                Stop the container. The service will be completely unavailable until started again.
              </p>
            </div>
            <Button
              variant="outline"
              onClick={handleStop}
              disabled={stopApp.isPending || stopDb.isPending}
            >
              <PowerOff className="mr-2 h-4 w-4" />
              Stop Container
            </Button>
          </div>
        </div>

        <div className="p-6">
          <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
            <div>
              <h3 className="font-semibold text-lg text-red-500">Delete Service</h3>
              <p className="mt-1 text-muted-foreground text-sm">
                Permanently delete this service, its database volumes, environment variables, and
                configuration.
              </p>
            </div>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteApp.isPending || deleteDb.isPending}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              Delete {isDatabase ? 'Database' : 'Application'}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
