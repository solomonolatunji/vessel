import {
  Mail01Icon,
  SecurityPasswordIcon,
  UserIcon,
  ViewIcon,
  ViewOffIcon,
} from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useRegister } from '#/hooks/useAuth';

export const RegisterForm = () => {
  const navigate = useNavigate();
  const { mutateAsync: register, isPending } = useRegister();

  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !email || !password) {
      toast.error('Please fill in all fields');
      return;
    }

    try {
      await register({ details: { name, email, password } });
      toast.success('Account created successfully');
      navigate({ to: '/' });
    } catch (error: unknown) {
      toast.error((error as Error)?.message || 'Failed to create account');
    }
  };

  return (
    <form onSubmit={handleRegister} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="name" className="text-sm font-medium">
          Full Name
        </Label>
        <div className="relative">
          <div className="absolute left-3 top-3.5 text-muted-foreground">
            <HugeiconsIcon icon={UserIcon} className="h-5 w-5" />
          </div>
          <Input
            id="name"
            type="text"
            placeholder="John Doe"
            className="h-12 pl-10 bg-background text-base"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="email" className="text-sm font-medium">
          Email
        </Label>
        <div className="relative">
          <div className="absolute left-3 top-3.5 text-muted-foreground">
            <HugeiconsIcon icon={Mail01Icon} className="h-5 w-5" />
          </div>
          <Input
            id="email"
            type="email"
            placeholder="name@example.com"
            className="h-12 pl-10 bg-background text-base"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="password" className="text-sm font-medium">
          Password
        </Label>
        <div className="relative">
          <div className="absolute left-3 top-3.5 text-muted-foreground">
            <HugeiconsIcon icon={SecurityPasswordIcon} className="h-5 w-5" />
          </div>
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            className="h-12 pl-10 pr-12 bg-background text-base"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <button
            type="button"
            className="absolute right-3 top-3.5 text-muted-foreground hover:text-foreground transition-colors"
            onClick={() => setShowPassword(!showPassword)}
          >
            <HugeiconsIcon icon={showPassword ? ViewOffIcon : ViewIcon} className="h-5 w-5" />
          </button>
        </div>
        <p className="text-xs text-muted-foreground pt-1">Must be at least 8 characters long.</p>
      </div>

      <Button type="submit" className="h-12 w-full text-base font-medium mt-2" disabled={isPending}>
        {isPending ? 'Creating account...' : 'Create Account'}
      </Button>
    </form>
  );
};
