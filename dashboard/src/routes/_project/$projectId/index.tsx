import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_project/$projectId/')({
  component: () => <div>Route /_project/$projectId/</div>,
});
