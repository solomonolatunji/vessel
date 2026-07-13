import { Layers01Icon } from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';

import { authStore } from '#/stores/authStore';

export const Route = createFileRoute('/_auth')({
  beforeLoad: () => {
    if (authStore.state.isAuthenticated) {
      throw redirect({ to: '/' });
    }
  },
  component: AuthLayout,
});

function AuthLayout() {
  return (
    <div className="flex min-h-screen w-full flex-col items-center justify-center bg-background font-sans text-foreground p-4">
      {/* Form Container */}
      <div className="w-full max-w-lg sm:max-w-120 relative z-10 flex flex-col items-center">
        {/* Top Logo */}
        <div className="flex items-center justify-center gap-2 font-bold text-xl mb-8">
          <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-primary text-primary-foreground">
            <HugeiconsIcon icon={Layers01Icon} className="h-5 w-5" />
          </div>
          Vessl
        </div>

        <div className="w-full">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
