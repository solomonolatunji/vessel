import { createFileRoute, Link } from '@tanstack/react-router';
import { LoginForm, OAuthButtons } from '#/features/auth';

import { useGetPublicSettings } from '#/hooks/useSettings';

export const Route = createFileRoute('/_auth/login')({
  component: LoginPage,
});

function LoginPage() {
  const { data: publicSettings } = useGetPublicSettings();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col space-y-2 text-center mb-8">
        <h1 className="text-3xl font-semibold tracking-tight text-foreground">Welcome back</h1>
        <p className="text-sm text-muted-foreground">Sign in to your Vessl workspace.</p>
      </div>

      <OAuthButtons />
      <LoginForm />

      {registrationEnabled && (
        <p className="mt-8 text-center text-sm text-muted-foreground">
          Need an Account?{' '}
          <Link to="/register" className="font-semibold text-primary hover:underline">
            Sign up
          </Link>
        </p>
      )}
    </div>
  );
}
