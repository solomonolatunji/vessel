import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/databases/$dbId/backups')({
  component: () => <div>Route Component</div>,
});
