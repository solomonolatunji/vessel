import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/services/$serviceId/build')({
  component: () => <div>Route Component</div>,
});
