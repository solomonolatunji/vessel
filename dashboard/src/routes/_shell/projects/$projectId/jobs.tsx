import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/projects/$projectId/jobs')({
  component: () => <div>Route Component</div>,
});
