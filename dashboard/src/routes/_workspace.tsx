import { createFileRoute, Outlet } from '@tanstack/react-router';
import { Shell } from '#/components/layout/Shell';

export const Route = createFileRoute('/_workspace')({
  component: () => (
    <Shell>
      <Outlet />
    </Shell>
  ),
});
