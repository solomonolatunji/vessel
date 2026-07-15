import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/import/railway')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_shell/import/railway"!</div>;
}
