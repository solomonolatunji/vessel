import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_workspace/settings')({
  component: () => <div>Route /_workspace/settings</div>,
});
