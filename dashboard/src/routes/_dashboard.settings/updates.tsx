import { createFileRoute } from '@tanstack/react-router';
import { UpdatesPage } from '#/features/instance/update-settings';

export const Route = createFileRoute('/_dashboard/settings/updates')({
  component: UpdatesPage,
});
