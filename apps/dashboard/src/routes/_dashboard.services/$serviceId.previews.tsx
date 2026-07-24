import { createFileRoute } from '@tanstack/react-router';
import { PRPreviews } from '#/features/services/pr-previews';

export const Route = createFileRoute('/_dashboard/services/$serviceId/previews')({
  component: ServicePRPreviewsRoute,
});

function ServicePRPreviewsRoute() {
  const { serviceId } = Route.useParams();

  return (
    <div className="space-y-6">
      <PRPreviews serviceId={serviceId} />
    </div>
  );
}
