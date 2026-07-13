import { Link } from '@tanstack/react-router';
import { Mail } from 'lucide-react';
import { useState } from 'react';

import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

import { useForgotPassword } from '#/hooks/useAuth';

export const ForgotPasswordForm = () => {
  const [email, setEmail] = useState('');
  const { mutate, isPending, isSuccess } = useForgotPassword();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;
    mutate(email);
  };

  if (isSuccess) {
    return (
      <div className="text-center space-y-4">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
          <Mail className="h-6 w-6 text-primary" />
        </div>
        <h3 className="text-xl font-medium">Check your email</h3>
        <p className="text-sm text-muted-foreground">
          If an account with that email exists, we've sent you instructions to reset your password.
        </p>
        <div className="mt-8">
          <Link to="/login" className="text-sm font-medium text-primary hover:underline">
            Back to sign in
          </Link>
        </div>
      </div>
    );
  }

  return (
    <>
      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="space-y-2">
          <Label htmlFor="email" className="text-sm font-medium">
            Email
          </Label>
          <div className="relative">
            <div className="absolute left-3 top-3.5 text-muted-foreground">
              <Mail className="h-5 w-5" />
            </div>
            <Input
              id="email"
              type="email"
              placeholder="name@example.com"
              className="h-12 pl-10 bg-background text-base"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isPending}
            />
          </div>
        </div>

        <Button
          type="submit"
          className="h-12 w-full text-base font-medium mt-2"
          disabled={isPending || !email}
        >
          {isPending ? 'Sending...' : 'Send Reset Link'}
        </Button>
      </form>

      <div className="mt-8 text-center text-sm">
        <span className="text-muted-foreground">Remember your password? </span>
        <Link to="/login" className="font-medium text-primary hover:underline">
          Sign in
        </Link>
      </div>
    </>
  );
};
