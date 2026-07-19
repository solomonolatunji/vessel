import { createFileRoute } from '@tanstack/react-router';
import { Box, LayoutGrid, Loader2, LogOut, Plus, SearchIcon } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { CreateProjectModal } from '#/features/projects/create-project-modal';
import { ProjectCard } from '#/features/projects/project-card';
import { useListCanvasSummaries } from '#/hooks/useCanvas';

export const Route = createFileRoute('/_dashboard/')({
  component: DashboardIndex,
});

function DashboardIndex() {
  const [createOpen, setCreateOpen] = useState(false);
  const { data, isLoading } = useListCanvasSummaries();

  const projects = data?.data || [];

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Box className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Projects</h1>
            <p className="text-muted-foreground text-sm">Manage your Vessl projects.</p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            className="gap-2"
            onClick={() => {
              window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', metaKey: true }));
            }}
          >
            <SearchIcon className="h-4 w-4" />
            SEARCH
          </Button>

          <Button onClick={() => setCreateOpen(true)} className="gap-2">
            <Plus className="h-4 w-4" />
            NEW PROJECT
          </Button>

          <Button variant="destructive" size="icon">
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : projects.length === 0 ? (
        <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card/40">
          <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
            <LayoutGrid className="h-5 w-5 text-primary" />
          </div>
          <h3 className="mt-4 font-bold text-foreground text-lg tracking-tight">No projects yet</h3>
          <p className="mt-1 max-w-sm text-center text-muted-foreground text-sm">
            Create a project first, then attach services inside it.
          </p>
          <Button className="mt-6 gap-2" onClick={() => setCreateOpen(true)}>
            <Plus className="h-4 w-4" />
            CREATE PROJECT
          </Button>
        </div>
      ) : (
        <div className="space-y-6">
          <div className="flex items-center justify-between border-border/20 border-b pb-4">
            <div className="flex items-center gap-4 text-sm">
              <div className="flex items-center gap-2 font-medium">
                <LayoutGrid className="h-4 w-4 text-muted-foreground" />
                {projects.length} Project{projects.length !== 1 && 's'}
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3">
            {projects.map((project) => (
              <div key={project.id}>
                <ProjectCard project={project} />
              </div>
            ))}
          </div>
        </div>
      )}

      <CreateProjectModal open={createOpen} onOpenChange={setCreateOpen} />
    </div>
  );
}
