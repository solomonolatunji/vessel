import { Loader2 } from 'lucide-react';
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
import { useGetRailwayProjects, useImportRailwayProject } from '#/hooks/useSystem';
import type { RailwayProject } from '#/interfaces/system';

export function RailwayImporter({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [step, setStep] = useState<'token' | 'projects' | 'options'>('token');
  const [token, setToken] = useState('');
  const [selectedProjectId, setSelectedProjectId] = useState<string>('');

  const [excludeRailwayVars, setExcludeRailwayVars] = useState(false);
  const [recreateDatabases, setRecreateDatabases] = useState(true);
  const [importData, setImportData] = useState(false);

  const {
    data: projectsData,
    isLoading: isLoadingProjects,
    refetch: fetchProjects,
  } = useGetRailwayProjects(token);
  const { mutateAsync: importProject, isPending: isImporting } = useImportRailwayProject();

  const handleNext = async (e: React.FormEvent) => {
    e.preventDefault();
    if (step === 'token') {
      if (!token) return toast.error('Token required');
      await fetchProjects();
      setStep('projects');
    } else if (step === 'projects') {
      if (!selectedProjectId) return toast.error('Select a project');
      setStep('options');
    } else if (step === 'options') {
      try {
        await importProject({
          token,
          projectId: selectedProjectId,
          excludeRailwayVars,
          recreateDatabases,
          importData,
        });
        toast.success('Railway project import started!');
        onOpenChange(false);
        setStep('token');
        setToken('');
        setSelectedProjectId('');
      } catch {
        toast.error('Failed to import project');
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
                  Import from Railway
                </DialogTitle>
                <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  {step === 'token' && 'AUTHENTICATE WITH YOUR RAILWAY TOKEN.'}
                  {step === 'projects' && 'SELECT A PROJECT TO IMPORT.'}
                  {step === 'options' && 'CONFIGURE IMPORT OPTIONS.'}
                </DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => {
                    if (step === 'options') setStep('projects');
                    else if (step === 'projects') setStep('token');
                  }}
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="h-px w-full bg-border/50" />

          <div className="p-8">
            {step === 'token' && (
              <div className="space-y-3">
                <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  RAILWAY API TOKEN
                </Label>
                <Input
                  type="password"
                  value={token}
                  onChange={(e) => setToken(e.target.value)}
                  placeholder="Paste your Railway project token..."
                  required
                  className="h-12 rounded-xl border-border/50 bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
                />
              </div>
            )}

            {step === 'projects' && (
              <div className="max-h-[300px] space-y-4 overflow-y-auto pr-2">
                {isLoadingProjects ? (
                  <div className="flex justify-center p-4">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  </div>
                ) : projectsData && projectsData.length > 0 ? (
                  projectsData.map((project: RailwayProject) => (
                    <button
                      type="button"
                      key={project.id}
                      onClick={() => setSelectedProjectId(project.id)}
                      className={`w-full cursor-pointer rounded-lg border p-3 text-left transition-colors ${
                        selectedProjectId === project.id
                          ? 'border-[#E93D82] bg-[#E93D82]/5'
                          : 'border-border/50 hover:bg-muted/50'
                      }`}
                    >
                      <div className="font-medium text-sm">{project.name}</div>
                      <div className="mt-1 text-muted-foreground text-xs">
                        {project.description || 'No description'}
                      </div>
                    </button>
                  ))
                ) : (
                  <div className="p-4 text-center text-muted-foreground text-sm">
                    No projects found.
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
                    checked={!excludeRailwayVars}
                    onChange={(e) => setExcludeRailwayVars(!e.target.checked)}
                  />
                  <div>
                    <div className="font-medium text-sm">Include Railway Variables</div>
                    <div className="mt-0.5 text-muted-foreground text-xs">
                      Import RAILWAY_* environment variables
                    </div>
                  </div>
                </label>

                <label className="flex cursor-pointer items-start gap-3 rounded-lg border border-border/50 p-3 hover:bg-muted/50">
                  <input
                    type="checkbox"
                    className="mt-1"
                    checked={recreateDatabases}
                    onChange={(e) => setRecreateDatabases(e.target.checked)}
                  />
                  <div>
                    <div className="font-medium text-sm">Recreate Databases</div>
                    <div className="mt-0.5 text-muted-foreground text-xs">
                      Provision managed databases for Railway plugins
                    </div>
                  </div>
                </label>

                <label className="flex cursor-pointer items-start gap-3 rounded-lg border border-border/50 p-3 hover:bg-muted/50">
                  <input
                    type="checkbox"
                    className="mt-1"
                    checked={importData}
                    onChange={(e) => setImportData(e.target.checked)}
                    disabled={!recreateDatabases}
                  />
                  <div className={!recreateDatabases ? 'opacity-50' : ''}>
                    <div className="font-medium text-sm">Import Data (Coming Soon)</div>
                    <div className="mt-0.5 text-muted-foreground text-xs">
                      Migrate data from Railway plugins
                    </div>
                  </div>
                </label>
              </div>
            )}
          </div>

          <div className="flex items-center justify-end gap-6 p-8 pt-6">
            {step !== 'token' && (
              <Button
                type="button"
                variant="ghost"
                onClick={() => (step === 'options' ? setStep('projects') : setStep('token'))}
                className="mr-auto"
              >
                BACK
              </Button>
            )}
            <Button type="button" variant="ghost" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isImporting || isLoadingProjects}>
              {isImporting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {step === 'options' ? 'START IMPORT' : 'CONTINUE'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
