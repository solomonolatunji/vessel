import { Link } from '@tanstack/react-router';
import { ChevronDown, Plus, Server } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { useServers } from '#/features/workspaces/hooks/useServers';
import { serverActions, useServerState } from '#/stores/serverStore';

export function ServerSwitcher() {
  const serverState = useServerState();
  const { data: serversData } = useServers();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (serversData && serversData.length > 0) {
      serverActions.setServers(serversData);
    }
  }, [serversData]);

  const servers = serverState.servers;
  const active = serverState.activeServer;

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  return (
    <div ref={ref} className="relative px-3 pb-2">
      <div className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground mb-1.5 px-1">
        Active Server
      </div>
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        className="flex w-full items-center gap-2.5 rounded-lg border border-sidebar-border bg-sidebar-accent/40 px-2.5 py-2 text-left hover:bg-sidebar-accent/70 transition-colors duration-100"
      >
        <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-emerald-500/10 text-[10px] font-bold text-emerald-500">
          <Server className="w-3.5 h-3.5" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="truncate text-[12px] font-semibold text-sidebar-foreground leading-none">
            {active?.name ?? 'Select server'}
          </p>
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
            {servers.length === 0 && (
              <p className="px-2 py-3 text-center text-xs text-muted-foreground">
                No servers connected
              </p>
            )}
            {servers.map((srv) => (
              <button
                type="button"
                key={srv.id}
                onClick={() => {
                  serverActions.switchServer(srv);
                  setOpen(false);
                }}
                className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-muted transition-colors text-left"
              >
                <Server className="w-3.5 h-3.5 text-emerald-500 shrink-0" />
                <span className="truncate flex-1">{srv.name}</span>
                {active?.id === srv.id && (
                  <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0" />
                )}
              </button>
            ))}
          </div>
          <div className="border-t border-border p-1">
            <Link
              to="/settings"
              onClick={() => setOpen(false)}
              className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-muted text-muted-foreground transition-colors"
            >
              <Plus className="h-4 w-4" />
              <span>Connect new server</span>
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}
