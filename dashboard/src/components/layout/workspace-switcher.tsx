import { Link } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import { CheckCircle, ChevronDown, Plus } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useListWorkspaces } from '#/hooks/useWorkspaces';
import type { Workspace } from '#/interfaces/workspace';
import { workspaceActions, workspaceStore } from '#/stores/workspaceStore';

export function WorkspaceSwitcher() {
  const { data } = useListWorkspaces();
  const workspaceState = useStore(workspaceStore);
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const workspaces: Workspace[] = useMemo(
    () => (data as { data?: Workspace[] } | undefined)?.data ?? [],
    [data]
  );
  const active = workspaceState.activeWorkspace;

  useEffect(() => {
    if (workspaces.length > 0) {
      workspaceActions.setWorkspaces(workspaces);
    }
  }, [workspaces]);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  const initials = active?.name ? active.name.slice(0, 2).toUpperCase() : '??';

  return (
    <div ref={ref} className="relative px-3 pb-2">
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        className="flex w-full items-center gap-2.5 rounded-lg border border-sidebar-border bg-sidebar-accent/40 px-2.5 py-2 text-left hover:bg-sidebar-accent/70 transition-colors duration-100"
      >
        <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-primary/20 text-[10px] font-bold text-primary">
          {initials}
        </div>
        <div className="min-w-0 flex-1">
          <p className="truncate text-[12px] font-semibold text-sidebar-foreground leading-none">
            {active?.name ?? 'Select workspace'}
          </p>
          {active?.preferredRegion && (
            <p className="truncate text-[10px] text-sidebar-foreground/40 mt-0.5 leading-none">
              {active.preferredRegion}
            </p>
          )}
        </div>
        <ChevronDown
          className={[
            'h-3.5 w-3.5 shrink-0 text-sidebar-foreground/40 transition-transform duration-150',
            open ? 'rotate-180' : '',
          ].join(' ')}
        />
      </button>

      {open && (
        <div className="absolute left-3 right-3 top-full z-50 mt-1 rounded-lg border border-border bg-popover shadow-lg overflow-hidden">
          <div className="p-1 max-h-48 overflow-y-auto">
            {workspaces.length === 0 && (
              <p className="px-2 py-3 text-center text-xs text-muted-foreground">
                No workspaces yet
              </p>
            )}
            {workspaces.map((ws) => (
              <button
                key={ws.id}
                type="button"
                onClick={() => {
                  workspaceActions.switchWorkspace(ws);
                  setOpen(false);
                }}
                className="flex w-full items-center gap-2.5 rounded-md px-2 py-1.5 text-sm hover:bg-accent transition-colors"
              >
                <div className="flex h-5 w-5 shrink-0 items-center justify-center rounded text-[9px] font-bold bg-primary/15 text-primary">
                  {ws.name.slice(0, 2).toUpperCase()}
                </div>
                <span className="flex-1 truncate text-[12px] font-medium text-foreground">
                  {ws.name}
                </span>
                {ws.id === active?.id && (
                  <CheckCircle className="h-3.5 w-3.5 text-primary shrink-0" />
                )}
              </button>
            ))}
          </div>
          <div className="border-t border-border p-1">
            <Link
              to={'/workspaces' as never}
              onClick={() => setOpen(false)}
              className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-xs font-medium text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            >
              <Plus className="h-3.5 w-3.5" />
              New workspace
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}
