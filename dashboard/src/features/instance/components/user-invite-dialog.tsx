import { Check, Mail } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useInviteUser } from '#/hooks/useUsers';

interface UserInviteDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function UserInviteDialog({ open, onOpenChange }: UserInviteDialogProps) {
  const { mutateAsync: inviteUser, isPending: inviting } = useInviteUser();
  const [inviteEmail, setInviteEmail] = useState('');

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inviteEmail) return;
    try {
      await inviteUser(inviteEmail);
      toast.success('User invited successfully');
      onOpenChange(false);
      setInviteEmail('');
    } catch (err) {
      const error = err as Error;
      toast.error(error.message || 'Failed to invite user');
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-[400px] [&>button]:hidden">
        <form onSubmit={handleInvite}>
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                  <Mail className="h-5 w-5 text-primary" />
                  Invite User
                </DialogTitle>
                <DialogDescription className="mt-1.5 flex items-center gap-1.5 font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                  <Mail className="h-3 w-3" />
                  Send an email invitation
                </DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  type="button"
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="h-px w-full bg-border/50" />

          <div className="px-5 pt-4 pb-5">
            <div className="space-y-2.5">
              <Label
                htmlFor="email"
                className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]"
              >
                EMAIL ADDRESS
              </Label>
              <Input
                id="email"
                type="email"
                placeholder="name@example.com"
                value={inviteEmail}
                onChange={(e) => setInviteEmail(e.target.value)}
                required
                className="h-10 rounded-lg border-border/50 bg-background/80 px-3 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-1 focus:ring-primary/20"
              />
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              type="button"
              variant="ghost"
              onClick={() => onOpenChange(false)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={inviting}
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              <Check className="h-3.5 w-3.5" />
              {inviting ? 'Inviting...' : 'Send Invite'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
