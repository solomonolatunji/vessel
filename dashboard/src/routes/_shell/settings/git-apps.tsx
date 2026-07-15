import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/settings/git-apps')({
  component: () => <div>Route Component</div>,
});
