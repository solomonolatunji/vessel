import { Link } from '@tanstack/react-router';
import {
  Heart,
  HelpCircle,
  LogOut,
  MessageSquare,
  Moon,
  MoreVertical,
  Sun,
  UserCircle,
} from 'lucide-react';
import { useTheme } from 'next-themes';
import { useEffect, useRef, useState } from 'react';
import { useLogout } from '#/hooks/useAuth';
import { useAuthState } from '#/stores/authStore';

interface UserMenuProps {
  collapsed: boolean;
}

export function UserMenu({ collapsed }: UserMenuProps) {
  const { theme, setTheme } = useTheme();
  const authState = useAuthState();
  const { mutateAsync: logout } = useLogout();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const user = authState.user;
  const initials = user?.name
    ? user.name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
    : 'U';

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  return (
    <div
      ref={ref}
      className={`relative bg-sidebar-accent/20 ${collapsed ? 'px-1 py-1' : 'px-2 py-2'}`}
    >
      {open && (
        <div
          className={`fade-in slide-in-from-bottom-2 absolute z-50 mb-3 animate-in rounded-2xl border border-border/80 bg-card p-1.5 shadow-2xl shadow-black/10 backdrop-blur-xl duration-200 dark:shadow-black/40 ${
            collapsed ? 'bottom-0 left-full ml-3 w-56' : 'right-2 bottom-full left-2'
          }`}
        >
          <div className="border-sidebar-border/50 border-b px-3 pt-2 pb-2.5">
            <p className="font-semibold text-sidebar-foreground text-sm">{user?.name}</p>
            <p className="text-[12px] text-sidebar-foreground/60">{user?.email}</p>
          </div>

          <div className="space-y-0.5 py-1.5">
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="group flex items-center gap-3 rounded-xl px-2.5 py-2 font-medium text-sidebar-foreground/60 text-sm transition-all duration-150 hover:bg-sidebar-accent hover:text-sidebar-foreground active:scale-[0.985]"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg border border-transparent text-muted-foreground transition-all duration-150 group-hover:border-sidebar-border/60 group-hover:bg-sidebar-accent group-hover:text-sidebar-foreground">
                <UserCircle className="h-3.5 w-3.5" />
              </div>
              Account Settings
            </Link>
          </div>

          <div className="mx-3 my-1 h-px bg-sidebar-border/40" />

          <div className="space-y-0.5 py-1.5">
            <button
              type="button"
              onClick={() => {
                window.open('https://github.com/anomalyco/vessl', '_blank');
                setOpen(false);
              }}
              className="group flex w-full items-center gap-3 rounded-xl px-2.5 py-2 font-medium text-sidebar-foreground/60 text-sm transition-all duration-150 hover:bg-sidebar-accent hover:text-sidebar-foreground active:scale-[0.985]"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg border border-transparent text-muted-foreground transition-all duration-150 group-hover:border-sidebar-border/60 group-hover:bg-sidebar-accent group-hover:text-sidebar-foreground">
                <HelpCircle className="h-3.5 w-3.5" />
              </div>
              Get Help
            </button>
            <button
              type="button"
              onClick={() => {
                window.open('https://github.com/sponsors/anomalyco', '_blank');
                setOpen(false);
              }}
              className="group flex w-full items-center gap-3 rounded-xl px-2.5 py-2 font-medium text-sidebar-foreground/60 text-sm transition-all duration-150 hover:bg-sidebar-accent hover:text-sidebar-foreground active:scale-[0.985]"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg border border-transparent text-muted-foreground transition-all duration-150 group-hover:border-sidebar-border/60 group-hover:bg-sidebar-accent group-hover:text-sidebar-foreground">
                <Heart className="h-3.5 w-3.5" />
              </div>
              Sponsor
            </button>
            <button
              type="button"
              onClick={() => {
                window.open('https://github.com/anomalyco/vessl/issues', '_blank');
                setOpen(false);
              }}
              className="group flex w-full items-center gap-3 rounded-xl px-2.5 py-2 font-medium text-sidebar-foreground/60 text-sm transition-all duration-150 hover:bg-sidebar-accent hover:text-sidebar-foreground active:scale-[0.985]"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg border border-transparent text-muted-foreground transition-all duration-150 group-hover:border-sidebar-border/60 group-hover:bg-sidebar-accent group-hover:text-sidebar-foreground">
                <MessageSquare className="h-3.5 w-3.5" />
              </div>
              Share Feedback
            </button>
          </div>

          <div className="mx-3 my-1 h-px bg-sidebar-border/40" />

          <div className="space-y-0.5 py-1.5">
            <button
              type="button"
              onClick={() => {
                setTheme(theme === 'dark' ? 'light' : 'dark');
                setOpen(false);
              }}
              className="group flex w-full items-center gap-3 rounded-xl px-2.5 py-2 font-medium text-sidebar-foreground/60 text-sm transition-all duration-150 hover:bg-sidebar-accent hover:text-sidebar-foreground active:scale-[0.985]"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg border border-transparent text-muted-foreground transition-all duration-150 group-hover:border-sidebar-border/60 group-hover:bg-sidebar-accent group-hover:text-sidebar-foreground">
                {theme === 'dark' ? (
                  <Sun className="h-3.5 w-3.5" />
                ) : (
                  <Moon className="h-3.5 w-3.5" />
                )}
              </div>
              {theme === 'dark' ? 'Light Theme' : 'Dark Theme'}
            </button>
          </div>

          <div className="mx-3 my-1 h-px bg-sidebar-border/40" />

          <div className="space-y-0.5 py-1.5">
            <button
              type="button"
              onClick={() => logout()}
              className="group flex w-full items-center gap-3 rounded-xl px-2.5 py-2 font-medium text-sm transition-all duration-150 hover:bg-destructive/10 active:scale-[0.985]"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg border border-transparent text-destructive/70 transition-all duration-150 group-hover:border-destructive/20 group-hover:bg-destructive/10 group-hover:text-destructive">
                <LogOut className="h-3.5 w-3.5" />
              </div>
              <span className="text-destructive/80 group-hover:text-destructive">Log out</span>
            </button>
          </div>
        </div>
      )}

      <button
        type="button"
        onClick={() => setOpen(!open)}
        className={`flex w-full items-center rounded-xl transition-colors hover:bg-sidebar-accent/50 ${
          collapsed ? 'justify-center gap-0 px-0 py-2' : 'gap-3 px-2.5 py-2'
        }`}
      >
        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary/15 font-bold text-[10px] text-primary">
          {initials}
        </div>
        {!collapsed && (
          <>
            <div className="min-w-0 flex-1 text-left">
              <p className="truncate font-medium text-[12px] text-sidebar-foreground leading-tight">
                {user?.name ?? 'User'}
              </p>
            </div>
            <MoreVertical className="h-3.5 w-3.5 shrink-0 text-sidebar-foreground/40" />
          </>
        )}
      </button>
    </div>
  );
}
