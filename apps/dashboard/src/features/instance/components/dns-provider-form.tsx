import { CheckCircle2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

interface DnsProviderFormProps {
  activeProvider: string;
  formData: Record<string, string>;
  setFormData: (data: Record<string, string>) => void;
  isPending: boolean;
  handleSaveProvider: (provider: string) => void;
}

export function DnsProviderForm({
  activeProvider,
  formData,
  setFormData,
  isPending,
  handleSaveProvider,
}: DnsProviderFormProps) {
  return (
    <div className="space-y-6">
      {activeProvider === 'cloudflare' && (
        <div className="fade-in-50 animate-in space-y-6 duration-300">
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div className="space-y-3">
              <Label
                htmlFor="cf-api-token"
                className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]"
              >
                API KEY / TOKEN
              </Label>
              <Input
                id="cf-api-token"
                type="password"
                placeholder="Cloudflare API key or DNS token"
                value={formData.cloudflareApiToken}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    cloudflareApiToken: e.target.value,
                  })
                }
                className="h-12 rounded-xl border-border/50 bg-background/80 px-4 font-medium"
              />
              <p className="text-[11px] text-muted-foreground/70">
                API Tokens only. Global API keys are not supported.
              </p>
            </div>
            <div className="space-y-3">
              <Label
                htmlFor="cf-email"
                className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]"
              >
                ACCOUNT EMAIL
              </Label>
              <Input
                id="cf-email"
                placeholder="Only needed for global API keys"
                value={formData.cloudflareEmail}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    cloudflareEmail: e.target.value,
                  })
                }
                className="h-12 rounded-xl border-border/50 bg-background/80 px-4 font-medium placeholder:text-muted-foreground/40"
              />
            </div>
          </div>
          <div className="space-y-3">
            <Label
              htmlFor="cf-zone-id"
              className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]"
            >
              ZONE ID
            </Label>
            <Input
              id="cf-zone-id"
              placeholder="Optional zone ID"
              value={formData.cloudflareZoneId}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  cloudflareZoneId: e.target.value,
                })
              }
              className="h-12 rounded-xl border-border/50 bg-background/80 px-4 font-medium placeholder:text-muted-foreground/40"
            />
          </div>
        </div>
      )}

      {activeProvider === 'namecheap' && (
        <div className="fade-in-50 animate-in space-y-6 duration-300">
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div className="space-y-3">
              <Label
                htmlFor="nc-api-user"
                className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]"
              >
                API USER
              </Label>
              <Input
                id="nc-api-user"
                placeholder="username"
                value={formData.namecheapApiUser}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    namecheapApiUser: e.target.value,
                  })
                }
                className="h-12 rounded-xl border-border/50 bg-background/80 px-4 font-medium"
              />
            </div>
            <div className="space-y-3">
              <Label
                htmlFor="nc-api-key"
                className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]"
              >
                API KEY
              </Label>
              <Input
                id="nc-api-key"
                type="password"
                placeholder="••••••••••••••••"
                value={formData.namecheapApiKey}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    namecheapApiKey: e.target.value,
                  })
                }
                className="h-12 rounded-xl border-border/50 bg-background/80 px-4 font-medium"
              />
            </div>
          </div>
          <div className="space-y-3">
            <Label
              htmlFor="nc-client-ip"
              className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]"
            >
              CLIENT IP
            </Label>
            <Input
              id="nc-client-ip"
              placeholder="Whitelisted server IP"
              value={formData.namecheapClientIp}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  namecheapClientIp: e.target.value,
                })
              }
              className="h-12 rounded-xl border-border/50 bg-background/80 px-4 font-medium placeholder:text-muted-foreground/40"
            />
          </div>
        </div>
      )}

      {activeProvider === 'spaceship' && (
        <div className="fade-in-50 animate-in space-y-6 duration-300">
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div className="space-y-3">
              <Label
                htmlFor="ss-api-key"
                className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]"
              >
                API KEY
              </Label>
              <Input
                id="ss-api-key"
                type="password"
                placeholder="Spaceship API key"
                value={formData.spaceshipApiKey}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    spaceshipApiKey: e.target.value,
                  })
                }
                className="h-12 rounded-xl border-border/50 bg-background/80 px-4 font-medium"
              />
            </div>
            <div className="space-y-3">
              <Label
                htmlFor="ss-api-secret"
                className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]"
              >
                API SECRET
              </Label>
              <Input
                id="ss-api-secret"
                type="password"
                placeholder="Spaceship API secret"
                value={formData.spaceshipApiSecret}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    spaceshipApiSecret: e.target.value,
                  })
                }
                className="h-12 rounded-xl border-border/50 bg-background/80 px-4 font-medium placeholder:text-muted-foreground/40"
              />
            </div>
          </div>
        </div>
      )}

      <div className="pt-4">
        <Button
          onClick={() => handleSaveProvider(activeProvider)}
          disabled={isPending}
          className="flex h-11 items-center gap-2 rounded-xl border border-primary/20 bg-primary/10 px-6 font-semibold text-primary text-xs uppercase tracking-widest shadow-none transition-all hover:bg-primary/20 hover:text-primary"
        >
          <CheckCircle2 className="h-4 w-4" /> SAVE CREDENTIALS
        </Button>
      </div>
    </div>
  );
}
