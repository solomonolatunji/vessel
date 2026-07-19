import { Bell, Check, Mail, MessageSquare, Phone, Send, Webhook, Zap } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Skeleton } from '#/components/ui/skeleton';
import { Switch } from '#/components/ui/switch';
import {
  useGetNotificationSettings,
  useTestNotification,
  useUpdateNotificationSettings,
} from '#/hooks/useSettings';
import type { ServerSettings } from '#/interfaces/settings';

type NotifSettings = Pick<
  ServerSettings,
  | 'discordWebhookUrl'
  | 'discordPingEnabled'
  | 'discordEnabled'
  | 'slackWebhookUrl'
  | 'slackEnabled'
  | 'telegramBotToken'
  | 'telegramChatId'
  | 'telegramEnabled'
  | 'smtpHost'
  | 'smtpPort'
  | 'smtpUser'
  | 'smtpPassword'
  | 'smtpFromName'
  | 'smtpFromAddress'
  | 'smtpEnabled'
  | 'resendApiKey'
  | 'resendEnabled'
  | 'pushoverUserKey'
  | 'pushoverApiToken'
  | 'pushoverEnabled'
  | 'genericWebhookUrl'
  | 'genericWebhookEnabled'
  | 'notificationAlerts'
>;

const EMPTY: NotifSettings = {
  discordWebhookUrl: '',
  discordPingEnabled: false,
  discordEnabled: false,
  slackWebhookUrl: '',
  slackEnabled: false,
  telegramBotToken: '',
  telegramChatId: '',
  telegramEnabled: false,
  smtpHost: '',
  smtpPort: 587,
  smtpUser: '',
  smtpPassword: '',
  smtpFromName: '',
  smtpFromAddress: '',
  smtpEnabled: false,
  resendApiKey: '',
  resendEnabled: false,
  pushoverUserKey: '',
  pushoverApiToken: '',
  pushoverEnabled: false,
  genericWebhookUrl: '',
  genericWebhookEnabled: false,
  notificationAlerts: true,
};

type SectionProps = {
  icon: React.ReactNode;
  title: string;
  provider: string;
  enabled: boolean;
  onToggle: (v: boolean) => void;
  onTest: () => void;
  testing: boolean;
  children: React.ReactNode;
};

const Section = ({
  icon,
  title,
  provider: _provider,
  enabled,
  onToggle,
  onTest,
  testing,
  children,
}: SectionProps) => (
  <div className="rounded-xl border border-border/50 bg-card/40 p-6">
    <div className="mb-4 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
          {icon}
        </div>
        <span className="font-semibold text-sm">{title}</span>
      </div>
      <div className="flex items-center gap-3">
        <Button
          size="sm"
          variant="outline"
          className="text-xs"
          disabled={!enabled || testing}
          onClick={onTest}
        >
          <Zap className="mr-1.5 h-3.5 w-3.5" />
          {testing ? 'Sending…' : 'Test'}
        </Button>
        <Switch checked={enabled} onCheckedChange={onToggle} />
      </div>
    </div>
    <div className={`space-y-4 ${!enabled ? 'pointer-events-none opacity-40' : ''}`}>
      {children}
    </div>
  </div>
);

