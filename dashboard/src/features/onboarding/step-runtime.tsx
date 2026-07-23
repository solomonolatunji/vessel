import { useFormContext } from 'react-hook-form';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { SetupSchema } from './schema';

export const StepRuntime = () => {
  const {
    register,
    formState: { errors },
  } = useFormContext<SetupSchema>();

  return (
    <div className="slide-in-from-right-4 animate-in duration-300">
      <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
        <div className="space-y-1.5 md:col-span-2">
          <Label className="font-medium text-foreground/90 text-sm">CODEDOCK_JWT_SECRET</Label>
          <Input
            {...register('env.jwtSecret')}
            placeholder="Generated if blank or enter a 32-character secret"
            className="h-11 rounded-xl border-border bg-card px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-1 focus:ring-primary/50"
          />
          {errors.env?.jwtSecret && (
            <p className="text-destructive text-xs">{errors.env.jwtSecret.message}</p>
          )}
        </div>

        <div className="space-y-1.5">
          <Label className="font-medium text-foreground/90 text-sm">CODEDOCK_DATA_DIR</Label>
          <Input
            {...register('env.dataDir')}
            className="h-11 rounded-xl border-border bg-card px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-1 focus:ring-primary/50"
          />
          {errors.env?.dataDir && (
            <p className="text-destructive text-xs">{errors.env.dataDir.message}</p>
          )}
        </div>

        <div className="space-y-1.5">
          <Label className="font-medium text-foreground/90 text-sm">PORT</Label>
          <Input
            type="number"
            {...register('env.port', { valueAsNumber: true })}
            className="h-11 rounded-xl border-border bg-card px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-1 focus:ring-primary/50"
          />
          {errors.env?.port && (
            <p className="text-destructive text-xs">{errors.env.port.message}</p>
          )}
        </div>

        <div className="space-y-1.5 md:col-span-2">
          <Label className="font-medium text-foreground/90 text-sm">CODEDOCK_DASHBOARD_URL</Label>
          <Input
            {...register('env.dashboardUrl')}
            placeholder="http://localhost:3000"
            className="h-11 rounded-xl border-border bg-card px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-1 focus:ring-primary/50"
          />
          {errors.env?.dashboardUrl && (
            <p className="text-destructive text-xs">{errors.env.dashboardUrl.message}</p>
          )}
        </div>
      </div>
    </div>
  );
};
