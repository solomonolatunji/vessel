import { createFileRoute, Link } from '@tanstack/react-router';
import { AlertCircle } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { ForgotPasswordForm } from '#/features/auth/forgot-password-form';
import { useGetPublicSettings } from '#/hooks/useSettings';
import { useSystemStore } from '#/stores/systemStore';

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPasswordPage,
  head: () => {
    const siteName = useSystemStore.getState().siteName;
    return { meta: [{ title: `Reset Password - ${siteName}` }] };
  },
});

function ForgotPasswordPage() {
  const { data, isLoading } = useGetPublicSettings();
  const emailEnabled = data?.data?.emailEnabled;

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in space-y-6 duration-700">
      <div className="relative rounded-2xl border border-border/80 bg-card/70 p-5 shadow-2xl shadow-black/10 backdrop-blur-xl sm:p-6 dark:shadow-black/40">
        <div className="mb-4 flex items-center gap-3">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-linear-to-br from-primary via-purple-600 to-violet-600 shadow-md shadow-primary/25">
            <span className="font-bold text-base text-white tracking-tighter">V</span>
          </div>
          <div>
            <p className="text-muted-foreground/70 text-xs uppercase tracking-wider">
              VESSL ACCESS
            </p>
            <p className="font-semibold text-foreground text-lg tracking-tight">Reset password</p>
          </div>
        </div>

        {!isLoading && emailEnabled === false ? (
          <div className="flex flex-col items-center gap-4 py-5 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
              <AlertCircle className="h-6 w-6 text-destructive" />
            </div>
            <div className="space-y-1">
              <p className="font-semibold text-base text-foreground tracking-tight">
                Email not configured
              </p>
              <p className="text-muted-foreground text-sm">
                Your team is yet to set or enable email. Please contact your administrator.
              </p>
            </div>
            <Button asChild variant="outline" className="mt-2">
              <Link to="/signin">Back to sign in</Link>
            </Button>
          </div>
        ) : (
          <ForgotPasswordForm />
        )}
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
