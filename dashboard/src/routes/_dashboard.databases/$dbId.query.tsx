import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/databases/$dbId/query')({
  component: () => <div>Route Component</div>,
});
