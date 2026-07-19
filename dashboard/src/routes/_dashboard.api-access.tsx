import { createFileRoute } from '@tanstack/react-router';
import { ApiKeysList } from '#/features/instance/api-keys-list';

export const Route = createFileRoute('/_dashboard/api-access')({
  component: () => <ApiKeysList />,
});
