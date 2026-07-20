import { Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '#/components/ui/dialog';
import { useDeleteUser } from '#/hooks/useUsers';
import type { User } from '#/interfaces/users';

interface UserDeleteDialogProps {
  target: User | null;
  onClose: () => void;
}

export function UserDeleteDialog({ target, onClose }: UserDeleteDialogProps) {
  const { mutateAsync: deleteUser, isPending: deleting } = useDeleteUser();

  const confirmDelete = async () => {
    if (!target) return;
    try {
      await deleteUser(target.id);
      toast.success(`${target.name} removed`);
      onClose();
    } catch {
      toast.error('Failed to remove user');
    }
  };

  return (
    <Dialog open={!!target} onOpenChange={(o) => !o && onClose()}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-100 [&>button]:hidden">
        <div className="p-5">
          <div className="flex items-start justify-between">
            <div className="flex flex-col">
              <DialogTitle className="flex items-center gap-2 font-bold text-destructive text-xl tracking-tight">
                <Trash2 className="h-5 w-5" />
                Remove User
              </DialogTitle>
              <DialogDescription>This will permanently remove {target?.email}</DialogDescription>
            </div>
            <DialogClose asChild>
              <Button
                variant="ghost"
                className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
              >
                CLOSE
              </Button>
            </DialogClose>
          </div>
        </div>
        <div className="flex items-center justify-end gap-3 p-5 pt-0">
          <Button
            variant="ghost"
            onClick={onClose}
            className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
          >
            Cancel
          </Button>
          <Button
            onClick={confirmDelete}
            disabled={deleting}
            variant="destructive"
            className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
          >
            <Trash2 className="h-3.5 w-3.5" />
            {deleting ? 'Removing...' : 'Remove User'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
