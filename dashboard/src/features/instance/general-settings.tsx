import { Check, Clock, Cpu, Globe, Info, Lock } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Skeleton } from '#/components/ui/skeleton';
import { Switch } from '#/components/ui/switch';
import { useGetSettings, useUpdateSettings } from '#/hooks/useSettings';
import type { ServerSettings } from '#/interfaces/settings';

type GeneralFields = Pick<
  ServerSettings,
  | 'siteName'
  | 'dashboardDomain'
  | 'defaultWildcardDomain'
  | 'publicIpv4'
  | 'publicIpv6'
  | 'traefikWildcardIp'
  | 'registrationEnabled'
  | 'registrationDomainAllowlist'
  | 'ipAllowlist'
  | 'mcpServerEnabled'
  | 'disableTwoStepConfirmation'
  | 'telemetryEnabled'
  | 'concurrentBuilds'
  | 'deploymentTimeout'
  | 'serverTimezone'
>;

const EMPTY: GeneralFields = {
  siteName: '',
  dashboardDomain: '',
  defaultWildcardDomain: '',
  publicIpv4: '',
  publicIpv6: '',
  traefikWildcardIp: '',
  registrationEnabled: true,
  registrationDomainAllowlist: '',
  ipAllowlist: '',
  mcpServerEnabled: false,
  disableTwoStepConfirmation: false,
  telemetryEnabled: true,
  concurrentBuilds: 2,
  deploymentTimeout: 300,
  serverTimezone: 'UTC',
};

type RowProps = { label: string; description?: string; children: React.ReactNode };
const Row = ({ label, description, children }: RowProps) => (
  <div className="flex items-start justify-between gap-8 border-border/30 border-b py-4 last:border-0">
    <div className="min-w-0 flex-1">
      <Label className="font-medium text-sm">{label}</Label>
      {description && <p className="mt-0.5 text-muted-foreground text-xs">{description}</p>}
    </div>
    <div className="w-80 shrink-0">{children}</div>
  </div>
);

type SectionProps = { icon: React.ReactNode; title: string; children: React.ReactNode };
const Section = ({ icon, title, children }: SectionProps) => (
  <div className="rounded-xl border border-border/50 bg-card/40 p-6">
    <div className="mb-4 flex items-center gap-3">
      <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
        {icon}
      </div>
      <span className="font-semibold text-sm">{title}</span>
    </div>
    {children}
  </div>
);

