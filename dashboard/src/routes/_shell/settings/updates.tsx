import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/settings/updates')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_shell/settings/updates"!</div>;
}
