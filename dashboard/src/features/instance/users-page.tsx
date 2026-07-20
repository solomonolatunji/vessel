import { format } from 'date-fns';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { useListUsers } from '#/hooks/useUsers';
import type { User } from '#/interfaces/users';
import { UserDeleteDialog } from './components/user-delete-dialog';
import { UserInviteDialog } from './components/user-invite-dialog';

const formatLastLogin = (dateStr?: string) => {
  if (!dateStr) return 'NEVER';
  try {
    return format(new Date(dateStr), 'MMM d, h:mm a').toUpperCase();
  } catch {
    return 'UNKNOWN';
  }
};

const StatBox = ({ icon: Icon, label, value }: { icon: any; label: string; value: number }) => (
  <div className="flex items-center justify-between rounded-lg border border-border/50 bg-card/20 px-4 py-3">
    <div className="flex items-center gap-3">
      <Icon className="h-4 w-4 text-muted-foreground" />
      <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
        {label}
      </span>
    </div>
    <span className="font-bold text-sm">{value || 0}</span>
  </div>
);

const UserRow = ({ user, onDelete }: { user: User; onDelete: (u: User) => void }) => (
  <div className="group relative rounded-2xl border border-border/50 bg-card/40 p-6">
    <div className="flex items-start justify-between">
      <div className="space-y-2">
        <div className="flex items-center gap-3">
          <h3 className="font-bold text-foreground text-xl">
            {user.name || user.email.split('@')[0]}
          </h3>
          <div className="rounded border border-primary/30 bg-primary/10 px-2 py-0.5 font-bold text-[10px] text-primary uppercase tracking-widest">
            {user.role}
          </div>
        </div>
        <p className="font-mono text-muted-foreground text-sm">{user.email}</p>
      </div>

      <div className="flex flex-col items-end gap-2">
        <div className="font-mono text-[10px] text-muted-foreground uppercase tracking-widest">
          LAST LOGIN {formatLastLogin(user.lastLogin)}
        </div>
        <Button
          variant="outline"
          onClick={() => onDelete(user)}
          className="absolute top-6 right-6 h-8 w-8 border-border/50 bg-transparent p-0 text-muted-foreground opacity-0 transition-colors hover:border-destructive/30 hover:bg-destructive/10 hover:text-destructive group-hover:opacity-100"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
    </div>

    <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
      <StatBox icon={Users} label="Projects" value={user.projectsCount || 0} />
      <StatBox icon={Server} label="Active Services" value={user.servicesCount || 0} />
      <StatBox icon={Key} label="API Keys" value={user.apiKeysCount || 0} />
    </div>
  </div>
);

export const UsersPage = () => {
  const { data, isLoading } = useListUsers();
  const [target, setTarget] = useState<User | null>(null);
  const [inviteOpen, setInviteOpen] = useState(false);

  const users = data?.data?.records ?? [];

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
        <Button onClick={() => setInviteOpen(true)} className="gap-2">
          <Plus className="h-4 w-4" />
          INVITE USER
        </Button>
      </div>

      <div className="grid grid-cols-1 gap-6">
        {isLoading &&
          [1, 2, 3].map((i) => <Skeleton key={i} className="h-30 w-full rounded-2xl" />)}
        {!isLoading && users.length === 0 && (
          <div className="flex items-center gap-3 rounded-2xl border border-border/50 bg-card/40 p-6">
            <Users className="h-4 w-4 text-muted-foreground" />
            <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              NO USERS
            </span>
          </div>
        )}
        {users.map((user) => (
          <UserRow key={user.id} user={user} onDelete={setTarget} />
        ))}
      </div>

      <UserInviteDialog open={inviteOpen} onOpenChange={setInviteOpen} />
      <UserDeleteDialog target={target} onClose={() => setTarget(null)} />
    </div>
  );
};
