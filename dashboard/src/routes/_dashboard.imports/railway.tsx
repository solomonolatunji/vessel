import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/imports/railway')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_shell/imports/railway"!</div>;
}
