import { Check, Clock, FolderOpen, Key } from 'lucide-react';
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
import { useCreateToken } from '#/hooks/useProfile';

interface ApiKeyCreateDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: (plainKey: string) => void;
}

export function ApiKeyCreateDialog({ open, onOpenChange, onSuccess }: ApiKeyCreateDialogProps) {
  const createToken = useCreateToken();

  const [name, setName] = useState('');
  const [accessLevel, setAccessLevel] = useState<'read' | 'read_write'>('read');
  const [projectScope, setProjectScope] = useState<'all' | 'specific'>('all');
  const [expirationDays, setExpirationDays] = useState<number | null>(30);

  const handleCreate = () => {
    if (!name.trim()) {
      toast.error('Please enter a name for the API key');
      return;
    }

    let expiresAt: string | undefined;
    if (expirationDays !== null) {
      const date = new Date();
      date.setDate(date.getDate() + expirationDays);
      expiresAt = date.toISOString();
    }

    createToken.mutate(
      {
        payload: {
          name,
          accessLevel,
          projectScope,
          allowedProjects: [],
          expiresAt,
        },
      },
      {
        onSuccess: (data) => {
          onSuccess(data.plain);
          setName('');
          setAccessLevel('read');
          setProjectScope('all');
          setExpirationDays(30);
        },
        onError: (err) => {
          toast.error(err.message || 'Failed to create API key');
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-2xl [&>button]:hidden">
        <div className="px-5 pt-5 pb-4">
          <div className="flex items-start justify-between">
            <div className="flex flex-col">
              <DialogTitle className="flex items-center gap-2 font-bold text-xl tracking-tight">
                <Key className="h-5 w-5 text-primary" />
                Create API key
              </DialogTitle>
              <DialogDescription className="mt-1.5 flex items-center gap-1.5 font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                <Key className="h-3 w-3" />
                Bearer token access
              </DialogDescription>
            </div>
            <DialogClose asChild>
              <Button className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground">
                CLOSE
              </Button>
            </DialogClose>
          </div>
        </div>

        <div className="h-px w-full bg-border/50" />

        <div className="space-y-5 px-5 pt-4 pb-5">
          <div className="space-y-2.5">
            <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
              NAME
            </Label>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Production deploys"
              className="h-10 rounded-lg border-border/50 bg-background/50 px-3 text-sm"
            />
          </div>

          <div className="space-y-2.5">
            <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
              ACCESS
            </Label>
            <div className="grid grid-cols-2 gap-4">
              <button
                type="button"
                onClick={() => setAccessLevel('read')}
                className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                  accessLevel === 'read'
                    ? 'border-primary/50 bg-primary/10 text-primary'
                    : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                }`}
              >
                <Key className="h-4 w-4" /> READ
              </button>
              <button
                type="button"
                onClick={() => setAccessLevel('read_write')}
                className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                  accessLevel === 'read_write'
                    ? 'border-primary/50 bg-primary/10 text-primary'
                    : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                }`}
              >
                <Key className="h-4 w-4" /> READ AND WRITE
              </button>
            </div>
          </div>

          <div className="space-y-2.5">
            <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
              PROJECTS
            </Label>
            <div className="grid grid-cols-2 gap-4">
              <button
                type="button"
                onClick={() => setProjectScope('all')}
                className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                  projectScope === 'all'
                    ? 'border-primary/50 bg-primary/10 text-primary'
                    : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                }`}
              >
                <FolderOpen className="h-4 w-4" /> ALL PROJECTS
              </button>
              <button
                type="button"
                onClick={() => setProjectScope('specific')}
                className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                  projectScope === 'specific'
                    ? 'border-primary/50 bg-primary/10 text-primary'
                    : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                }`}
              >
                <FolderOpen className="h-4 w-4" /> SPECIFIC PROJECTS
              </button>
            </div>
          </div>

          <div className="space-y-2.5">
            <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
              EXPIRATION
            </Label>
            <div className="grid grid-cols-2 gap-4">
              <button
                type="button"
                onClick={() => setExpirationDays(7)}
                className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                  expirationDays === 7
                    ? 'border-primary/50 bg-primary/10 text-primary'
                    : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                }`}
              >
                <Clock className="h-4 w-4" /> 7 DAYS
              </button>
              <button
                type="button"
                onClick={() => setExpirationDays(30)}
                className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                  expirationDays === 30
                    ? 'border-primary/50 bg-primary/10 text-primary'
                    : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                }`}
              >
                <Clock className="h-4 w-4" /> 30 DAYS
              </button>
              <button
                type="button"
                onClick={() => setExpirationDays(90)}
                className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                  expirationDays === 90
                    ? 'border-primary/50 bg-primary/10 text-primary'
                    : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                }`}
              >
                <Clock className="h-4 w-4" /> 90 DAYS
              </button>
              <button
                type="button"
                onClick={() => setExpirationDays(null)}
                className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                  expirationDays === null
                    ? 'border-primary/50 bg-primary/10 text-primary'
                    : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                }`}
              >
                <Clock className="h-4 w-4" /> NO EXPIRATION
              </button>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 pt-2">
            <Button
              onClick={() => onOpenChange(false)}
              variant="ghost"
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              onClick={handleCreate}
              disabled={createToken.isPending}
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              <Check className="h-3.5 w-3.5" />
              {createToken.isPending ? 'Creating...' : 'Create Key'}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
