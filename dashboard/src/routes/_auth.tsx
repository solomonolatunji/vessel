import { createFileRoute, Outlet } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth')({
  component: () => (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-background px-4 py-12">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(#6d28d9_0.8px,transparent_1px)] bg-size-[40px_40px] opacity-10 dark:opacity-20" />
      <div className="relative w-full max-w-md">
        <Outlet />
      </div>
    </div>
  ),
});
