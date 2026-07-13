import { Link, useRouterState } from '@tanstack/react-router';
import { ExternalLink } from 'lucide-react';
import type React from 'react';

export type NavItemProps = {
  title: string;
  url: string;
  icon: React.ComponentType<{ className?: string }>;
  external?: boolean;
};

export function NavItem({ item, exact = false }: { item: NavItemProps; exact?: boolean }) {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const isActive = exact
    ? pathname === item.url
    : pathname.startsWith(item.url) && item.url !== '/';

  const className = [
    'group flex items-center gap-3 rounded-md px-3 py-2 text-[13px] font-medium transition-all duration-100',
    isActive
      ? 'bg-sidebar-accent text-sidebar-accent-foreground'
      : 'text-sidebar-foreground/60 hover:text-sidebar-foreground hover:bg-sidebar-accent/60',
  ].join(' ');

  const IconComponent = (
    <item.icon
      className={[
        'h-4 w-4 shrink-0',
        isActive ? 'text-primary' : 'group-hover:text-sidebar-foreground/70',
      ].join(' ')}
    />
  );

  if (item.external) {
    return (
      <a href={item.url} target="_blank" rel="noopener noreferrer" className={className}>
        {IconComponent}
        <span className="truncate flex-1">{item.title}</span>
        <ExternalLink className="h-3.5 w-3.5 opacity-50 shrink-0" />
      </a>
    );
  }

  return (
    <Link to={item.url as never} className={className}>
      {IconComponent}
      <span className="truncate">{item.title}</span>
    </Link>
  );
}
