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
        <p className="mb-4 font-semibold text-foreground text-lg tracking-tight">Invalid Request</p>
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
    <div className="fade-in slide-in-from-bottom-4 animate-in space-y-8 duration-700">
      <div className="relative rounded-2xl border border-border/80 bg-card/70 p-6 shadow-2xl shadow-black/10 backdrop-blur-xl dark:shadow-black/40">
        <div className="mb-5 flex items-center gap-3">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-linear-to-br from-primary via-purple-600 to-violet-600 shadow-md shadow-primary/25">
            <span className="font-bold text-base text-white tracking-tighter">V</span>
          </div>
          <div>
            <p className="text-muted-foreground/70 text-xs uppercase tracking-wider">
              VESSL ACCESS
            </p>
            <p className="font-semibold text-foreground text-lg tracking-tight">New password</p>
          </div>
        </div>

        <ResetPasswordForm token={token} />
      </div>

      <div className="text-center">
        <div className="inline-flex items-center gap-2 text-[10px] text-muted-foreground/60 uppercase tracking-[2px]">
          <div className="h-px w-8 bg-border" />
          SECURE ACCESS
          <div className="h-px w-8 bg-border" />
        </div>
      </div>
    </div>
  );
}
