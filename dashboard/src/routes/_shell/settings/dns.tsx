import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/settings/dns')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_shell/settings/dns"!</div>;
}
