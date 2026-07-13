import { useRouterState } from '@tanstack/react-router';
import { BellIcon, PlusIcon, SearchIcon } from 'lucide-react';

const routeLabels: Record<string, { title: string; description: string }> = {
  '/': { title: 'Dashboard', description: 'Overview of your infrastructure' },
  '/projects': { title: 'Projects', description: 'Manage your projects and services' },
  '/databases': { title: 'Databases', description: 'Your databases and storage' },
  '/deployments': { title: 'Deployments', description: 'Recent and ongoing deployments' },
  '/teams': { title: 'Teams', description: 'Manage team members and access' },
  '/settings': { title: 'Settings', description: 'Account and workspace settings' },
  '/support': { title: 'Support', description: 'Get help and documentation' },
};

export function Topbar() {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const current = routeLabels[pathname] ?? { title: 'Dashboard', description: '' };

  return (
    <header className="flex h-14 shrink-0 items-center justify-between gap-4 border-b border-border bg-background/95 px-6 backdrop-blur-sm">
      <div className="flex items-center gap-3">
        <div>
          <h1 className="text-sm font-semibold text-foreground leading-none">{current.title}</h1>
          {current.description && (
            <p className="text-xs text-muted-foreground mt-0.5 leading-none">
              {current.description}
            </p>
          )}
        </div>
      </div>

      <div className="flex items-center gap-2">
        <button
          type="button"
          className="flex items-center gap-2 rounded-md border border-border bg-muted/40 px-3 py-1.5 text-sm text-muted-foreground hover:bg-muted/70 hover:text-foreground transition-colors duration-150"
        >
          <SearchIcon className="h-3.5 w-3.5" />
          <span className="hidden sm:inline text-xs">Search...</span>
          <kbd className="hidden sm:inline-flex h-4 items-center gap-0.5 rounded border border-border bg-background px-1 font-mono text-[9px] text-muted-foreground">
            ⌘K
          </kbd>
        </button>

        <button
          type="button"
          className="flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-xs font-semibold text-primary-foreground hover:bg-primary/90 transition-colors duration-150"
        >
          <PlusIcon className="h-3.5 w-3.5" />
          <span className="hidden sm:inline">New</span>
        </button>

        <button
          type="button"
          className="relative flex h-8 w-8 items-center justify-center rounded-md border border-border text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors duration-150"
        >
          <BellIcon className="h-4 w-4" />
          <span className="absolute right-1.5 top-1.5 h-1.5 w-1.5 rounded-full bg-primary" />
        </button>
      </div>
    </header>
  );
}
