import { createFileRoute, Outlet } from '@tanstack/react-router';

export const Route = createFileRoute('/_project')({
  component: () => <Outlet />,
});
