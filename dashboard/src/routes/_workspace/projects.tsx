import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_workspace/projects')({
  component: () => <div>Route /_workspace/projects</div>,
});
