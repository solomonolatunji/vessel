import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/settings/users')({
  component: () => <div>Route Component</div>,
});
