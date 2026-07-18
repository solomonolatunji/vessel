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
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-[500px] [&>button]:hidden">
        <form onSubmit={handleSubmit}>
          <div className="flex flex-col p-8 pb-6">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="font-bold text-2xl tracking-tight">New Project</DialogTitle>
                <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  CREATE THE PROJECT FIRST, THEN ADD SERVICES INSIDE IT.
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

          <div className="space-y-6 p-8">
            <div className="space-y-3">
              <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                PROJECT NAME
              </Label>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                placeholder="Acme platform"
                className="h-12 rounded-xl border-border/50 bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              />
            </div>
            <div className="space-y-3">
              <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                DESCRIPTION
              </Label>
              <Input
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Internal tools and APIs"
                className="h-12 rounded-xl border-border/50 bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
              />
            </div>
          </div>

          <div className="flex items-center justify-end gap-6 p-8 pt-6">
            <Button type="button" variant="ghost" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending ? 'CREATING...' : 'CREATE PROJECT'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
