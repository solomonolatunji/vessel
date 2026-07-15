import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/projects/$projectId/compose')({
  component: () => <div>Route Component</div>,
});
