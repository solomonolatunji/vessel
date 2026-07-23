import { createFileRoute, Link, Navigate } from '@tanstack/react-router';
import { OAuthButtons } from '#/features/auth/o-auth-buttons';
import { RegisterForm } from '#/features/auth/register-form';
import { useGetPublicSettings, useGetSetupStatus } from '#/hooks/useSettings';
import { useSystemStore } from '#/stores/systemStore';

export const Route = createFileRoute('/_auth/signup')({
  component: RegisterPage,
  head: () => {
    const siteName = useSystemStore.getState().siteName;
    return { meta: [{ title: `Sign Up - ${siteName}` }] };
  },
});

function RegisterPage() {
  const { data: publicSettings } = useGetPublicSettings();
  const { data: setupStatus, isLoading } = useGetSetupStatus();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;

  if (!isLoading && setupStatus?.data?.setupRequired) {
    return <Navigate to="/onboarding" replace />;
  }

  if (!isLoading && !registrationEnabled) {
    return <Navigate to="/signin" replace />;
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in space-y-6 duration-700">
      <div className="relative rounded-2xl border border-border/80 bg-card/70 p-5 shadow-2xl shadow-black/10 backdrop-blur-xl sm:p-6 dark:shadow-black/40">
        <div className="mb-4 flex items-center gap-3">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-linear-to-br from-primary via-purple-600 to-violet-600 shadow-md shadow-primary/25">
            <span className="font-bold text-base text-white tracking-tighter">V</span>
          </div>
          <div>
            <p className="text-muted-foreground/70 text-xs uppercase tracking-wider">
              CODEDOCK ACCESS
            </p>
            <p className="font-semibold text-foreground text-lg tracking-tight">Sign up</p>
          </div>
        </div>

        <OAuthButtons />
        <RegisterForm />

        <p className="mt-5 text-center text-muted-foreground text-sm">
          Already have an account?{' '}
          <Link
            to="/signin"
            className="font-medium text-primary underline-offset-4 transition-colors hover:text-primary-hover hover:underline"
          >
            Sign in
          </Link>
        </p>
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
