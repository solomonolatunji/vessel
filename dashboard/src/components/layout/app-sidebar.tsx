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

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '#/components/ui/sidebar';

// Menu items
const items = [
  { title: 'Dashboard', url: '/', icon: DashboardSquare01Icon },
  { title: 'Projects', url: '/projects', icon: Folder01Icon },
  { title: 'Databases & Storage', url: '/databases', icon: Database01Icon },
  { title: 'Jobs & Backups', url: '/jobs', icon: Activity01Icon },
  { title: 'Teams', url: '/teams', icon: UserGroupIcon },
  { title: 'Settings', url: '/settings', icon: Settings01Icon },
];

export function AppSidebar() {
  return (
    <Sidebar variant="inset">
      <SidebarHeader className="h-16 flex items-center justify-center border-b">
        <div className="flex items-center gap-2 font-bold text-lg tracking-tight text-foreground w-full px-4">
          <HugeiconsIcon icon={CloudServerIcon} className="h-5 w-5 text-primary" />
          <span>Vessl</span>
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel className="text-muted-foreground font-medium tracking-wider">
            Platform
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild tooltip={item.title}>
                    <Link
                      to={item.url}
                      activeProps={{
                        className: 'bg-sidebar-accent text-sidebar-accent-foreground font-medium',
                      }}
                      inactiveProps={{
                        className:
                          'text-sidebar-foreground/70 hover:text-sidebar-foreground hover:bg-sidebar-accent/50',
                      }}
                    >
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

      <SidebarFooter className="border-t p-4">
        <div className="text-xs text-muted-foreground text-center">Vessl Control Panel v0.1.0</div>
      </SidebarFooter>
    </Sidebar>
  );
}
