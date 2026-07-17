import {
  Cloud,
  Code,
  Database,
  Download,
  FolderKanban,
  HardDrive,
  LayoutDashboard,
  LayoutTemplate,
  PanelLeft,
  ScrollText,
  Settings,
  Terminal,
  Users,
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
      { title: 'Projects', url: '/projects', icon: FolderKanban },
    ],
  },
  {
    title: 'Resources',
    items: [
      { title: 'Databases', url: '/databases', icon: Database },
      { title: 'Storage', url: '/storage', icon: HardDrive },
      { title: 'Sources', url: '/settings/git-apps', icon: Code },
    ],
  },
  {
    title: 'Discover',
    items: [
      { title: 'Templates', url: '/templates', icon: LayoutTemplate },
      { title: 'Import', url: '/imports/railway', icon: Download },
    ],
  },
  {
    title: 'System',
    items: [
      { title: 'Audit Logs', url: '/audit-logs', icon: ScrollText },
      { title: 'Terminal', url: '/terminal', icon: Terminal },
      { title: 'Users', url: '/settings/users', icon: Users },
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
}

export function AppSidebar({ collapsed, onToggle }: AppSidebarProps) {
  return (
    <aside
      className={`fixed inset-y-0 left-0 z-20 flex flex-col border-sidebar-border/50 border-r bg-sidebar/90 backdrop-blur-xl transition-all duration-300 ${
        collapsed ? 'w-16' : 'w-64'
      }`}
    >
      <div className="flex h-13 shrink-0 items-center gap-2 px-3">
        <div className="flex items-center gap-2.5">
          <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-primary/10">
            <Cloud className="h-4 w-4 text-primary" />
          </div>
          {!collapsed && (
            <span className="font-semibold text-sidebar-foreground text-sm tracking-tight">
              Vessl
            </span>
          )}
        </div>
        <div className="flex-1" />
        {!collapsed && (
          <span className="rounded bg-sidebar-accent/80 px-1.5 py-0.5 font-medium text-[10px] text-muted-foreground">
            v0.1
          </span>
        )}
        <button
          type="button"
          onClick={onToggle}
          className="flex h-6 w-6 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
        >
          <PanelLeft
            className={`h-4 w-4 transition-transform duration-300 ${collapsed ? 'scale-x-[-1]' : ''}`}
          />
        </button>
      </div>

      <nav className="flex flex-1 flex-col gap-5 overflow-y-auto px-2 pt-3 pb-3">
        {navGroups.map((group, i) => (
          <div key={i} className="flex flex-col gap-0.5">
            {!collapsed && group.title && (
              <h4 className="px-2 pb-1 font-medium text-[10px] text-sidebar-foreground/40 uppercase tracking-widest">
                {group.title}
              </h4>
            )}
            {group.items.map((item) => (
              <NavItem key={item.url} item={item} exact={item.exact} collapsed={collapsed} />
            ))}
          </div>
        ))}
      </nav>

      <div
        className={`mt-auto flex flex-col gap-0.5 bg-sidebar-accent/20 ${collapsed ? 'px-1 py-1' : 'px-2 py-2'}`}
      >
        {bottomNav.map((item) => (
          <NavItem key={item.url} item={item} collapsed={collapsed} />
        ))}
      </div>

      <UserMenu collapsed={collapsed} />
    </aside>
  );
}
