import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/databases/$dbId/')({
  component: () => <div>Route Component</div>,
});
