import { createFileRoute } from '@tanstack/react-router';
import { DnsSettings } from '#/features/instance/dns-settings';

export const Route = createFileRoute('/_dashboard/settings/dns')({
  component: DnsSettings,
});
