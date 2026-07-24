import { createFileRoute, Navigate } from '@tanstack/react-router';
import { BackgroundPattern } from '#/components/layout/background-pattern';
import { OnboardingForm } from '#/features/onboarding/onboarding-form';
import { useGetSetupStatus } from '#/hooks/useSettings';
import { useSystemStore } from '#/stores/systemStore';

export const Route = createFileRoute('/onboarding')({
  component: OnboardingPage,
  head: () => {
    const siteName = useSystemStore.getState().siteName;
    return { meta: [{ title: `Onboarding - ${siteName}` }] };
  },
});

function OnboardingPage() {
  const { data: setupStatus, isLoading } = useGetSetupStatus();

  if (isLoading) return null;

  if (!setupStatus?.data?.setupRequired) {
    return <Navigate to="/signin" replace />;
  }

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-background px-4 py-12">
      <BackgroundPattern />
      <div className="fade-in slide-in-from-bottom-4 relative z-10 w-full animate-in duration-700">
        <OnboardingForm cwd={setupStatus?.data?.cwd} />
      </div>
    </div>
  );
}
