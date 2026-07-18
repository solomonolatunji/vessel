import { LayoutGrid, List, Loader2, Plus } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { useListCanvasSummaries } from '#/hooks/useCanvas';
import { ProjectCard } from './project-card';

export const ProjectList = ({ onCreateClick }: { onCreateClick?: () => void }) => {
  const { data, isLoading } = useListCanvasSummaries();
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const projects = data?.data || [];

  if (projects.length === 0) {
    return (
      <section className="border border-border/40 bg-card/40 px-6 py-10 sm:px-8">
        <div className="max-w-2xl">
          <div className="inline-flex items-center gap-2 border border-border/50 bg-muted/30 px-3 py-1 font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
            <LayoutGrid className="h-3.5 w-3.5" />
            Empty registry
          </div>
          <h2 className="mt-6 font-bold text-3xl text-foreground tracking-tight">
            No projects yet
          </h2>
          <p className="mt-3 max-w-lg font-mono text-muted-foreground text-sm leading-relaxed">
            Create a project first, then attach services inside it. Each service gets its own
            deployment timeline, runtime logs, variables, and domains.
          </p>
          <Button
            className="mt-8 h-auto gap-2 rounded-none border border-primary/20 bg-primary/10 py-2.5 font-mono font-semibold text-[11px] text-primary uppercase tracking-wider hover:bg-primary/20"
            onClick={onCreateClick}
          >
            <Plus className="h-4 w-4" />
            Create project
          </Button>
        </div>
      </section>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between border-border/20 border-b pb-4">
        <div className="flex items-center gap-4 text-sm">
          <div className="flex items-center gap-2 font-medium">
            <LayoutGrid className="h-4 w-4 text-muted-foreground" />
            {projects.length} Project{projects.length !== 1 && 's'}
          </div>
          <div className="h-4 w-[1px] bg-border/60"></div>
          <div className="flex cursor-pointer items-center gap-1 text-muted-foreground transition-colors hover:text-foreground">
            Sort By: <span className="text-foreground">Recent Activity</span>
            <span className="ml-0.5 text-[10px] opacity-60">▼</span>
          </div>
        </div>

        <div className="flex items-center rounded-lg border border-border/50 bg-muted/30 p-1">
          <Button
            variant="ghost"
            size="sm"
            className={`h-7 rounded-md px-2.5 ${viewMode === 'grid' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:bg-transparent hover:text-foreground'}`}
            onClick={() => setViewMode('grid')}
          >
            <LayoutGrid className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className={`h-7 rounded-md px-2.5 ${viewMode === 'list' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:bg-transparent hover:text-foreground'}`}
            onClick={() => setViewMode('list')}
          >
            <List className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div
        className={
          viewMode === 'grid'
            ? 'grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3'
            : 'flex flex-col gap-4'
        }
      >
        {projects.map((project) => (
          <div key={project.id}>
            <ProjectCard project={project} mode={viewMode} />
          </div>
        ))}
      </div>
    </div>
  );
};
