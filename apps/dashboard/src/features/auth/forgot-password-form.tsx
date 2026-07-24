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

  const handleSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!email) return;
    mutate(email);
  };

  if (isSuccess) {
    return (
      <div className="space-y-4 text-center">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
          <Mail className="h-6 w-6 text-primary" />
        </div>
        <p className="font-medium text-foreground text-lg tracking-tight">Check your email</p>
        <p className="text-muted-foreground text-sm">
          If an account with that email exists, we've sent you instructions to reset your password.
        </p>
        <div className="mt-6">
          <Link to="/signin" className="font-medium text-primary text-sm hover:underline">
            Back to sign in
          </Link>
        </div>
      </div>
    );
  }

  return (
    <>
      <form onSubmit={handleSubmit} className="space-y-4">
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
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isPending}
            />
          </div>
        </div>

        <Button
          type="submit"
          disabled={isPending || !email}
          className="h-11 w-full rounded-xl bg-linear-to-r from-primary to-purple-600 font-semibold text-sm shadow-lg shadow-primary/30 transition-all duration-200 hover:brightness-110 active:scale-[0.985]"
        >
          {isPending ? 'Sending...' : 'Send Reset Link'}
        </Button>
      </form>

      <div className="mt-6 text-center text-sm">
        <span className="text-muted-foreground">Remember your password? </span>
        <Link
          to="/signin"
          className="font-medium text-primary underline-offset-4 transition-colors hover:text-primary-hover hover:underline"
        >
          Sign in
        </Link>
      </div>
    </>
  );
};
