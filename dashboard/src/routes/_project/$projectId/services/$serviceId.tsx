import { createFileRoute } from '@tanstack/react-router';
import { ServiceMetricsPanel } from '#/features/services';

export const Route = createFileRoute('/_project/$projectId/services/$serviceId')({
  component: ServiceDetailsPage,
});

function ServiceDetailsPage() {
  const { serviceId } = Route.useParams();

  return (
    <div className="flex flex-col space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold tracking-tight">Service Details</h1>
      </div>
      
      <div className="grid gap-6">
        <section>
          <h2 className="text-lg font-semibold mb-4">Real-time Metrics</h2>
          <ServiceMetricsPanel serviceId={serviceId} />
        </section>
      </div>
    </div>
  );
}
