import { createFileRoute } from '@tanstack/react-router';
import { GithubIntegration } from '#/features/sources';

export const Route = createFileRoute('/_dashboard/sources')({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      code: (search.code as string) || undefined,
    };
  },
  component: () => <GithubIntegration />,
});
