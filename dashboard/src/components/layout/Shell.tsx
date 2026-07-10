import {
  Activity01Icon,
  CloudServerIcon,
  DashboardSquare01Icon,
  Database01Icon,
  Folder01Icon,
  Settings01Icon,
  UserGroupIcon,
} from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { Link } from '@tanstack/react-router';
import type * as React from 'react';

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from '#/components/ui/sidebar.tsx';

// Menu items
const items = [
  {
    title: 'Dashboard',
    url: '/',
    icon: DashboardSquare01Icon,
  },
  {
    title: 'Projects',
    url: '/projects',
    icon: Folder01Icon,
  },
  {
    title: 'Databases & Storage',
    url: '/databases',
    icon: Database01Icon,
  },
  {
    title: 'Jobs & Backups',
    url: '/jobs',
    icon: Activity01Icon,
  },
  {
    title: 'Teams',
    url: '/teams',
    icon: UserGroupIcon,
  },
  {
    title: 'Settings',
    url: '/settings',
    icon: Settings01Icon,
  },
];

export function Shell({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider>
      <Sidebar
        variant="inset"
        className="dark bg-zinc-950/50 backdrop-blur-xl border-r border-zinc-800"
      >
        <SidebarHeader className="h-16 flex items-center justify-center border-b border-zinc-800/50">
          <div className="flex items-center gap-2 font-bold text-lg tracking-tight text-zinc-100 w-full px-4">
            <HugeiconsIcon icon={CloudServerIcon} className="h-5 w-5 text-indigo-400" />
            <span>Vessel</span>
          </div>
        </SidebarHeader>
        <SidebarContent className="bg-transparent">
          <SidebarGroup>
            <SidebarGroupLabel className="text-zinc-400 font-medium tracking-wider">
              Platform
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      asChild
                      tooltip={item.title}
                      className="hover:bg-zinc-800/50 text-zinc-300 hover:text-zinc-50 transition-colors"
                    >
                      <Link to={item.url}>
                        <HugeiconsIcon icon={item.icon} className="h-4 w-4" />
                        <span>{item.title}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter className="border-t border-zinc-800/50 p-4">
          <div className="text-xs text-zinc-500 text-center">Vessel Control Panel v0.1.0</div>
        </SidebarFooter>
      </Sidebar>

      <SidebarInset className="bg-zinc-950 flex flex-col min-h-screen">
        <header className="h-16 shrink-0 border-b border-zinc-800 flex items-center justify-between px-4 sticky top-0 bg-zinc-950/80 backdrop-blur-md z-10">
          <div className="flex items-center gap-4">
            <SidebarTrigger className="text-zinc-400 hover:text-zinc-50" />
          </div>

          {/* System Health Indicators */}
          <div className="flex items-center gap-4 text-sm font-medium">
            <div className="flex items-center gap-2 text-zinc-400">
              <span className="flex h-2 w-2 rounded-full bg-emerald-500 ring-2 ring-emerald-500/20"></span>
              Docker: Online
            </div>
            <div className="h-4 w-px bg-zinc-800"></div>
            <div className="flex items-center gap-2 text-zinc-400">
              <span>
                CPU: <span className="text-zinc-200">12%</span>
              </span>
            </div>
            <div className="h-4 w-px bg-zinc-800"></div>
            <div className="flex items-center gap-2 text-zinc-400">
              <span>
                RAM: <span className="text-zinc-200">4.2GB</span>
              </span>
            </div>
            <div className="h-4 w-px bg-zinc-800"></div>
            <div className="bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 px-2 py-1 rounded-md text-xs">
              v0.1.0 Available
            </div>
          </div>
        </header>

        <main className="flex-1 overflow-auto bg-zinc-950 text-zinc-50 p-6">{children}</main>
      </SidebarInset>
    </SidebarProvider>
  );
}
