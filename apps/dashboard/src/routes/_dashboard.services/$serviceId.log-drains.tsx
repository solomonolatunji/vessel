import { createFileRoute } from '@tanstack/react-router';
import { LogDrainsCard } from '#/features/services/log-drains-card';

export const Route = createFileRoute('/_dashboard/services/$serviceId/log-drains')({
  component: LogDrainsRoute,
});

function LogDrainsRoute() {
  const { serviceId } = Route.useParams();

  return (
    <div className="max-w-5xl space-y-6">
      <div>
        <h1 className="font-bold text-2xl">Log Drains</h1>
        <p className="mt-1 text-muted-foreground">
          Forward your container logs to external observability tools.
        </p>
      </div>

      <LogDrainsCard serviceId={serviceId} />
    </div>
  );
}
