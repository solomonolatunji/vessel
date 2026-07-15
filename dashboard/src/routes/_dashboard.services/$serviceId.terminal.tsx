import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/services/$serviceId/terminal')({
  component: () => <div>Route Component</div>,
});
