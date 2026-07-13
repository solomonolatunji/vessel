import { createFileRoute, redirect } from '@tanstack/react-router';
import { OnboardingPage } from '#/features/onboarding/onboarding-page';
import { authStore } from '#/stores/authStore';

export const Route = createFileRoute('/onboarding')({
  beforeLoad: () => {
    if (!authStore.state.isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: OnboardingPage,
});
