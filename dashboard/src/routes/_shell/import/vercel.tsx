import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/import/vercel')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_shell/import/vercel"!</div>;
}
