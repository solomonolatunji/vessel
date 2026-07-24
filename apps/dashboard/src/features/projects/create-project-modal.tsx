import { FolderKanban } from 'lucide-react';
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
import { useCreateProject } from '#/hooks/useProjects';

export function CreateProjectModal({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const { mutateAsync: createProject, isPending } = useCreateProject();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createProject({ payload: { name, description } });
      toast.success('Project created');
      onOpenChange(false);
      setName('');
      setDescription('');
    } catch {
      toast.error('Failed to create project');
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-125 [&>button]:hidden">
        <form onSubmit={handleSubmit}>
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-xl tracking-tight">
                  <FolderKanban className="h-5 w-5 text-primary" />
                  New Project
                </DialogTitle>
                <DialogDescription>Create the project first, then add services</DialogDescription>
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

          <div className="space-y-5 px-5 pt-4 pb-5">
            <div className="space-y-2.5">
              <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                PROJECT NAME
              </Label>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                placeholder="Acme platform"
                className="h-10 rounded-lg border-border/50 bg-background/80 px-3 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-1 focus:ring-primary/20"
              />
            </div>
            <div className="space-y-2.5">
              <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                DESCRIPTION
              </Label>
              <Input
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Internal tools and APIs"
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
              disabled={isPending}
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              <FolderKanban className="h-3.5 w-3.5" />
              {isPending ? 'Creating...' : 'Create Project'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
