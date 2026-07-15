import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/services/$serviceId/deployments')({
  component: () => <div>Route Component</div>,
});
