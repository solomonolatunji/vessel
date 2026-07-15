import { createFileRoute, Navigate } from '@tanstack/react-router';
import { SetupForm } from '#/features/auth/setup-form';
import { useGetSetupStatus } from '#/hooks/useSettings';

export const Route = createFileRoute('/_auth/setup')({
  component: SetupPage,
});

function SetupPage() {
  const { data: setupStatus, isLoading } = useGetSetupStatus();

  if (!isLoading && !setupStatus?.data?.setupRequired) {
    return <Navigate to="/signin" replace />;
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-500">
      <div className="mb-8 flex flex-col space-y-2 text-center">
        <h1 className="font-semibold text-3xl text-foreground tracking-tight">Welcome to Vessl</h1>
        <p className="text-muted-foreground text-sm">
          Set up the initial owner account to manage your cluster.
        </p>
      </div>

      <SetupForm />
    </div>
  );
}
