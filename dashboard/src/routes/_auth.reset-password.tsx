import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { ResetPasswordForm } from '#/features/auth/reset-password-form';

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
      <div className="fade-in slide-in-from-bottom-4 animate-in text-center duration-500">
        <h1 className="mb-4 font-semibold text-3xl text-foreground tracking-tight">
          Invalid Request
        </h1>
        <p className="mb-8 text-muted-foreground text-sm">
          The password reset token is missing. Please check your email link again.
        </p>
        <button
          type="button"
          onClick={() => navigate({ to: '/signin' })}
          className="font-medium text-primary text-sm hover:underline"
        >
          Return to sign in
        </button>
      </div>
    );
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-500">
      <div className="mb-8 flex flex-col space-y-2 text-center">
        <h1 className="font-semibold text-3xl text-foreground tracking-tight">
          Create new password
        </h1>
        <p className="text-muted-foreground text-sm">
          Your new password must be different from previous used passwords.
        </p>
      </div>

      <ResetPasswordForm token={token} />
    </div>
  );
}
