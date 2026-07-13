import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';
import { Shell } from '#/components/layout/shell';
import { useListWorkspaces } from '#/hooks/useWorkspaces';
import { authStore } from '#/stores/authStore';

export const Route = createFileRoute('/_workspace')({
  beforeLoad: () => {
    if (!authStore.state.isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: WorkspaceLayout,
});

function WorkspaceLayout() {
  useListWorkspaces();

  return (
    <Shell>
      <Outlet />
    </Shell>
  );
}
