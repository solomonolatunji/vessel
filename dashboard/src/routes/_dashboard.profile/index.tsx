import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/profile/')({
  component: () => <div>Route Component</div>,
});
