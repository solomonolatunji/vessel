import { zodResolver } from '@hookform/resolvers/zod';
import { Eye, EyeOff, Lock, Mail, User } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useRegister } from '#/hooks/useAuth';

const registerSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters long'),
});

type RegisterSchema = z.infer<typeof registerSchema>;

export const RegisterForm = () => {
  const { mutateAsync: registerUser, isPending } = useRegister();
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterSchema>({
    resolver: zodResolver(registerSchema),
    defaultValues: { name: '', email: '', password: '' },
  });

  const onSubmit = async (data: RegisterSchema) => {
    try {
      await registerUser(data);
    } catch {}
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-1.5">
        <Label htmlFor="name" className="font-medium text-foreground/90 text-sm">
          Full Name
        </Label>
        <div className="group relative">
          <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
            <User className="h-4 w-4" />
          </div>
          <Input
            id="name"
            type="text"
            placeholder="John Doe"
            className="h-11 rounded-xl border-border bg-background/80 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...register('name')}
          />
        </div>
        {errors.name && <p className="pl-1 text-destructive text-xs">{errors.name.message}</p>}
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
        {errors.email && <p className="pl-1 text-destructive text-xs">{errors.email.message}</p>}
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
        {errors.password ? (
          <p className="pl-1 text-destructive text-xs">{errors.password.message}</p>
        ) : (
          <p className="pt-1 text-muted-foreground text-xs">Must be at least 8 characters long.</p>
        )}
      </div>

      <Button
        type="submit"
        disabled={isPending}
        className="h-11 w-full rounded-xl bg-linear-to-r from-primary to-purple-600 font-semibold text-sm shadow-lg shadow-primary/30 transition-all duration-200 hover:brightness-110 active:scale-[0.985]"
      >
        {isPending ? 'Creating account...' : 'Create Account'}
      </Button>
    </form>
  );
};
