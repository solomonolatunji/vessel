import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/settings/maintenance')({
  component: () => <div>Route Component</div>,
});
