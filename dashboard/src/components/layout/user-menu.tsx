import { Link } from '@tanstack/react-router';
import {
  BarChart3,
  FileText,
  HelpCircle,
  LogOut,
  Moon,
  MoreVertical,
  Settings,
  Sun,
  UserCircle,
} from 'lucide-react';
import { useTheme } from 'next-themes';
import { useEffect, useRef, useState } from 'react';
import { useLogout } from '#/hooks/useAuth';
import { authStore, useAuthState } from '#/stores/authStore';

export function UserMenu() {
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
    <div ref={ref} className="relative border-t border-sidebar-border p-3">
      {open && (
        <div className="absolute bottom-full left-3 right-3 mb-2 rounded-lg border border-border bg-popover shadow-xl overflow-hidden py-1 z-50 animate-in fade-in zoom-in-95 duration-100">
          <div className="px-3 py-2 border-b border-border mb-1">
            <p className="text-[13px] font-semibold text-foreground">{user?.name}</p>
            <p className="text-[11px] text-muted-foreground">{user?.email}</p>
          </div>

          <div className="px-1 py-1 space-y-0.5">
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <UserCircle className="h-4 w-4 text-muted-foreground" />
              Account Settings
            </Link>
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <Settings className="h-4 w-4 text-muted-foreground" />
              Workspace Settings
            </Link>
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
              Project Usage
            </Link>
          </div>

          <div className="my-1 h-px bg-border" />

          <div className="px-1 py-1 space-y-0.5">
            <Link
              to={'/docs' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <FileText className="h-4 w-4 text-muted-foreground" />
              Documentation
            </Link>
            <Link
              to={'/support' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <HelpCircle className="h-4 w-4 text-muted-foreground" />
              Support
            </Link>
          </div>

          <div className="my-1 h-px bg-border" />

          <div className="px-1 py-1 space-y-0.5">
            <button
              type="button"
              onClick={() => {
                setTheme(theme === 'dark' ? 'light' : 'dark');
                setOpen(false);
              }}
              className="flex w-full items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              {theme === 'dark' ? (
                <Sun className="h-4 w-4 text-muted-foreground" />
              ) : (
                <Moon className="h-4 w-4 text-muted-foreground" />
              )}
              {theme === 'dark' ? 'Light Theme' : 'Dark Theme'}
            </button>
            <button
              type="button"
              onClick={() => logout()}
              className="flex w-full items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] text-destructive hover:bg-destructive/10 transition-colors"
            >
              <LogOut className="h-4 w-4" />
              Log out
            </button>
          </div>
        </div>
      )}

      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex w-full items-center gap-2.5 rounded-lg hover:bg-sidebar-accent/50 p-1.5 transition-colors"
      >
        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary/20 text-[11px] font-bold text-primary">
          {initials}
        </div>
        <div className="min-w-0 text-left flex-1">
          <p className="truncate text-[12px] font-semibold text-sidebar-foreground leading-none">
            {user?.name ?? 'User'}
          </p>
          <p className="truncate text-[10px] text-sidebar-foreground/50 mt-0.5 leading-none">
            Workspace Owner
          </p>
        </div>
        <MoreVertical className="h-4 w-4 text-sidebar-foreground/50 shrink-0" />
      </button>
    </div>
  );
}
