import { Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Dialog, DialogContent, DialogTitle } from '#/components/ui/dialog';
import { useDeleteToken } from '#/hooks/useProfile';

interface ApiKeyDeleteDialogProps {
  deleteId: string | null;
  onClose: () => void;
}

export function ApiKeyDeleteDialog({ deleteId, onClose }: ApiKeyDeleteDialogProps) {
  const deleteToken = useDeleteToken();

  const handleDelete = () => {
    if (!deleteId) return;
    deleteToken.mutate(
      { id: deleteId },
      {
        onSuccess: () => {
          toast.success('API key deleted');
          onClose();
        },
        onError: (err) => {
          toast.error(err.message || 'Failed to delete API key');
        },
      }
    );
  };

  return (
    <Dialog open={!!deleteId} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
        <div className="p-5">
          <div className="flex items-start justify-between">
            <div className="flex flex-col">
              <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                <Trash2 className="h-5 w-5 text-destructive" />
                Delete API Key
              </DialogTitle>
            </div>
          </div>
          <p className="mt-4 text-muted-foreground text-sm">
            Are you sure you want to delete this API key? Any applications or scripts using it will
            immediately lose access. This action cannot be undone.
          </p>
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
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteToken.isPending}
            className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
          >
            <Trash2 className="h-3.5 w-3.5" />
            {deleteToken.isPending ? 'Deleting...' : 'Delete Key'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
