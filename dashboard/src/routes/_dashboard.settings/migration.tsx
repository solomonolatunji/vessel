import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/settings/migration')({
  component: () => <div>Route Component</div>,
});
