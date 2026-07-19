import { createFileRoute } from '@tanstack/react-router';
import { GithubIntegration } from '#/features/sources';

export const Route = createFileRoute('/_dashboard/sources')({
  component: () => <GithubIntegration />,
});
