import { useFormContext } from 'react-hook-form';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { SetupSchema } from './schema';

export const StepDomain = () => {
  const { register } = useFormContext<SetupSchema>();

  return (
    <div className="slide-in-from-right-4 animate-in duration-300">
      <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
        <div className="space-y-1.5">
          <Label htmlFor="dashboardDomain" className="font-medium text-foreground/90 text-sm">
            Dashboard domain
          </Label>
          <Input
            id="dashboardDomain"
            placeholder="app.vessl.dev"
            className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...register('dashboardDomain')}
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="defaultWildcardDomain" className="font-medium text-foreground/90 text-sm">
            Wildcard root domain
          </Label>
          <Input
            id="defaultWildcardDomain"
            placeholder="vessl.dev"
            className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...register('defaultWildcardDomain')}
          />
        </div>
      </div>
    </div>
  );
};
