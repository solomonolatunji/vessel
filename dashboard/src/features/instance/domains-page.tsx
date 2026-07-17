import { Globe, Server } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Skeleton } from '#/components/ui/skeleton';
import { useGetSettings, useUpdateSettings } from '#/hooks/useSettings';

export const DomainsPage = () => {
  const { data, isLoading } = useGetSettings();
  const { mutateAsync: updateSettings, isPending } = useUpdateSettings();

  const [formData, setFormData] = useState({
    dashboardDomain: '',
    defaultWildcardDomain: '',
  });

  useEffect(() => {
    if (data?.data) {
      setFormData({
        dashboardDomain: data.data.dashboardDomain || '',
        defaultWildcardDomain: data.data.defaultWildcardDomain || '',
      });
    }
  }, [data?.data]);

  const handleSave = async (field: keyof typeof formData) => {
    try {
      await updateSettings({
        payload: {
          [field]: formData[field],
        } as any,
      });
      toast.success('Settings updated successfully');
    } catch {
      toast.error('Failed to update settings');
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

  return (
    <div className="space-y-6">
      <div className="flex flex-col justify-between gap-6 rounded-2xl border border-border/50 bg-card/40 p-8 md:flex-row md:items-start">
        <div className="flex-1 space-y-4">
          <div className="space-y-1">
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              INSTANCE DOMAINS
            </p>
            <h1 className="font-bold text-3xl tracking-tight">Domains</h1>
          </div>
          <p className="max-w-2xl text-muted-foreground text-sm leading-relaxed">
            Configure the dashboard domain and wildcard root domain for deployed services.
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
        {/* Dashboard Domain */}
        <div className="space-y-6 rounded-2xl border border-border/50 bg-card/40 p-8">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
              <Server className="h-5 w-5" />
            </div>
            <div>
              <h2 className="font-semibold text-lg">Dashboard Domain</h2>
              <p className="text-muted-foreground text-sm">
                No dashboard domain, point hostname to this server
              </p>
            </div>
          </div>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="font-bold text-muted-foreground text-xs uppercase tracking-wider">
                DASHBOARD DOMAIN
              </Label>
              <Input
                placeholder="vessl.example.com"
                value={formData.dashboardDomain}
                onChange={(e) => setFormData({ ...formData, dashboardDomain: e.target.value })}
                className="h-11 bg-background/50 font-mono"
              />
            </div>
            <Button
              onClick={() => handleSave('dashboardDomain')}
              disabled={isPending}
              className="h-11 w-full bg-primary font-bold text-primary-foreground text-xs uppercase tracking-wider"
            >
              Save Dashboard Domain
            </Button>
          </div>
        </div>

        {/* Root Domain */}
        <div className="space-y-6 rounded-2xl border border-border/50 bg-card/40 p-8">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-blue-500/20 bg-blue-500/10 text-blue-500">
              <Globe className="h-5 w-5" />
            </div>
            <div>
              <h2 className="font-semibold text-lg">Set root domain</h2>
              <p className="text-muted-foreground text-sm">
                Used as default wildcard for user deployments
              </p>
            </div>
          </div>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="font-bold text-muted-foreground text-xs uppercase tracking-wider">
                WILDCARD ROOT DOMAIN
              </Label>
              <Input
                placeholder="apps.example.com"
                value={formData.defaultWildcardDomain}
                onChange={(e) =>
                  setFormData({ ...formData, defaultWildcardDomain: e.target.value })
                }
                className="h-11 bg-background/50 font-mono"
              />
            </div>
            <Button
              onClick={() => handleSave('defaultWildcardDomain')}
              disabled={isPending}
              className="h-11 w-full font-bold text-xs uppercase tracking-wider"
            >
              Save Root Domain
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};
