import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/settings/users')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_shell/settings/users"!</div>;
}
