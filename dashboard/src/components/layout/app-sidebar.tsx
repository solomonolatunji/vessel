import {
  Bot,
  Cloud,
  Code,
  Download,
  Globe,
  HardDrive,
  Key,
  LayoutDashboard,
  LayoutTemplate,
  Network,
  PanelLeft,
  RefreshCw,
  ScrollText,
  Settings,
  Users,
  Wrench,
  X,
} from 'lucide-react';
import { NavItem, type NavItemProps } from './nav-item';
import { UserMenu } from './user-menu';

type NavGroup = {
  title?: string;
  items: (NavItemProps & { exact?: boolean })[];
};

const navGroups: NavGroup[] = [
  {
    title: 'Overview',
    items: [
      { title: 'Dashboard', url: '/', icon: LayoutDashboard, exact: true },
      { title: 'Users', url: '/users', icon: Users },
    ],
  },
  {
    title: 'Resources',
    items: [
      { title: 'Storage', url: '/storage', icon: HardDrive },
      { title: 'S3/R2 Destinations', url: '/s3-destinations', icon: Cloud },
      { title: 'Domains', url: '/domains', icon: Globe },
    ],
  },
  {
    title: 'Discover',
    items: [
      { title: 'Templates', url: '/templates', icon: LayoutTemplate },
      { title: 'Sources', url: '/sources', icon: Code },
      { title: 'AI', url: '/ai', icon: Bot },
    ],
  },
  {
    title: 'System & Settings',
    items: [
      { title: 'API Access', url: '/api-access', icon: Key },
      { title: 'DNS', url: '/dns', icon: Network },
      { title: 'Migration', url: '/migrations', icon: Download },
      { title: 'Maintenance', url: '/maintenance', icon: Wrench },
      { title: 'Updates', url: '/updates', icon: RefreshCw },
      { title: 'Settings', url: '/settings', icon: Settings, exact: true },
    ],
  },
];

const bottomNav = [
  {
    title: 'Docs',
    url: 'https://docs.vessl.com',
    icon: ScrollText,
    external: true,
  },
];

interface AppSidebarProps {
  collapsed: boolean;
  onToggle: () => void;
  mobileOpen: boolean;
  onMobileClose: () => void;
}

export function AppSidebar({ collapsed, onToggle, mobileOpen, onMobileClose }: AppSidebarProps) {
  const navCollapsed = collapsed && !mobileOpen;

  return (
    <>
      {mobileOpen && (
        <button
          type="button"
          className="fixed inset-0 z-30 cursor-default bg-black/50 md:hidden"
          onClick={onMobileClose}
          aria-label="Close menu"
        />
      )}

      <aside
        className={`fixed inset-y-0 left-0 z-40 flex flex-col border-sidebar-border/50 border-r bg-sidebar transition-all duration-300 md:z-20 ${
          collapsed ? 'md:w-16' : 'md:w-64'
        } ${mobileOpen ? 'w-64 translate-x-0' : 'w-64 -translate-x-full md:translate-x-0'}`}
      >
        <div className="flex items-center justify-between px-2 py-2 md:hidden">
          <div className="flex items-center gap-3 px-2.5 py-2">
            <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg border border-transparent text-muted-foreground transition-all duration-150">
              <Cloud className="h-4 w-4 text-primary" />
            </div>
            <span className="truncate font-medium text-sidebar-foreground text-sm">Vessl</span>
          </div>
          <button
            type="button"
            onClick={onMobileClose}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="hidden md:block">
          <div className={collapsed ? 'px-2 py-2' : 'px-2'}>
            <div
              className={`flex items-center ${collapsed ? 'justify-center' : 'gap-3 px-2.5 py-2'}`}
            >
              <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg border border-transparent text-muted-foreground transition-all duration-150">
                <Cloud className="h-4 w-4 text-primary" />
              </div>
              {!collapsed && (
                <>
                  <span className="flex-1 truncate font-medium text-sidebar-foreground text-sm">
                    Vessl
                  </span>
                  <span className="rounded bg-sidebar-accent/80 px-1.5 py-0.5 font-medium text-[10px] text-muted-foreground">
                    v0.1
                  </span>
                  <button
                    type="button"
                    onClick={onToggle}
                    className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-all duration-150 hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  >
                    <PanelLeft className="h-4 w-4" />
                  </button>
                </>
              )}
            </div>
          </div>

          {collapsed && (
            <button
              type="button"
              onClick={onToggle}
              className="absolute top-2 right-0 z-30 hidden h-7 w-7 translate-x-1/2 items-center justify-center rounded-lg border border-border/60 bg-card text-muted-foreground shadow-md transition-all duration-300 hover:bg-sidebar-accent hover:text-sidebar-foreground active:scale-[0.95] md:flex"
            >
              <PanelLeft className="h-4 w-4" />
            </button>
          )}
        </div>

        <nav className="flex flex-1 flex-col gap-5 overflow-y-auto px-2 pt-3 pb-3">
          {navGroups.map((group, i) => (
            <div key={i} className="flex flex-col gap-0.5">
              {!navCollapsed && group.title && (
                <h4 className="px-2 pb-1 font-medium text-[10px] text-sidebar-foreground/40 uppercase tracking-widest">
                  {group.title}
                </h4>
              )}
              {group.items.map((item) => (
                <NavItem key={item.url} item={item} exact={item.exact} collapsed={navCollapsed} />
              ))}
            </div>
          ))}
        </nav>

        <div
          className={`mt-auto flex flex-col gap-0.5 bg-sidebar-accent/20 ${navCollapsed ? 'px-1 py-1' : 'px-2 py-2'}`}
        >
          {bottomNav.map((item) => (
            <NavItem key={item.url} item={item} collapsed={navCollapsed} />
          ))}
        </div>

        <UserMenu collapsed={navCollapsed} />
      </aside>
    </>
  );
}
