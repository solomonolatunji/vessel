import { Check, Globe } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Skeleton } from '#/components/ui/skeleton';
import { useGetSettings, useUpdateSettings } from '#/hooks/useSettings';

type SectionProps = {
  icon: React.ReactNode;
  title: string;
  action?: React.ReactNode;
  children: React.ReactNode;
};

const Section = ({ icon, title, action, children }: SectionProps) => (
  <div className="rounded-xl border border-border/50 bg-card/40 p-6">
    <div className="mb-4 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
          {icon}
        </div>
        <span className="font-semibold text-sm">{title}</span>
      </div>
      {action && <div className="flex shrink-0">{action}</div>}
    </div>
    <div className="divide-y divide-border/50">{children}</div>
  </div>
);

type RowProps = { label: string; description?: string; children: React.ReactNode };
const Row = ({ label, description, children }: RowProps) => (
  <div className="flex flex-col gap-4 py-4 md:flex-row md:items-center md:justify-between">
    <div className="flex-1 pr-4">
      <Label className="font-medium text-sm">{label}</Label>
      {description && <p className="mt-1 text-muted-foreground text-sm">{description}</p>}
    </div>
    <div className="flex w-full shrink-0 md:w-1/2 md:justify-end">{children}</div>
  </div>
);

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

  const handleSaveAll = async () => {
    try {
      await updateSettings({
        payload: {
          dashboardDomain: formData.dashboardDomain,
          defaultWildcardDomain: formData.defaultWildcardDomain,
        },
      });
      toast.success('Settings updated successfully');
    } catch {
      toast.error('Failed to update settings');
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-50 w-full rounded-2xl" />
        <Skeleton className="h-75 w-full rounded-2xl" />
      </div>
    );
  }

  return (
    <div className="space-y-6 pb-12">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Globe className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Domains</h1>
            <p className="text-muted-foreground text-sm">
              Configure the dashboard domain and wildcard root domain for deployed services.
            </p>
          </div>
        </div>

        <div className="flex shrink-0 items-center gap-3">
          <Button onClick={handleSaveAll} disabled={isPending}>
            <Check className="mr-2 h-4 w-4" />
            {isPending ? 'Saving...' : 'Save Changes'}
          </Button>
        </div>
      </div>

      <Section icon={<Globe className="h-4 w-4" />} title="Domain Configuration">
        <Row
          label="Dashboard Domain"
          description="Point your custom domain's A record to this server's IP address."
        >
          <Input
            placeholder="vessl.example.com"
            value={formData.dashboardDomain}
            onChange={(e) => setFormData({ ...formData, dashboardDomain: e.target.value })}
            className="font-mono"
          />
        </Row>
        <Row
          label="Wildcard Root Domain"
          description="Used as the default wildcard domain for user deployments (*.apps.example.com)."
        >
          <Input
            placeholder="apps.example.com"
            value={formData.defaultWildcardDomain}
            onChange={(e) => setFormData({ ...formData, defaultWildcardDomain: e.target.value })}
            className="font-mono"
          />
        </Row>
      </Section>
    </div>
  );
};
