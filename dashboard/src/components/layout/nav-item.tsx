import { Link, useRouterState } from '@tanstack/react-router';
import { ExternalLink } from 'lucide-react';
import type React from 'react';

export type NavItemProps = {
  title: string;
  url: string;
  icon: React.ComponentType<{ className?: string }>;
  external?: boolean;
  badge?: string;
};

export function NavItem({
  item,
  exact = false,
  collapsed = false,
}: {
  item: NavItemProps;
  exact?: boolean;
  collapsed?: boolean;
}) {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const isActive = exact
    ? pathname === item.url
    : pathname.startsWith(item.url) && item.url !== '/';

  return (
    <Link
      to={item.url as never}
      className={`group relative flex items-center rounded-xl font-medium text-sm transition-all duration-150 active:scale-[0.985] ${
        collapsed ? 'justify-center gap-0 px-0 py-2' : 'gap-3 px-2.5 py-2'
      } ${
        isActive
          ? 'nav-active-glow bg-sidebar-accent text-sidebar-accent-foreground shadow-inner'
          : 'text-sidebar-foreground/60 hover:bg-sidebar-accent/60 hover:text-sidebar-foreground'
      }`}
      target={item.external ? '_blank' : undefined}
      rel={item.external ? 'noopener noreferrer' : undefined}
    >
      {!collapsed && isActive && (
        <div className="absolute top-1/2 -left-2 h-5 w-0.75 -translate-y-1/2 rounded-r-full bg-primary" />
      )}

      <div
        className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-lg border transition-all duration-150 ${
          isActive
            ? 'border-primary/25 bg-primary/10 text-primary'
            : 'border-transparent text-muted-foreground group-hover:border-sidebar-border/60 group-hover:bg-sidebar-accent group-hover:text-sidebar-foreground'
        }`}
      >
        <item.icon className="h-4 w-4" />
      </div>

      {!collapsed && (
        <>
          <span className="flex-1 truncate">{item.title}</span>
          {item.external && <ExternalLink className="h-3 w-3 shrink-0 opacity-50" />}
          {item.badge && (
            <span className="rounded-full bg-primary/10 px-1.5 py-0.5 font-medium text-[9px] text-primary">
              {item.badge}
            </span>
          )}
        </>
      )}
    </Link>
  );
}
