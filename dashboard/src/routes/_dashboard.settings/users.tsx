import { createFileRoute } from '@tanstack/react-router';
import { UsersPage } from '#/features/instance/users-page';

export const Route = createFileRoute('/_dashboard/settings/users')({
  component: UsersPage,
});
