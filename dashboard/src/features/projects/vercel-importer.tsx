import { Check, Loader2, Triangle } from 'lucide-react';
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
import { useCreateProject, useSetVars } from '#/hooks/useProjects';
import { useGetProjectEnv, useListProjects as useListVercelProjects } from '#/hooks/useVercel';
import type { VercelProject } from '#/interfaces/vercel';

export function VercelImporter({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [step, setStep] = useState<'projects' | 'options'>('projects');
  const [selectedProjectId, setSelectedProjectId] = useState<string>('');
  const [importEnvVars, setImportEnvVars] = useState(true);

  const { data: projectsResponse, isLoading: isLoadingProjects } = useListVercelProjects(open);
  const projectsData = projectsResponse?.data || [];

  const { refetch: fetchEnvVars } = useGetProjectEnv(selectedProjectId, false);
  const { mutateAsync: createProject, isPending: isCreating } = useCreateProject();
  const { mutateAsync: setVars, isPending: isSettingVars } = useSetVars();

  const isWorking = isCreating || isSettingVars;

  const handleNext = async (e: React.FormEvent) => {
    e.preventDefault();
    if (step === 'projects') {
      if (!selectedProjectId) return toast.error('Select a project');
      setStep('options');
    } else if (step === 'options') {
      try {
        const project = projectsData.find((p) => p.id === selectedProjectId);
        if (!project) throw new Error('Project not found');

        const newProject = await createProject({
          payload: { name: project.name, description: `Imported from Vercel: ${project.name}` },
        });
        const createdProjectId = newProject.data.id;

        if (importEnvVars) {
          const envResponse = await fetchEnvVars();
          const envVars = envResponse.data?.data || [];
          if (envVars.length > 0) {
            const variables: Record<string, string> = {};
            for (const env of envVars) {
              if (env.type === 'plain' || env.type === 'encrypted') {
                variables[env.key] = env.value || '';
              }
            }
            await setVars({ id: createdProjectId, payload: { variables } });
          }
        }
        toast.success('Vercel project imported successfully!');
        onOpenChange(false);
        setStep('projects');
      } catch {
        toast.error('Failed to import Vercel project');
      }
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-[500px] [&>button]:hidden">
        <form onSubmit={handleNext}>
          <div className="flex flex-col p-8 pb-6">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="font-bold text-2xl tracking-tight">
                  Import from Vercel
                </DialogTitle>
                <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  {step === 'projects'
                    ? 'SELECT A PROJECT TO IMPORT.'
                    : 'CONFIGURE IMPORT OPTIONS.'}
                </DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => setStep('projects')}
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="h-px w-full bg-border/50" />

          <div className="p-8">
            {step === 'projects' && (
              <div className="max-h-[300px] space-y-4 overflow-y-auto pr-2">
                {isLoadingProjects ? (
                  <div className="flex justify-center p-4">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  </div>
                ) : projectsData.length > 0 ? (
                  projectsData.map((project: VercelProject) => (
                    <button
                      type="button"
                      key={project.id}
                      onClick={() => setSelectedProjectId(project.id)}
                      className={`w-full cursor-pointer rounded-lg border p-3 text-left transition-colors ${
                        selectedProjectId === project.id
                          ? 'border-zinc-400 bg-zinc-400/10'
                          : 'border-border/50 hover:bg-muted/50'
                      }`}
                    >
                      <div className="font-medium text-sm">{project.name}</div>
                      {project.framework && (
                        <div className="mt-1 text-muted-foreground text-xs">
                          Framework: {project.framework}
                        </div>
                      )}
                    </button>
                  ))
                ) : (
                  <div className="p-4 text-center text-muted-foreground text-sm">
                    No Vercel projects found. Make sure your account is linked.
                  </div>
                )}
              </div>
            )}

            {step === 'options' && (
              <div className="space-y-4">
                <label className="flex cursor-pointer items-start gap-3 rounded-lg border border-border/50 p-3 hover:bg-muted/50">
                  <input
                    type="checkbox"
                    className="mt-1"
                    checked={importEnvVars}
                    onChange={(e) => setImportEnvVars(e.target.checked)}
                  />
                  <div>
                    <div className="font-medium text-sm">Import Environment Variables</div>
                    <div className="mt-0.5 text-muted-foreground text-xs">
                      Copy all environment variables from Vercel to Vessl
                    </div>
                  </div>
                </label>
              </div>
            )}
          </div>

          <div className="flex items-center justify-end gap-6 p-8 pt-6">
            {step !== 'projects' && (
              <Button
                type="button"
                variant="ghost"
                onClick={() => setStep('projects')}
                className="mr-auto"
              >
                BACK
              </Button>
            )}
            <Button type="button" variant="ghost" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isWorking || isLoadingProjects}>
              {isWorking && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {step === 'options' ? 'START IMPORT' : 'CONTINUE'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
