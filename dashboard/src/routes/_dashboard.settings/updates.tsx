import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/settings/updates')({
  component: () => <div>Route Component</div>,
});