export const NotificationsSettings = () => {
  const { data, isLoading } = useGetNotificationSettings();
  const { mutateAsync: update, isPending } = useUpdateNotificationSettings();
  const { mutateAsync: testNotif, isPending: testing } = useTestNotification();

  const [form, setForm] = useState<NotifSettings>(EMPTY);
  const [testingProvider, setTestingProvider] = useState<string | null>(null);

  useEffect(() => {
    if (data?.data) {
      const s = data.data as Record<string, unknown>;
      setForm({
        discordWebhookUrl: (s.discordWebhookUrl as string) ?? '',
        discordPingEnabled: (s.discordPingEnabled as boolean) ?? false,
        discordEnabled: (s.discordEnabled as boolean) ?? false,
        slackWebhookUrl: (s.slackWebhookUrl as string) ?? '',
        slackEnabled: (s.slackEnabled as boolean) ?? false,
        telegramBotToken: (s.telegramBotToken as string) ?? '',
        telegramChatId: (s.telegramChatId as string) ?? '',
        telegramEnabled: (s.telegramEnabled as boolean) ?? false,
        smtpHost: (s.smtpHost as string) ?? '',
        smtpPort: (s.smtpPort as number) ?? 587,
        smtpUser: (s.smtpUser as string) ?? '',
        smtpPassword: (s.smtpPassword as string) ?? '',
        smtpFromName: (s.smtpFromName as string) ?? '',
        smtpFromAddress: (s.smtpFromAddress as string) ?? '',
        smtpEnabled: (s.smtpEnabled as boolean) ?? false,
        resendApiKey: (s.resendApiKey as string) ?? '',
        resendEnabled: (s.resendEnabled as boolean) ?? false,
        pushoverUserKey: (s.pushoverUserKey as string) ?? '',
        pushoverApiToken: (s.pushoverApiToken as string) ?? '',
        pushoverEnabled: (s.pushoverEnabled as boolean) ?? false,
        genericWebhookUrl: (s.genericWebhookUrl as string) ?? '',
        genericWebhookEnabled: (s.genericWebhookEnabled as boolean) ?? false,
        notificationAlerts: (s.notificationAlerts as boolean) ?? true,
      });
    }
  }, [data]);

  const set = (k: keyof NotifSettings, v: unknown) => setForm((f) => ({ ...f, [k]: v }));

  const handleSave = async () => {
    try {
      await update(form as Record<string, unknown>);
      toast.success('Notification settings saved');
    } catch {
      toast.error('Failed to save notification settings');
    }
  };

  const handleTest = async (provider: string) => {
    setTestingProvider(provider);
    try {
      await testNotif({ provider });
      toast.success(`Test notification sent via ${provider}`);
    } catch {
      toast.error(`Failed to send test via ${provider}`);
    } finally {
      setTestingProvider(null);
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[...Array(4)].map((_, i) => (
          <Skeleton key={i} className="h-40 w-full rounded-xl" />
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-lg">Notification Channels</h2>
          <p className="text-muted-foreground text-sm">
            Configure where Vessl sends alerts for deployments, errors, and system events.
          </p>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground text-sm">Global alerts</span>
            <Switch
              checked={form.notificationAlerts}
              onCheckedChange={(v) => set('notificationAlerts', v)}
            />
          </div>
          <Button onClick={handleSave} disabled={isPending}>
            <Check className="mr-2 h-4 w-4" />
            {isPending ? 'Saving…' : 'Save Changes'}
          </Button>
        </div>
      </div>

      <Section
        icon={<MessageSquare className="h-4 w-4" />}
        title="Discord"
        provider="discord"
        enabled={form.discordEnabled ?? false}
        onToggle={(v) => set('discordEnabled', v)}
        onTest={() => handleTest('discord')}
        testing={testingProvider === 'discord' && testing}
      >
        <div className="space-y-2">
          <Label className="text-xs">Webhook URL</Label>
          <Input
            value={form.discordWebhookUrl ?? ''}
            onChange={(e) => set('discordWebhookUrl', e.target.value)}
            placeholder="https://discord.com/api/webhooks/..."
            className="font-mono text-xs"
          />
        </div>
        <div className="flex items-center gap-2">
          <Switch
            checked={form.discordPingEnabled ?? false}
            onCheckedChange={(v) => set('discordPingEnabled', v)}
          />
          <Label className="text-muted-foreground text-xs">@here ping on critical alerts</Label>
        </div>
      </Section>

      <Section
        icon={<Bell className="h-4 w-4" />}
        title="Slack"
        provider="slack"
        enabled={form.slackEnabled ?? false}
        onToggle={(v) => set('slackEnabled', v)}
        onTest={() => handleTest('slack')}
        testing={testingProvider === 'slack' && testing}
      >
        <div className="space-y-2">
          <Label className="text-xs">Webhook URL</Label>
          <Input
            value={form.slackWebhookUrl ?? ''}
            onChange={(e) => set('slackWebhookUrl', e.target.value)}
            placeholder="https://hooks.slack.com/services/..."
            className="font-mono text-xs"
          />
        </div>
      </Section>

      <Section
        icon={<Send className="h-4 w-4" />}
        title="Telegram"
        provider="telegram"
        enabled={form.telegramEnabled ?? false}
        onToggle={(v) => set('telegramEnabled', v)}
        onTest={() => handleTest('telegram')}
        testing={testingProvider === 'telegram' && testing}
      >
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-xs">Bot Token</Label>
            <Input
              type="password"
              value={form.telegramBotToken ?? ''}
              onChange={(e) => set('telegramBotToken', e.target.value)}
              placeholder="1234567890:AAF..."
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">Chat ID</Label>
            <Input
              value={form.telegramChatId ?? ''}
              onChange={(e) => set('telegramChatId', e.target.value)}
              placeholder="-100..."
              className="font-mono text-xs"
            />
          </div>
        </div>
      </Section>

      <Section
        icon={<Mail className="h-4 w-4" />}
        title="SMTP Email"
        provider="smtp"
        enabled={form.smtpEnabled ?? false}
        onToggle={(v) => set('smtpEnabled', v)}
        onTest={() => handleTest('smtp')}
        testing={testingProvider === 'smtp' && testing}
      >
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-xs">SMTP Host</Label>
            <Input
              value={form.smtpHost ?? ''}
              onChange={(e) => set('smtpHost', e.target.value)}
              placeholder="smtp.example.com"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">Port</Label>
            <Input
              type="number"
              value={form.smtpPort ?? 587}
              onChange={(e) => set('smtpPort', Number(e.target.value))}
              placeholder="587"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">Username</Label>
            <Input
              value={form.smtpUser ?? ''}
              onChange={(e) => set('smtpUser', e.target.value)}
              placeholder="user@example.com"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">Password</Label>
            <Input
              type="password"
              value={form.smtpPassword ?? ''}
              onChange={(e) => set('smtpPassword', e.target.value)}
              placeholder="••••••••"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">From Name</Label>
            <Input
              value={form.smtpFromName ?? ''}
              onChange={(e) => set('smtpFromName', e.target.value)}
              placeholder="Vessl"
              className="text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">From Address</Label>
            <Input
              value={form.smtpFromAddress ?? ''}
              onChange={(e) => set('smtpFromAddress', e.target.value)}
              placeholder="noreply@example.com"
              className="font-mono text-xs"
            />
          </div>
        </div>
      </Section>

      <Section
        icon={<Mail className="h-4 w-4" />}
        title="Resend"
        provider="resend"
        enabled={form.resendEnabled ?? false}
        onToggle={(v) => set('resendEnabled', v)}
        onTest={() => handleTest('resend')}
        testing={testingProvider === 'resend' && testing}
      >
        <div className="space-y-2">
          <Label className="text-xs">API Key</Label>
          <Input
            type="password"
            value={form.resendApiKey ?? ''}
            onChange={(e) => set('resendApiKey', e.target.value)}
            placeholder="re_..."
            className="font-mono text-xs"
          />
        </div>
      </Section>

      <Section
        icon={<Phone className="h-4 w-4" />}
        title="Pushover"
        provider="pushover"
        enabled={form.pushoverEnabled ?? false}
        onToggle={(v) => set('pushoverEnabled', v)}
        onTest={() => handleTest('pushover')}
        testing={testingProvider === 'pushover' && testing}
      >
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-xs">User Key</Label>
            <Input
              value={form.pushoverUserKey ?? ''}
              onChange={(e) => set('pushoverUserKey', e.target.value)}
              placeholder="u..."
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">API Token</Label>
            <Input
              type="password"
              value={form.pushoverApiToken ?? ''}
              onChange={(e) => set('pushoverApiToken', e.target.value)}
              placeholder="a..."
              className="font-mono text-xs"
            />
          </div>
        </div>
      </Section>

      <Section
        icon={<Webhook className="h-4 w-4" />}
        title="Generic Webhook"
        provider="webhook"
        enabled={form.genericWebhookEnabled ?? false}
        onToggle={(v) => set('genericWebhookEnabled', v)}
        onTest={() => handleTest('webhook')}
        testing={testingProvider === 'webhook' && testing}
      >
        <div className="space-y-2">
          <Label className="text-xs">Webhook URL</Label>
          <Input
            value={form.genericWebhookUrl ?? ''}
            onChange={(e) => set('genericWebhookUrl', e.target.value)}
            placeholder="https://example.com/webhook"
            className="font-mono text-xs"
          />
        </div>
      </Section>
    </div>
  );
};
