import { Eye, EyeOff, Lock, Mail, User } from 'lucide-react';
import { useState } from 'react';
import { useFormContext } from 'react-hook-form';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { SetupSchema } from './schema';

export const StepOwner = () => {
  const [showPassword, setShowPassword] = useState(false);
  const {
    register,
    formState: { errors },
  } = useFormContext<SetupSchema>();

  return (
    <div className="slide-in-from-right-4 animate-in duration-300">
      <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
        <div className="space-y-1.5">
          <Label htmlFor="name" className="font-medium text-foreground/90 text-sm">
            Name
          </Label>
          <div className="group relative">
            <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
              <User className="h-4 w-4" />
            </div>
            <Input
              id="name"
              placeholder="John Doe"
              className="h-11 rounded-xl border-border bg-background/80 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              {...register('name')}
            />
          </div>
          {errors.name && <p className="text-destructive text-xs">{errors.name.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="email" className="font-medium text-foreground/90 text-sm">
            Email
          </Label>
          <div className="group relative">
            <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
              <Mail className="h-4 w-4" />
            </div>
            <Input
              id="email"
              type="email"
              placeholder="name@example.com"
              className="h-11 rounded-xl border-border bg-background/80 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              {...register('email')}
            />
          </div>
          {errors.email && <p className="text-destructive text-xs">{errors.email.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="password" className="font-medium text-foreground/90 text-sm">
            Password
          </Label>
          <div className="group relative">
            <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
              <Lock className="h-4 w-4" />
            </div>
            <Input
              id="password"
              type={showPassword ? 'text' : 'password'}
              className="h-11 rounded-xl border-border bg-background/80 pr-10 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              {...register('password')}
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute top-1/2 right-3.5 -translate-y-1/2 text-muted-foreground transition-colors hover:text-foreground"
            >
              {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </button>
          </div>
          {errors.password && <p className="text-destructive text-xs">{errors.password.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="confirmPassword" className="font-medium text-foreground/90 text-sm">
            Confirm Password
          </Label>
          <div className="group relative">
            <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
              <Lock className="h-4 w-4" />
            </div>
            <Input
              id="confirmPassword"
              type={showPassword ? 'text' : 'password'}
              className="h-11 rounded-xl border-border bg-background/80 pr-10 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              {...register('confirmPassword')}
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute top-1/2 right-3.5 -translate-y-1/2 text-muted-foreground transition-colors hover:text-foreground"
            >
              {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </button>
          </div>
          {errors.confirmPassword && (
            <p className="text-destructive text-xs">{errors.confirmPassword.message}</p>
          )}
        </div>
      </div>
    </div>
  );
};
