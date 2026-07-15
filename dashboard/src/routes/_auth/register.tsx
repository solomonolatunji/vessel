import { createFileRoute, Link, Navigate } from '@tanstack/react-router';
import { OAuthButtons } from '#/features/auth/o-auth-buttons';
import { RegisterForm } from '#/features/auth/register-form';
import { useGetPublicSettings, useGetSetupStatus } from '#/hooks/useSettings';

export const Route = createFileRoute('/_auth/register')({
  component: RegisterPage,
});

function RegisterPage() {
  const { data: publicSettings } = useGetPublicSettings();
  const { data: setupStatus, isLoading } = useGetSetupStatus();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;

  if (!isLoading && setupStatus?.data?.setupRequired) {
    return <Navigate to="/setup" replace />;
  }

  if (!isLoading && !registrationEnabled) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-500">
      <div className="mb-8 flex flex-col space-y-2 text-center">
        <h1 className="font-semibold text-3xl text-foreground tracking-tight">Create an account</h1>
        <p className="text-muted-foreground text-sm">Enter your details below to get started.</p>
      </div>

      <OAuthButtons />
      <RegisterForm />

      <p className="mt-8 text-center text-muted-foreground text-sm">
        Already have an account?{' '}
        <Link to="/login" className="font-semibold text-primary hover:underline">
          Sign in
        </Link>
      </p>
    </div>
  );
}
