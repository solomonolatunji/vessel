import { createFileRoute } from '@tanstack/react-router';
import { StorageInstancesList } from '#/features/storage/storage-instances-list';

export const Route = createFileRoute('/_dashboard/storage')({
  component: () => <StorageInstancesList />,
});
