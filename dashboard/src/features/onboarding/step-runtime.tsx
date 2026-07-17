import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

export const StepRuntime = () => {
  return (
    <div className="slide-in-from-right-4 animate-in duration-300">
      <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
        {['VESSL_JWT_SECRET', 'DATA_DIR', 'PUBLIC_URL', 'PORT', 'BUILDKIT_HOST'].map((key) => (
          <div key={key} className="space-y-1.5">
            <Label className="font-medium text-foreground/90 text-sm">{key}</Label>
            <Input
              disabled
              value="******"
              className="h-11 rounded-xl border-border bg-muted px-4 text-sm transition-all duration-300"
            />
          </div>
        ))}
      </div>
    </div>
  );
};
