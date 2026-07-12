import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_project/$projectId/services/$serviceId')({
  component: () => <div>Route /_project/$projectId/services/$serviceId</div>,
});
