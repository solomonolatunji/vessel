import { createFileRoute, Outlet } from '@tanstack/react-router';
import { BackgroundPattern } from '#/components/layout/background-pattern';

export const Route = createFileRoute('/_auth')({
  component: () => (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-background px-4 py-12">
      <BackgroundPattern />
      <div className="relative w-full max-w-md">
        <Outlet />
      </div>
    </div>
  ),
});
