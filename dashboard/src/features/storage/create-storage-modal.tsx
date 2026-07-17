import { HardDrive } from 'lucide-react';
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
import { useCreate } from '#/hooks/useStorage';

interface CreateStorageModalProps {
  onClose: () => void;
}

export const CreateStorageModal = ({ onClose }: CreateStorageModalProps) => {
  const [name, setName] = useState('');
  // In a real implementation these might come from a context or selector
  const [projectId] = useState('global');
  const [environmentId] = useState('global');

  const createMutation = useCreate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Name is required');
      return;
    }

    try {
      await createMutation.mutateAsync({
        payload: {
          name,
          projectId,
          environmentId,
          type: 'minio',
        },
      });
      toast.success('Storage container created');
      onClose();
    } catch {
      toast.error('Failed to create storage container');
    }
  };

  return (
    <Dialog open onOpenChange={() => onClose()}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-2xl [&>button]:hidden">
        <div className="flex flex-col p-8 pb-6">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
                <HardDrive className="h-5 w-5 text-primary" />
              </div>
              <div className="flex flex-col">
                <DialogTitle className="font-bold text-2xl text-foreground tracking-tight">
                  Create instance
                </DialogTitle>
                <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  MINIO STORAGE CONTAINER
                </DialogDescription>
              </div>
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

        <div className="h-px w-full bg-border/50" />

        <form onSubmit={handleSubmit}>
          <div className="space-y-8 p-8">
            <div className="space-y-3">
              <Label
                htmlFor="name"
                className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]"
              >
                INSTANCE NAME
              </Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. main-storage"
                className="h-12 rounded-xl border-border/50 bg-background/50 px-4"
                autoFocus
              />
            </div>
          </div>

          <div className="h-px w-full bg-border/50" />

          <div className="flex items-center justify-end bg-muted/20 p-6">
            <Button
              type="submit"
              disabled={createMutation.isPending || !name.trim()}
              className="h-11 bg-primary px-8 font-bold text-primary-foreground text-xs uppercase tracking-wider"
            >
              {createMutation.isPending ? 'CREATING...' : 'CREATE CONTAINER'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
