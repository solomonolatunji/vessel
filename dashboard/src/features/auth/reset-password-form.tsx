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

  const handleSubmit = (e: React.FormEvent) => {
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
    <form onSubmit={handleSubmit} className="space-y-5">
      <div className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="new-password" className="text-sm font-medium">
            New Password
          </Label>
          <div className="relative">
            <div className="absolute left-3 top-3.5 text-muted-foreground">
              <Lock className="h-5 w-5" />
            </div>
            <Input
              id="new-password"
              type="password"
              placeholder="Enter new password"
              className="h-12 pl-10 bg-background text-base"
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

        <div className="space-y-2">
          <Label htmlFor="confirm-password" className="text-sm font-medium">
            Confirm Password
          </Label>
          <div className="relative">
            <div className="absolute left-3 top-3.5 text-muted-foreground">
              <Lock className="h-5 w-5" />
            </div>
            <Input
              id="confirm-password"
              type="password"
              placeholder="Confirm new password"
              className="h-12 pl-10 bg-background text-base"
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

        {error && <p className="text-sm font-medium text-destructive">{error}</p>}
      </div>

      <Button
        type="submit"
        className="h-12 w-full text-base font-medium mt-2"
        disabled={isPending || !newPassword || !confirmPassword}
      >
        {isPending ? 'Resetting Password...' : 'Reset Password'}
      </Button>
    </form>
  );
};
