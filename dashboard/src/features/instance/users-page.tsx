import { Check, Trash2, Users } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Skeleton } from '#/components/ui/skeleton';
import { useDeleteUser, useInviteUser, useListUsers } from '#/hooks/useUsers';
import type { User } from '#/interfaces/users';

const roleVariant: Record<string, 'default' | 'secondary' | 'outline'> = {
  admin: 'default',
  member: 'secondary',
  viewer: 'outline',
};

const UserRow = ({ user, onDelete }: { user: User; onDelete: (u: User) => void }) => (
  <div className="flex items-center justify-between rounded-lg border border-border/50 bg-card/40 px-4 py-3">
    <div className="flex items-center gap-3">
      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10 font-semibold text-primary text-xs uppercase">
        {user.name.charAt(0)}
      </div>
      <div>
        <p className="font-medium text-sm">{user.name}</p>
        <p className="text-muted-foreground text-xs">{user.email}</p>
      </div>
    </div>
    <div className="flex items-center gap-3">
      <Badge variant={roleVariant[user.role] ?? 'outline'}>{user.role}</Badge>
      <Button
        variant="outline"
        size="sm"
        className="h-7 text-destructive hover:bg-destructive/10 hover:text-destructive"
        onClick={() => onDelete(user)}
      >
        Remove
      </Button>
    </div>
  </div>
);

export const UsersPage = () => {
  const { data, isLoading } = useListUsers();
  const { mutateAsync: deleteUser, isPending: deleting } = useDeleteUser();
  const { mutateAsync: inviteUser, isPending: inviting } = useInviteUser();
  const [target, setTarget] = useState<User | null>(null);
  const [inviteOpen, setInviteOpen] = useState(false);
  const [inviteEmail, setInviteEmail] = useState('');

  const users = data?.data?.records ?? [];

  const confirmDelete = async () => {
    if (!target) return;
    try {
      await deleteUser(target.id);
      toast.success(`${target.name} removed`);
    } catch {
      toast.error('Failed to remove user');
    } finally {
      setTarget(null);
    }
  };

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inviteEmail) return;
    try {
      await inviteUser(inviteEmail);
      toast.success('User invited successfully');
      setInviteOpen(false);
      setInviteEmail('');
    } catch (err) {
      const error = err as Error;
      toast.error(error.message || 'Failed to invite user');
    }
  };

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Users className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Users</h1>
            <p className="text-muted-foreground text-sm">
              Manage who has access to this Vessl instance.
            </p>
          </div>
        </div>
        <Dialog open={inviteOpen} onOpenChange={setInviteOpen}>
          <DialogTrigger asChild>
            <Button className="flex h-11 items-center gap-2 rounded-xl px-6 font-semibold text-xs uppercase tracking-widest transition-all">
              Invite User
            </Button>
          </DialogTrigger>
          <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-[500px] [&>button]:hidden">
            <form onSubmit={handleInvite}>
              <div className="flex flex-col p-8 pb-6">
                <div className="flex items-start justify-between">
                  <div className="flex flex-col">
                    <DialogTitle className="font-bold text-2xl tracking-tight">
                      Invite User
                    </DialogTitle>
                    <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      SEND AN EMAIL INVITATION TO A NEW USER.
                    </DialogDescription>
                  </div>
                  <DialogClose asChild>
                    <Button
                      type="button"
                      className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                    >
                      CLOSE
                    </Button>
                  </DialogClose>
                </div>
              </div>

              <div className="h-px w-full bg-border/50" />

              <div className="p-8">
                <div className="space-y-3">
                  <Label
                    htmlFor="email"
                    className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]"
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
                    className="h-12 rounded-xl border-border/50 bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
                  />
                </div>
              </div>

              <div className="flex items-center justify-end gap-6 p-8 pt-6">
                <Button type="button" variant="ghost" onClick={() => setInviteOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={inviting}>
                  <Check className="mr-2 h-4 w-4" />
                  {inviting ? 'INVITING...' : 'SEND INVITE'}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="space-y-2">
        {isLoading && [1, 2, 3].map((i) => <Skeleton key={i} className="h-14 w-full rounded-lg" />)}
        {!isLoading && users.length === 0 && (
          <p className="py-10 text-center text-muted-foreground text-sm">No users found.</p>
        )}
        {users.map((user) => (
          <UserRow key={user.id} user={user} onDelete={setTarget} />
        ))}
      </div>

      <Dialog open={!!target} onOpenChange={(o) => !o && setTarget(null)}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-[500px] [&>button]:hidden">
          <div className="flex flex-col p-8 pb-6">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="font-bold text-2xl text-destructive tracking-tight">
                  Remove User
                </DialogTitle>
                <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  THIS WILL PERMANENTLY REMOVE {target?.email}.
                </DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  variant="outline"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="flex items-center justify-end gap-6 p-8 pt-6">
            <Button onClick={() => setTarget(null)}>Cancel</Button>
            <Button onClick={confirmDelete} disabled={deleting} variant="destructive">
              <Trash2 className="mr-2 h-4 w-4" />
              {deleting ? 'REMOVING...' : 'REMOVE USER'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};
