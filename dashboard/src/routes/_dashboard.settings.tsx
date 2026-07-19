import { createFileRoute } from '@tanstack/react-router';
import { SettingsLayout } from '#/features/instance/settings-page';

export const Route = createFileRoute('/_dashboard/settings')({
  component: SettingsLayout,
});
