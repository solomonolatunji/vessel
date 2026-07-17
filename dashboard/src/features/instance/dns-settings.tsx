import { CheckCircle2 } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Button } from "#/components/ui/button";
import { Input } from "#/components/ui/input";
import { Label } from "#/components/ui/label";
import { Skeleton } from "#/components/ui/skeleton";
import { useGetSettings, useUpdateSettings } from "#/hooks/useSettings";

const providers = [
  { id: "cloudflare", name: "Cloudflare", sub: "API KEY / TOKEN + ZONE" },
  { id: "namecheap", name: "Namecheap", sub: "API USER + KEY" },
  { id: "spaceship", name: "Spaceship", sub: "API KEY + SECRET" },
];

export const DnsSettings = () => {
  const { data, isLoading } = useGetSettings();
  const { mutateAsync: updateSettings, isPending } = useUpdateSettings();

  const [activeProvider, setActiveProvider] = useState("cloudflare");
  const [formData, setFormData] = useState({
    cloudflareApiToken: "",
    namecheapApiUser: "",
    namecheapApiKey: "",
    namecheapClientIp: "",
    spaceshipApiKey: "",
    // Added for UI fidelity, not currently saved to backend
    cloudflareEmail: "",
    cloudflareZoneId: "",
    spaceshipApiSecret: "",
  });

  useEffect(() => {
    if (data?.data) {
      setFormData((prev) => ({
        ...prev,
        cloudflareApiToken: data.data.cloudflareApiToken || "",
        namecheapApiUser: data.data.namecheapApiUser || "",
        namecheapApiKey: data.data.namecheapApiKey || "",
        namecheapClientIp: data.data.namecheapClientIp || "",
        spaceshipApiKey: data.data.spaceshipApiKey || "",
      }));
    }
  }, [data?.data]);

  const handleSaveProvider = async (provider: string) => {
    try {
      const payload: any = {};
      if (provider === "cloudflare") {
        payload.cloudflareApiToken = formData.cloudflareApiToken;
      } else if (provider === "namecheap") {
        payload.namecheapApiUser = formData.namecheapApiUser;
        payload.namecheapApiKey = formData.namecheapApiKey;
        payload.namecheapClientIp = formData.namecheapClientIp;
      } else if (provider === "spaceship") {
        payload.spaceshipApiKey = formData.spaceshipApiKey;
      }

      await updateSettings({ payload });
      toast.success("Provider credentials saved");
    } catch {
      toast.error("Failed to save provider credentials");
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
      {/* Header */}
      <div className="flex items-center justify-between border-border/50 border-b pb-6">
        <div className="flex items-center gap-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
            <span className="font-bold text-[10px] tracking-widest">API</span>
          </div>
          <div>
            <h1 className="font-bold text-2xl text-foreground">
              DNS Management API
            </h1>
            <p className="mt-1.5 font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
              PROVIDER CREDENTIALS FOR AUTOMATED DNS.
            </p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <img
            src={`/dns-providers/${activeProvider}.svg`}
            alt={activeProviderData?.name}
            className="h-4 w-auto mix-blend-screen brightness-150 contrast-125 grayscale"
          />
          <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            {activeProviderData?.name}
          </span>
        </div>
      </div>

      {/* Cards */}
      <div className="grid grid-cols-1 gap-6 pt-2 md:grid-cols-3">
        {providers.map((p) => {
          const isActive = activeProvider === p.id;
          return (
            <div
              key={p.id}
              onClick={() => setActiveProvider(p.id)}
              className={`group relative cursor-pointer rounded-2xl border p-5 transition-all duration-200 ${
                isActive
                  ? "border-primary/50 bg-card/40 shadow-sm"
                  : "border-border/50 bg-background/50 hover:border-border/80 hover:bg-card/40"
              }`}
            >
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                  <div
                    className={`flex h-10 w-10 items-center justify-center rounded-xl border transition-colors ${
                      isActive
                        ? "border-primary/30 bg-primary/5"
                        : "border-border/50 bg-background"
                    }`}
                  >
                    <img
                      src={`/dns-providers/${p.id}.svg`}
                      alt={p.name}
                      className="h-5 w-auto"
                    />
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
                      ? "border-primary/30 text-primary"
                      : "border-border/50 text-muted-foreground/50"
                  }`}
                >
                  NEW
                </span>
              </div>
            </div>
          );
        })}
      </div>

      {/* Active Content */}
      <div className="mt-4 space-y-10 rounded-2xl border border-border/50 bg-card/40 p-8">
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

        <div className="space-y-6">
          {activeProvider === "cloudflare" && (
            <div className="fade-in-50 animate-in space-y-6 duration-300">
              <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div className="space-y-3">
                  <Label className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]">
                    API KEY / TOKEN
                  </Label>
                  <Input
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
                </div>
                <div className="space-y-3">
                  <Label className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]">
                    ACCOUNT EMAIL
                  </Label>
                  <Input
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
                <Label className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]">
                  ZONE ID
                </Label>
                <Input
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

          {activeProvider === "namecheap" && (
            <div className="fade-in-50 animate-in space-y-6 duration-300">
              <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div className="space-y-3">
                  <Label className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]">
                    API USER
                  </Label>
                  <Input
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
                  <Label className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]">
                    API KEY
                  </Label>
                  <Input
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
                <Label className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]">
                  CLIENT IP
                </Label>
                <Input
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

          {activeProvider === "spaceship" && (
            <div className="fade-in-50 animate-in space-y-6 duration-300">
              <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div className="space-y-3">
                  <Label className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]">
                    API KEY
                  </Label>
                  <Input
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
                  <Label className="font-bold text-[10px] text-muted-foreground/90 uppercase tracking-[0.15em]">
                    API SECRET
                  </Label>
                  <Input
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
      </div>
    </div>
  );
};