export const GeneralSettings = () => {
  const { data, isLoading } = useGetSettings();
  const { mutateAsync: updateSettings, isPending } = useUpdateSettings();

  const [form, setForm] = useState<GeneralFields>(EMPTY);

  useEffect(() => {
    if (data?.data) {
      const s = data.data;
      setForm({
        siteName: s.siteName ?? '',
        dashboardDomain: s.dashboardDomain ?? '',
        defaultWildcardDomain: s.defaultWildcardDomain ?? '',
        publicIpv4: s.publicIpv4 ?? '',
        publicIpv6: s.publicIpv6 ?? '',
        traefikWildcardIp: s.traefikWildcardIp ?? '',
        registrationEnabled: s.registrationEnabled,
        registrationDomainAllowlist: s.registrationDomainAllowlist ?? '',
        ipAllowlist: s.ipAllowlist ?? '',
        mcpServerEnabled: s.mcpServerEnabled,
        disableTwoStepConfirmation: s.disableTwoStepConfirmation,
        telemetryEnabled: s.telemetryEnabled,
        concurrentBuilds: s.concurrentBuilds,
        deploymentTimeout: s.deploymentTimeout,
        serverTimezone: s.serverTimezone,
      });
    }
  }, [data]);

  const set = <K extends keyof GeneralFields>(k: K, v: GeneralFields[K]) =>
    setForm((f) => ({ ...f, [k]: v }));

  const handleSave = async () => {
    try {
      await updateSettings({ payload: form });
      toast.success('Settings saved');
    } catch {
      toast.error('Failed to save settings');
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[...Array(3)].map((_, i) => (
          <Skeleton key={i} className="h-48 w-full rounded-xl" />
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-lg">General</h2>
          <p className="text-muted-foreground text-sm">
            Core instance configuration and network settings.
          </p>
        </div>
        <Button onClick={handleSave} disabled={isPending}>
          <Check className="mr-2 h-4 w-4" />
          {isPending ? 'Saving…' : 'Save Changes'}
        </Button>
      </div>

      <Section icon={<Globe className="h-4 w-4" />} title="Instance Identity">
        <Row label="Site Name" description="Displayed in the browser tab and emails.">
          <Input
            value={form.siteName ?? ''}
            onChange={(e) => set('siteName', e.target.value)}
            placeholder="Vessl"
            className="text-sm"
          />
        </Row>
        <Row label="Dashboard Domain" description="The domain Vessl control panel is served from.">
          <Input
            value={form.dashboardDomain ?? ''}
            onChange={(e) => set('dashboardDomain', e.target.value)}
            placeholder="pilot.example.com"
            className="font-mono text-xs"
          />
        </Row>
        <Row
          label="Wildcard Domain"
          description="Root domain for generated app URLs (e.g. *.apps.example.com)."
        >
          <Input
            value={form.defaultWildcardDomain ?? ''}
            onChange={(e) => set('defaultWildcardDomain', e.target.value)}
            placeholder="apps.example.com"
            className="font-mono text-xs"
          />
        </Row>
      </Section>

      <Section icon={<Info className="h-4 w-4" />} title="Network">
        <Row label="Public IPv4" description="Server's public IPv4 address for DNS A records.">
          <Input
            value={form.publicIpv4 ?? ''}
            onChange={(e) => set('publicIpv4', e.target.value)}
            placeholder="1.2.3.4"
            className="font-mono text-xs"
          />
        </Row>
        <Row label="Public IPv6" description="Server's public IPv6 address (optional).">
          <Input
            value={form.publicIpv6 ?? ''}
            onChange={(e) => set('publicIpv6', e.target.value)}
            placeholder="2001:db8::1"
            className="font-mono text-xs"
          />
        </Row>
        <Row label="Traefik Wildcard IP" description="IP Traefik routes wildcard domains to.">
          <Input
            value={form.traefikWildcardIp ?? ''}
            onChange={(e) => set('traefikWildcardIp', e.target.value)}
            placeholder="1.2.3.4"
            className="font-mono text-xs"
          />
        </Row>
        <Row
          label="IP Allowlist"
          description="Comma-separated CIDRs that can access the control plane."
        >
          <Input
            value={form.ipAllowlist ?? ''}
            onChange={(e) => set('ipAllowlist', e.target.value)}
            placeholder="0.0.0.0/0"
            className="font-mono text-xs"
          />
        </Row>
      </Section>

      <Section icon={<Lock className="h-4 w-4" />} title="Access & Security">
        <Row
          label="Open Registration"
          description="Allow new users to register without an invitation."
        >
          <Switch
            checked={form.registrationEnabled}
            onCheckedChange={(v) => set('registrationEnabled', v)}
          />
        </Row>
        <Row
          label="Domain Allowlist"
          description="Restrict registration to specific email domains (comma-separated)."
        >
          <Input
            value={form.registrationDomainAllowlist ?? ''}
            onChange={(e) => set('registrationDomainAllowlist', e.target.value)}
            placeholder="example.com, company.org"
            className="text-xs"
          />
        </Row>
        <Row
          label="Disable Two-Step Confirmation"
          description="Skip destructive action confirmation dialogs (not recommended)."
        >
          <Switch
            checked={form.disableTwoStepConfirmation}
            onCheckedChange={(v) => set('disableTwoStepConfirmation', v)}
          />
        </Row>
        <Row
          label="MCP Server"
          description="Enable the Model Context Protocol server for AI tooling."
        >
          <Switch
            checked={form.mcpServerEnabled}
            onCheckedChange={(v) => set('mcpServerEnabled', v)}
          />
        </Row>
      </Section>

      <Section icon={<Cpu className="h-4 w-4" />} title="Build & Deployment">
        <Row label="Concurrent Builds" description="Max number of parallel build jobs.">
          <Input
            type="number"
            value={form.concurrentBuilds}
            onChange={(e) => set('concurrentBuilds', Number(e.target.value))}
            min={1}
            max={20}
            className="font-mono text-xs"
          />
        </Row>
        <Row
          label="Deployment Timeout (s)"
          description="Seconds before a deployment is considered failed."
        >
          <Input
            type="number"
            value={form.deploymentTimeout}
            onChange={(e) => set('deploymentTimeout', Number(e.target.value))}
            min={60}
            className="font-mono text-xs"
          />
        </Row>
      </Section>

      <Section icon={<Clock className="h-4 w-4" />} title="System">
        <Row label="Server Timezone" description="Timezone used for cron jobs and logs.">
          <Input
            value={form.serverTimezone ?? ''}
            onChange={(e) => set('serverTimezone', e.target.value)}
            placeholder="UTC"
            className="font-mono text-xs"
          />
        </Row>
        <Row label="Telemetry" description="Send anonymous usage statistics to help improve Vessl.">
          <Switch
            checked={form.telemetryEnabled}
            onCheckedChange={(v) => set('telemetryEnabled', v)}
          />
        </Row>
      </Section>
    </div>
  );
};
