import { useFormContext } from 'react-hook-form';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { SetupSchema } from './schema';

export const StepBackups = () => {
  const { register, watch } = useFormContext<SetupSchema>();
  const s3Skip = watch('s3Skip');

  return (
    <div className="slide-in-from-right-4 animate-in duration-300">
      <label className="mb-5 flex cursor-pointer items-center gap-2">
        <input
          type="checkbox"
          className="h-4 w-4 rounded border-input text-primary focus:ring-primary"
          {...register('s3Skip')}
        />
        <span className="font-medium text-foreground/90 text-sm">Skip for now</span>
      </label>

      {!s3Skip && (
        <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
          <div className="space-y-1.5 md:col-span-2">
            <Label htmlFor="s3AccountId" className="font-medium text-foreground/90 text-sm">
              R2 account ID
            </Label>
            <Input
              id="s3AccountId"
              className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              {...register('s3AccountId')}
            />
          </div>
          <div className="space-y-1.5 md:col-span-2">
            <Label htmlFor="s3Bucket" className="font-medium text-foreground/90 text-sm">
              R2 bucket
            </Label>
            <Input
              id="s3Bucket"
              className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              {...register('s3Bucket')}
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="s3AccessKeyId" className="font-medium text-foreground/90 text-sm">
              R2 access key ID
            </Label>
            <Input
              id="s3AccessKeyId"
              className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              {...register('s3AccessKeyId')}
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="s3SecretAccessKey" className="font-medium text-foreground/90 text-sm">
              R2 secret access key
            </Label>
            <Input
              id="s3SecretAccessKey"
              type="password"
              className="h-11 rounded-xl border-border bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              {...register('s3SecretAccessKey')}
            />
          </div>
        </div>
      )}
    </div>
  );
};
