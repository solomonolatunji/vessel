import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/settings/backups')({
  component: () => <div>Route Component</div>,
});
