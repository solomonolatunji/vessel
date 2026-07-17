import { createFileRoute } from '@tanstack/react-router';
import { DomainsPage } from '#/features/instance/domains-page';

export const Route = createFileRoute('/_dashboard/domains')({
  component: DomainsPage,
});
