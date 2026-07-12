import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_project/$projectId/deployments')({
  component: () => <div>Route /_project/$projectId/deployments</div>,
});
