import type { QueryClient } from '@tanstack/react-query';
import { createRootRouteWithContext, HeadContent, Outlet } from '@tanstack/react-router';
import { ThemeProvider } from '#/components/theme-provider';
import { Toaster } from '#/components/ui/sonner';
import { TooltipProvider } from '#/components/ui/tooltip';
import { PostHogProvider } from '#/integrations/posthog-provider';

interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      {
        name: 'description',
        content: 'Codedock - Deploy apps instantly',
      },
      { title: 'Codedock' },
    ],
  }),
  component: RootDocument,
});

function RootDocument() {
  return (
    <ThemeProvider>
      <PostHogProvider>
        <HeadContent />
        <TooltipProvider>
          <Outlet />
        </TooltipProvider>
        <Toaster />
      </PostHogProvider>
    </ThemeProvider>
  );
}
