import { Check } from 'lucide-react';
import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Skeleton } from '#/components/ui/skeleton';
import { Switch } from '#/components/ui/switch';
import { useListOAuthProviders, useSaveOAuthProvider } from '#/hooks/useOAuth';
import type { SaveOAuthProviderRequest } from '#/interfaces/oauth';

type ProviderConfig = {
  id: string;
  name: string;
  fields: {
    key: keyof SaveOAuthProviderRequest;
    label: string;
    placeholder: string;
    type?: string;
  }[];
};

const PROVIDERS: ProviderConfig[] = [
  {
    id: 'github',
    name: 'GitHub',
    fields: [
      { key: 'clientId', label: 'Client ID', placeholder: 'Iv1.abc123...' },
      { key: 'clientSecret', label: 'Client Secret', placeholder: '••••••••', type: 'password' },
    ],
  },
  {
    id: 'gitlab',
    name: 'GitLab',
    fields: [
      { key: 'clientId', label: 'Application ID', placeholder: 'abc123...' },
      { key: 'clientSecret', label: 'Secret', placeholder: '••••••••', type: 'password' },
      {
        key: 'baseUrl',
        label: 'Self-hosted URL (optional)',
        placeholder: 'https://gitlab.example.com',
      },
    ],
  },
  {
    id: 'google',
    name: 'Google',
    fields: [
      { key: 'clientId', label: 'Client ID', placeholder: '123456789.apps.googleusercontent.com' },
      { key: 'clientSecret', label: 'Client Secret', placeholder: 'GOCSPX-...', type: 'password' },
    ],
  },
  {
    id: 'microsoft',
    name: 'Microsoft',
    fields: [
      {
        key: 'clientId',
        label: 'Application (Client) ID',
        placeholder: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx',
      },
      { key: 'clientSecret', label: 'Client Secret', placeholder: '••••••••', type: 'password' },
      { key: 'tenant', label: 'Tenant ID', placeholder: 'common' },
    ],
  },
];

type ProviderState = Record<string, Partial<SaveOAuthProviderRequest>>;

export const OAuthProvidersList = () => {
  const { data, isLoading } = useListOAuthProviders();
  const { mutateAsync: save, isPending } = useSaveOAuthProvider();

  const [form, setForm] = useState<ProviderState>({});
  const [saving, setSaving] = useState<string | null>(null);

  const isInitialized = React.useRef(false);

  useEffect(() => {
    if (data?.data) {
      if (!isInitialized.current) {
        isInitialized.current = true;
        const initial: ProviderState = {};
        for (const p of PROVIDERS) {
          const existing = data.data.find((d) => d.providerName === p.id);
          initial[p.id] = {
            id: existing?.id,
            providerName: p.id,
            clientId: existing?.clientId ?? '',
            clientSecret: existing?.clientSecret ? '********' : '',
            baseUrl: existing?.baseUrl ?? '',
            tenant: existing?.tenant ?? '',
            enabled: existing?.enabled ?? false,
          };
        }
        setForm(initial);
      } else {
        setForm((f) => {
          const next = { ...f };
          let changed = false;
          for (const p of PROVIDERS) {
            const existing = data.data.find((d) => d.providerName === p.id);
            if (existing?.id && next[p.id]?.id !== existing.id) {
              next[p.id] = { ...next[p.id], id: existing.id };
              changed = true;
            }
          }
          return changed ? next : f;
        });
      }
    }
  }, [data]);

  const set = (providerId: string, key: keyof SaveOAuthProviderRequest, value: unknown) => {
    setForm((f) => ({ ...f, [providerId]: { ...f[providerId], [key]: value } }));
  };

  const handleSaveAll = async () => {
    setSaving('all');
    try {
      await Promise.all(
        PROVIDERS.map((p) => {
          const s = form[p.id] as SaveOAuthProviderRequest;
          const payload = {
            ...s,
            redirectUri: `${window.location.origin}/api/auth/oauth/${p.id}/callback`,
            clientSecret: s.clientSecret === '********' ? undefined : s.clientSecret,
            enabled:
              !s.clientId || (!s.clientSecret && s.clientSecret !== '********') ? false : s.enabled,
          };
          return save({ payload });
        })
      );
      toast.success('OAuth providers saved');
    } catch {
      toast.error('Failed to save OAuth providers');
    } finally {
      setSaving(null);
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[...Array(4)].map((_, i) => (
          <Skeleton key={i} className="h-32 w-full rounded-xl" />
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-lg">OAuth Providers</h2>
          <p className="text-muted-foreground text-sm">
            Enable OAuth login for users. The redirect URI is{' '}
            <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">
              {window.location.origin}/api/auth/oauth/[provider]/callback
            </code>
          </p>
        </div>
        <Button onClick={handleSaveAll} disabled={saving === 'all' || isPending}>
          <Check className="mr-2 h-4 w-4" />
          {saving === 'all' ? 'Saving…' : 'Save Changes'}
        </Button>
      </div>

      {PROVIDERS.map((provider) => {
        const state = form[provider.id] ?? {};
        return (
          <div key={provider.id} className="rounded-xl border border-border/50 bg-card/40 p-6">
            <div className="mb-5 flex items-center justify-between">
              <span className="font-semibold text-sm">{provider.name}</span>
              <div className="flex items-center gap-2.5">
                <span className="text-muted-foreground text-xs">
                  {state.enabled ? 'Enabled' : 'Disabled'}
                </span>
                <Switch
                  checked={state.enabled ?? false}
                  onCheckedChange={(v) => set(provider.id, 'enabled', v)}
                />
              </div>
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              {provider.fields.map((f) => (
                <div key={f.key} className="space-y-1.5">
                  <Label className="text-xs">{f.label}</Label>
                  <Input
                    type={f.type ?? 'text'}
                    value={(state[f.key] as string) ?? ''}
                    onChange={(e) => set(provider.id, f.key, e.target.value)}
                    placeholder={f.placeholder}
                    className="font-mono text-xs"
                  />
                </div>
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
};
