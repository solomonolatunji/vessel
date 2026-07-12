import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_workspace/databases')({
  component: () => <div>Route /_workspace/databases</div>,
});
