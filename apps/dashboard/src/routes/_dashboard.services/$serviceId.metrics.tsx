import { createFileRoute } from '@tanstack/react-router';
import { ServiceMetricsPage } from '#/features/services/service-metrics';

export const Route = createFileRoute('/_dashboard/services/$serviceId/metrics')({
  component: ServiceMetricsRoute,
});

function ServiceMetricsRoute() {
  const { serviceId } = Route.useParams();
  return <ServiceMetricsPage serviceId={serviceId} />;
}
