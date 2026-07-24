import { Button } from '#/components/ui/button';
import { useEnabledOAuthProviders } from '#/hooks/useOAuth';
import { oauthService } from '#/services/oauth';

export const OAuthButtons = () => {
  const { data: enabledProvidersData } = useEnabledOAuthProviders();
  const enabledProviders = (enabledProvidersData?.data || []).map((p) =>
    p.providerName.toLowerCase()
  );

  const handleOAuthLogin = (provider: string) => {
    oauthService.triggerOAuthLogin(provider);
  };

  if (!enabledProviders.length) return null;

  return (
    <>
      <div className="mb-5 space-y-2">
        {enabledProviders.map((provider) => (
          <Button
            key={provider}
            variant="outline"
            type="button"
            onClick={() => handleOAuthLogin(provider)}
            className="h-10 w-full rounded-xl border-border/80 font-medium text-sm capitalize transition-all hover:border-primary/30 hover:bg-muted/50"
          >
            Continue with {provider}
          </Button>
        ))}
      </div>

      <div className="relative mb-5">
        <div className="absolute inset-0 flex items-center">
          <div className="h-px w-full bg-linear-to-r from-transparent via-border to-transparent" />
        </div>
        <div className="relative flex justify-center">
          <span className="bg-card/80 px-3 text-muted-foreground text-xs uppercase tracking-widest">
            or
          </span>
        </div>
      </div>
    </>
  );
};
