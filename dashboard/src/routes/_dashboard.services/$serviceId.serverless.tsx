import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/services/$serviceId/serverless')({
  component: () => <div>Route Component</div>,
});
