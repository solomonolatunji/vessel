import { createFileRoute } from '@tanstack/react-router';
import { MaintenancePage } from '#/features/instance/maintenance-settings';

export const Route = createFileRoute('/_dashboard/settings/maintenance')({
  component: MaintenancePage,
});
