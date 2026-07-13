import type * as React from 'react';

import { SidebarInset, SidebarProvider } from '#/components/ui/sidebar';
import { AppSidebar } from './app-sidebar';
import { Topbar } from './topbar';

export function Shell({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset className="bg-background flex flex-col min-h-screen">
        <Topbar />
        <main className="flex-1 overflow-auto bg-background p-6">
          <div className="mx-auto max-w-7xl w-full">{children}</div>
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
