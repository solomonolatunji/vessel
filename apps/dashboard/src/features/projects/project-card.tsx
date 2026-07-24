import { Link } from '@tanstack/react-router';
import { Box, Cloud, Database, Folder } from 'lucide-react';

import type { CanvasSummary } from '#/interfaces/project';

const GithubIcon = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.2c3-.3 6-1.5 6-6.5 0-1.4-.5-2.5-1.5-3.4.1-.3.6-1.6-.1-3.3 0 0-1.2-.4-3.8 1.4a12.8 12.8 0 0 0-7 0C3.9 1.5 2.7 1.9 2.7 1.9c-.7 1.7-.2 3 .1 3.3-1 1-1.5 2-1.5 3.4 0 5 3 6.2 6 6.5-.4.4-.7 1-.8 2.2-.8.4-2.8.9-4-1.1 0 0-.7-1.3-2-1.4 0 0-1.3-.1-.1 1.2 0 0 1.2 1.8 3 2.5 1.5.5 3.3.4 3.3.4z" />
  </svg>
);

// Basic mapping of icons based on string value
const IconMap: Record<string, React.ReactNode> = {
  github: <GithubIcon className="h-5 w-5" />,
  postgres: <Database className="h-5 w-5 text-blue-500" />,
  mysql: <Database className="h-5 w-5 text-blue-400" />,
  redis: <Database className="h-5 w-5 text-red-500" />,
  s3: <Cloud className="h-5 w-5 text-amber-500" />,
  local: <Folder className="h-5 w-5 text-gray-500" />,
};

const getIcon = (iconName: string) => {
  return IconMap[iconName.toLowerCase()] || <Box className="h-5 w-5 text-primary" />;
};

export const ProjectCard = ({ project }: { project: CanvasSummary }) => {
  return (
    <Link
      to={`/projects/$projectId`}
      params={{ projectId: project.id }}
      className="group block h-full"
    >
      <div className="flex h-full flex-col rounded-2xl border border-border/50 bg-card/40 p-5 transition-colors hover:border-primary/50 hover:bg-card/80">
        <div className="flex items-start justify-between">
          <div>
            <h3 className="font-bold text-foreground text-lg transition-colors group-hover:text-primary">
              {project.name}
            </h3>
            <div className="mt-1 flex items-center gap-2 font-mono text-[10px] text-muted-foreground uppercase tracking-widest">
              {project.defaultEnvironment ? (
                <>
                  <div className="h-1.5 w-1.5 rounded-full bg-emerald-500/80 shadow-[0_0_8px_rgba(16,185,129,0.4)]"></div>
                  <span>{project.defaultEnvironment.name}</span>
                </>
              ) : (
                <>
                  <div className="h-1.5 w-1.5 rounded-full bg-zinc-500/80"></div>
                  <span>No Environment</span>
                </>
              )}
            </div>
          </div>
          <div className="rounded border border-primary/30 bg-primary/10 px-2 py-0.5 font-bold text-[10px] text-primary uppercase tracking-widest">
            {project.onlineServices}/{project.totalServices} ONLINE
          </div>
        </div>

        <div className="mt-4 flex-1 border-border/50 border-t pt-4">
          {project.totalServices === 0 ? (
            <div className="flex h-18 items-center justify-center rounded-xl border border-border/50 border-dashed bg-background/30">
              <span className="font-mono text-[10px] text-muted-foreground uppercase tracking-widest">
                No services attached
              </span>
            </div>
          ) : (
            <div className="flex flex-col gap-3">
              <span className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                Attached Services
              </span>
              <div className="flex flex-wrap gap-2">
                {project.serviceIcons?.slice(0, 5).map((icon, i) => (
                  <div
                    key={i}
                    className="flex h-10 w-10 items-center justify-center rounded-lg border border-border/50 bg-background/50 shadow-sm transition-colors group-hover:border-primary/30"
                  >
                    {getIcon(icon)}
                  </div>
                ))}
                {project.totalServices > 5 && (
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-border/50 bg-background/50 font-mono text-[10px] text-muted-foreground shadow-sm">
                    +{project.totalServices - 5}
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </Link>
  );
};
