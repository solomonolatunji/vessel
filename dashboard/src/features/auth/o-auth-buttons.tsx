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
      <div className="flex flex-col space-y-3 mb-6">
        {enabledProviders.map((provider) => (
          <Button
            key={provider}
            variant="outline"
            type="button"
            onClick={() => handleOAuthLogin(provider)}
            className="h-11 w-full capitalize text-base font-medium bg-background border-border hover:bg-muted/50"
          >
            Continue with {provider}
          </Button>
        ))}
      </div>
      <div className="relative mb-6">
        <div className="absolute inset-0 flex items-center">
          <span className="w-full border-t border-border" />
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-background px-2 text-muted-foreground font-semibold">
            Or continue with email
          </span>
        </div>
      </div>
    </>
  );
};
