import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/services/$serviceId/')({
  component: () => <div>Route Component</div>,
});
