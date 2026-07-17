import { Lock } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useResetPassword } from '#/hooks/useAuth';

interface ResetPasswordFormProps {
  token: string;
}

export const ResetPasswordForm = ({ token }: ResetPasswordFormProps) => {
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');

  const { mutate, isPending } = useResetPassword();

  const handleSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!newPassword || !confirmPassword) return;

    if (newPassword !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    setError('');

    mutate({ token, newPassword });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-1.5">
        <Label htmlFor="new-password" className="font-medium text-foreground/90 text-sm">
          New Password
        </Label>
        <div className="group relative">
          <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
            <Lock className="h-4 w-4" />
          </div>
          <Input
            id="new-password"
            type="password"
            placeholder="Enter new password"
            className="h-11 rounded-xl border-border bg-background/80 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            value={newPassword}
            onChange={(e) => {
              setNewPassword(e.target.value);
              setError('');
            }}
            required
            disabled={isPending}
          />
        </div>
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="confirm-password" className="font-medium text-foreground/90 text-sm">
          Confirm Password
        </Label>
        <div className="group relative">
          <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
            <Lock className="h-4 w-4" />
          </div>
          <Input
            id="confirm-password"
            type="password"
            placeholder="Confirm new password"
            className="h-11 rounded-xl border-border bg-background/80 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            value={confirmPassword}
            onChange={(e) => {
              setConfirmPassword(e.target.value);
              setError('');
            }}
            required
            disabled={isPending}
          />
        </div>
      </div>

      {error && <p className="font-medium text-destructive text-sm">{error}</p>}

      <Button
        type="submit"
        disabled={isPending || !newPassword || !confirmPassword}
        className="h-11 w-full rounded-xl bg-linear-to-r from-primary to-purple-600 font-semibold text-sm shadow-lg shadow-primary/30 transition-all duration-200 hover:brightness-110 active:scale-[0.985]"
      >
        {isPending ? 'Resetting...' : 'Reset Password'}
      </Button>
    </form>
  );
};
