import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/settings/maintenance')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_shell/settings/maintenance"!</div>;
}
