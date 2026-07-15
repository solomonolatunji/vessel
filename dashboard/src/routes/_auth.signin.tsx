import { createFileRoute, Link, Navigate } from '@tanstack/react-router';
import { LoginForm } from '#/features/auth/login-form';
import { OAuthButtons } from '#/features/auth/o-auth-buttons';

import { useGetPublicSettings, useGetSetupStatus } from '#/hooks/useSettings';

export const Route = createFileRoute('/_auth/signin')({
  component: LoginPage,
});

function LoginPage() {
  const { data: publicSettings } = useGetPublicSettings();
  const { data: setupStatus, isLoading } = useGetSetupStatus();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;

  if (!isLoading && setupStatus?.data?.setupRequired) {
    return <Navigate to="/setup" replace />;
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-500">
      <div className="mb-8 flex flex-col space-y-2 text-center">
        <h1 className="font-semibold text-3xl text-foreground tracking-tight">Welcome back</h1>
        <p className="text-muted-foreground text-sm">Sign in to your Vessl instance.</p>
      </div>

      <OAuthButtons />
      <LoginForm />

      {registrationEnabled && (
        <p className="mt-8 text-center text-muted-foreground text-sm">
          Need an Account?{' '}
          <Link to="/signup" className="font-semibold text-primary hover:underline">
            Sign up
          </Link>
        </p>
      )}
    </div>
  );
}
