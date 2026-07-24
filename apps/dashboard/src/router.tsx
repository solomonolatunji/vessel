import { createRouter as createTanStackRouter, RouterProvider } from '@tanstack/react-router';
import TanstackQueryProvider, { getContext } from './integrations/tanstack-query/root-provider';
import { routeTree } from './routeTree.gen';

const context = getContext();

export const router = createTanStackRouter({
  routeTree,
  context,
  scrollRestoration: true,
  defaultPreload: 'intent',
  defaultPreloadStaleTime: 0,
  defaultPreloadGcTime: 10_000,
});

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

export function AppRouter() {
  return (
    <TanstackQueryProvider>
      <RouterProvider router={router} />
    </TanstackQueryProvider>
  );
}
