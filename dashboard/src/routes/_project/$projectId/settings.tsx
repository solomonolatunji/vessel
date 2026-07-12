import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_project/$projectId/settings')({
  component: () => <div>Route /_project/$projectId/settings</div>,
});
