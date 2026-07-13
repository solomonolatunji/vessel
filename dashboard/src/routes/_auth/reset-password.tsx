import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { ResetPasswordForm } from '#/features/auth';

export const Route = createFileRoute('/_auth/reset-password')({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      token: (search.token as string) || '',
    };
  },
  component: ResetPasswordPage,
});

function ResetPasswordPage() {
  const { token } = Route.useSearch();
  const navigate = useNavigate();

  if (!token) {
    return (
      <div className="animate-in fade-in slide-in-from-bottom-4 duration-500 text-center">
        <h1 className="text-3xl font-semibold tracking-tight text-foreground mb-4">
          Invalid Request
        </h1>
        <p className="text-sm text-muted-foreground mb-8">
          The password reset token is missing. Please check your email link again.
        </p>
        <button
          type="button"
          onClick={() => navigate({ to: '/login' })}
          className="text-sm font-medium text-primary hover:underline"
        >
          Return to sign in
        </button>
      </div>
    );
  }

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col space-y-2 text-center mb-8">
        <h1 className="text-3xl font-semibold tracking-tight text-foreground">
          Create new password
        </h1>
        <p className="text-sm text-muted-foreground">
          Your new password must be different from previous used passwords.
        </p>
      </div>

      <ResetPasswordForm token={token} />
    </div>
  );
}
