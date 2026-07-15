import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/settings/oauth')({
  component: () => <div>Route Component</div>,
});
