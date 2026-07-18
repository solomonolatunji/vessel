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

export const ProjectCard = ({
  project,
  mode = 'grid',
}: {
  project: CanvasSummary;
  mode?: 'grid' | 'list';
}) => {
  if (mode === 'list') {
    return (
      <Link to={`/projects/$projectId`} params={{ projectId: project.id }} className="group block">
        <div className="rounded-2xl border border-border/50 bg-card/40 p-6 shadow-sm transition-all hover:border-primary/50 hover:bg-card/80">
          <div className="flex items-start justify-between">
            <div>
              <div className="flex items-center gap-3">
                <h3 className="font-bold text-foreground text-xl transition-colors group-hover:text-primary">
                  {project.name}
                </h3>
                {project.defaultEnvironment && (
                  <div className="rounded border border-primary/30 bg-primary/10 px-2 py-0.5 font-bold text-[10px] text-primary uppercase tracking-widest">
                    {project.defaultEnvironment.name}
                  </div>
                )}
              </div>
              <div className="mt-2 flex items-center gap-2 font-mono text-[10px] text-muted-foreground uppercase tracking-widest">
                {project.defaultEnvironment ? (
                  <>
                    <div className="h-1.5 w-1.5 rounded-full bg-emerald-500/80 shadow-[0_0_8px_rgba(16,185,129,0.4)]"></div>
                    <span>Production</span>
                  </>
                ) : (
                  <>
                    <div className="h-1.5 w-1.5 rounded-full bg-zinc-500/80"></div>
                    <span>No environment</span>
                  </>
                )}
                <span className="opacity-50">•</span>
                <span>
                  {project.onlineServices}/{project.totalServices} services online
                </span>
              </div>
            </div>

            {project.totalServices > 0 && (
              <div className="flex items-center gap-1.5">
                {project.serviceIcons?.slice(0, 3).map((icon, i) => (
                  <div
                    key={i}
                    className="flex h-8 w-8 items-center justify-center rounded-md border border-border/40 bg-background shadow-sm"
                  >
                    {getIcon(icon)}
                  </div>
                ))}
                {project.totalServices > 3 && (
                  <div className="flex h-8 w-8 items-center justify-center rounded-md border border-border/40 bg-background font-mono text-[10px] text-muted-foreground shadow-sm">
                    +{project.totalServices - 3}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </Link>
    );
  }

  return (
    <Link to={`/projects/$projectId`} params={{ projectId: project.id }} className="group block">
      <div className="relative overflow-hidden border border-border/40 bg-card/60 p-5 text-left shadow-sm transition-colors hover:border-primary/50 hover:bg-card">
        <div className="relative z-10">
          <div className="flex items-start justify-between gap-4">
            <div className="min-w-0">
              <h2 className="truncate font-semibold text-[15px] transition-colors group-hover:text-primary">
                {project.name}
              </h2>
            </div>
            <span className="shrink-0 border border-border/50 bg-background px-2.5 py-1 font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.18em]">
              {project.totalServices} service{project.totalServices !== 1 ? 's' : ''}
            </span>
          </div>

          <div className="mt-5">
            <div className="border border-border/40 bg-muted/20 p-2">
              <div
                className="flex min-h-[140px] items-center justify-center bg-[#0d0d10]"
                style={{
                  backgroundImage:
                    'radial-gradient(circle at 1px 1px, rgba(255,255,255,0.08) 1px, transparent 0)',
                  backgroundSize: '16px 16px',
                }}
              >
                {project.totalServices === 0 ? (
                  <div className="flex h-full min-h-[140px] items-center justify-center text-muted-foreground text-xs">
                    No services yet.
                  </div>
                ) : (
                  <div className="flex max-w-[11.5rem] flex-wrap items-center justify-center gap-2">
                    {project.serviceIcons?.slice(0, 7).map((icon, i) => (
                      <div
                        key={i}
                        className="flex h-10 w-10 items-center justify-center rounded-md border border-border/40 bg-background p-2.5 shadow-sm"
                      >
                        {getIcon(icon)}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>

          <div className="mt-5 flex items-center gap-2 text-[13px] text-muted-foreground">
            {project.defaultEnvironment && (
              <>
                <div className="h-1.5 w-1.5 rounded-full bg-emerald-500/80 shadow-[0_0_8px_rgba(16,185,129,0.4)]"></div>
                <span>{project.defaultEnvironment.name}</span>
                <span className="px-1 opacity-40">•</span>
              </>
            )}
            {!project.defaultEnvironment && project.totalServices > 0 && (
              <>
                <div className="h-1.5 w-1.5 rounded-full bg-emerald-500/80 shadow-[0_0_8px_rgba(16,185,129,0.4)]"></div>
                <span className="pr-1 opacity-40">•</span>
              </>
            )}
            <span>
              {project.onlineServices}/{project.totalServices} services online
            </span>
          </div>
        </div>
      </div>
    </Link>
  );
};
