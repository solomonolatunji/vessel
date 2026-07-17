import { useFormContext } from 'react-hook-form';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { SetupSchema } from './schema';

export const StepGithub = () => {
  const { register } = useFormContext<SetupSchema>();

  return (
    <div className="slide-in-from-right-4 animate-in duration-300">
      <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
        <div className="space-y-1.5 md:col-span-2">
          <Label htmlFor="githubAppId" className="font-medium text-foreground/90 text-sm">
            GitHub App ID
          </Label>
          <Input
            id="githubAppId"
            className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...register('githubAppId')}
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="githubClientId" className="font-medium text-foreground/90 text-sm">
            Client ID
          </Label>
          <Input
            id="githubClientId"
            className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...register('githubClientId')}
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="githubClientSecret" className="font-medium text-foreground/90 text-sm">
            Client Secret
          </Label>
          <Input
            id="githubClientSecret"
            type="password"
            className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...register('githubClientSecret')}
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="githubPrivateKey" className="font-medium text-foreground/90 text-sm">
            Private Key
          </Label>
          <Input
            id="githubPrivateKey"
            type="password"
            className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...register('githubPrivateKey')}
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="githubWebhookSecret" className="font-medium text-foreground/90 text-sm">
            Webhook Secret
          </Label>
          <Input
            id="githubWebhookSecret"
            type="password"
            className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...register('githubWebhookSecret')}
          />
        </div>
      </div>
    </div>
  );
};
