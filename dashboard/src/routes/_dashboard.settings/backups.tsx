import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/settings/backups')({
  component: () => <div>Route Component</div>,
});
