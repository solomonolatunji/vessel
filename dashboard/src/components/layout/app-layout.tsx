import { useRouterState } from '@tanstack/react-router';
import { Menu, X } from 'lucide-react';
import type * as React from 'react';
import { useEffect, useState } from 'react';
import { AppSidebar } from './app-sidebar';
import { BackgroundPattern } from './background-pattern';
import { CommandPalette } from './command-palette';

export function AppLayout({ children }: { children: React.ReactNode }) {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [commandOpen, setCommandOpen] = useState(false);
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });

  // biome-ignore lint/correctness/useExhaustiveDependencies: run on pathname change
  useEffect(() => {
    setMobileMenuOpen(false);
  }, [pathname]);

  return (
    <div className="relative flex min-h-screen bg-background">
      <BackgroundPattern />
      <AppSidebar
        collapsed={sidebarCollapsed}
        onToggle={() => setSidebarCollapsed((p) => !p)}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />
      <div
        className={`relative flex flex-1 flex-col ${sidebarCollapsed ? 'md:pl-16' : 'md:pl-64'}`}
      >
        <div className="flex h-14 items-center px-4 md:hidden">
          <button
            type="button"
            onClick={() => setMobileMenuOpen((p) => !p)}
            className="flex h-9 w-9 items-center justify-center rounded-xl border border-border/60 text-muted-foreground transition-colors hover:bg-muted"
          >
            {mobileMenuOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
          </button>
        </div>
        <main className="flex-1 overflow-auto p-4 md:p-8 md:pt-12">
          <div key={pathname} className="page-transition mx-auto w-full max-w-7xl">
            {children}
          </div>
        </main>
      </div>
      <CommandPalette open={commandOpen} onOpenChange={setCommandOpen} />
    </div>
  );
}
