import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_workspace/')({
  component: () => <div>Route /_workspace/</div>,
});
