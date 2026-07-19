import { Globe } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { useGetSettings, useUpdateSettings } from '#/hooks/useSettings';
import { Skeleton } from '#/components/ui/skeleton';
import { DnsProviderForm } from './components/dns-provider-form';

const providers = [
  { id: 'cloudflare', name: 'Cloudflare', sub: 'API KEY / TOKEN + ZONE' },
  { id: 'namecheap', name: 'Namecheap', sub: 'API USER + KEY' },
  { id: 'spaceship', name: 'Spaceship', sub: 'API KEY + SECRET' },
];

export const DnsSettings = () => {
  const { data, isLoading } = useGetSettings();
  const { mutateAsync: updateSettings, isPending } = useUpdateSettings();

  const [activeProvider, setActiveProvider] = useState('cloudflare');
  const [formData, setFormData] = useState<Record<string, string>>({
    cloudflareApiToken: '',
    namecheapApiUser: '',
    namecheapApiKey: '',
    namecheapClientIp: '',
    spaceshipApiKey: '',
    cloudflareEmail: '',
    cloudflareZoneId: '',
    spaceshipApiSecret: '',
  });

  useEffect(() => {
    if (data?.data) {
      setFormData((prev) => ({
        ...prev,
        cloudflareApiToken: data.data.cloudflareApiToken || '',
        namecheapApiUser: data.data.namecheapApiUser || '',
        namecheapApiKey: data.data.namecheapApiKey || '',
        namecheapClientIp: data.data.namecheapClientIp || '',
        spaceshipApiKey: data.data.spaceshipApiKey || '',
      }));
    }
  }, [data?.data]);

  const handleSaveProvider = async (provider: string) => {
    try {
      const payload: Record<string, string> = {};
      if (provider === 'cloudflare') {
        payload.cloudflareApiToken = formData.cloudflareApiToken;
      } else if (provider === 'namecheap') {
        payload.namecheapApiUser = formData.namecheapApiUser;
        payload.namecheapApiKey = formData.namecheapApiKey;
        payload.namecheapClientIp = formData.namecheapClientIp;
      } else if (provider === 'spaceship') {
        payload.spaceshipApiKey = formData.spaceshipApiKey;
      }

      await updateSettings({ payload });
      toast.success('Provider credentials saved');
    } catch {
      toast.error('Failed to save provider credentials');
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-[200px] w-full rounded-2xl" />
        <Skeleton className="h-[300px] w-full rounded-2xl" />
      </div>
    );
  }

  const activeProviderData = providers.find((p) => p.id === activeProvider);

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Globe className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">API credentials</h1>
            <p className="text-muted-foreground text-sm">
              Vessl manages your domains using various DNS providers. You can add your API
              credentials here.
            </p>
          </div>
        </div>
        <div className="flex shrink-0 flex-col items-end gap-4">
          <div className="flex items-center gap-3">
            <img
              src={`/dns-providers/${activeProvider}.svg`}
              alt={activeProviderData?.name}
              className="h-4 w-auto"
            />
            <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              {activeProviderData?.name}
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 pt-2 md:grid-cols-3">
        {providers.map((p) => {
          const isActive = activeProvider === p.id;
          return (
            <button
              type="button"
              key={p.id}
              onClick={() => setActiveProvider(p.id)}
              className={`group relative w-full cursor-pointer rounded-2xl border p-6 text-left transition-all duration-200 ${
                isActive
                  ? 'border-primary/50 bg-card/40 shadow-sm'
                  : 'border-border/50 bg-background/50 hover:border-border/80 hover:bg-card/40'
              }`}
            >
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                  <div
                    className={`flex h-10 w-10 items-center justify-center rounded-xl border transition-colors ${
                      isActive ? 'border-primary/30 bg-primary/5' : 'border-border/50 bg-background'
                    }`}
                  >
                    <img src={`/dns-providers/${p.id}.svg`} alt={p.name} className="h-5 w-auto" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-[15px]">{p.name}</h3>
                    <p className="mt-1 max-w-[120px] text-[9px] text-muted-foreground/70 uppercase leading-relaxed tracking-[0.15em]">
                      {p.sub}
                    </p>
                  </div>
                </div>
                <span
                  className={`rounded-md border px-2 py-0.5 font-bold text-[8px] uppercase tracking-wider ${
                    isActive
                      ? 'border-primary/30 text-primary'
                      : 'border-border/50 text-muted-foreground/50'
                  }`}
                >
                  NEW
                </span>
              </div>
            </button>
          );
        })}
      </div>

      <div className="mt-4 space-y-10 rounded-2xl border border-border/50 bg-card/40 p-6">
        <div className="flex items-center gap-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-2xl border border-border/50 bg-background/50">
            <img
              src={`/dns-providers/${activeProvider}.svg`}
              alt={activeProvider}
              className="h-6 w-auto"
            />
          </div>
          <div>
            <h2 className="font-bold text-2xl tracking-tight">
              Connect {activeProviderData?.name}
            </h2>
            <p className="mt-1.5 font-medium text-muted-foreground text-sm">
              Store provider credentials for DNS record automation.
            </p>
          </div>
        </div>

        <DnsProviderForm
          activeProvider={activeProvider}
          formData={formData}
          setFormData={setFormData}
          isPending={isPending}
          handleSaveProvider={handleSaveProvider}
        />
      </div>
    </div>
  );
};
