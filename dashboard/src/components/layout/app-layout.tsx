import { useRouterState } from '@tanstack/react-router';
import type * as React from 'react';
import { useState } from 'react';
import { AppSidebar } from './app-sidebar';
import { CommandPalette } from './command-palette';
import { Topbar } from './topbar';

export function AppLayout({ children }: { children: React.ReactNode }) {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [commandOpen, setCommandOpen] = useState(false);
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });

  return (
    <div className="relative flex min-h-screen bg-background">
      <div className="pointer-events-none fixed inset-0 bg-[radial-gradient(#6d28d9_0.8px,transparent_1px)] bg-size-[40px_40px] opacity-10 dark:opacity-20" />
      <AppSidebar collapsed={sidebarCollapsed} onToggle={() => setSidebarCollapsed((p) => !p)} />
      <div className={`relative flex flex-1 flex-col ${sidebarCollapsed ? 'pl-16' : 'pl-64'}`}>
        <Topbar onOpenCommand={() => setCommandOpen(true)} />
        <main className="flex-1 overflow-auto p-8">
          <div key={pathname} className="page-transition mx-auto w-full max-w-[1280px]">
            {children}
          </div>
        </main>
      </div>
      <CommandPalette open={commandOpen} onOpenChange={setCommandOpen} />
    </div>
  );
}
