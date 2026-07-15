import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/projects/$projectId/services/$serviceId/build')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_shell/projects/$projectId/services/$serviceId/build"!</div>;
}
