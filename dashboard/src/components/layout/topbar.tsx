import { Logout01Icon, Settings01Icon } from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { useNavigate } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';

import { Avatar, AvatarFallback, AvatarImage } from '#/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { SidebarTrigger } from '#/components/ui/sidebar';
import { useLogout } from '#/hooks/useAuth';
import { authStore } from '#/stores/authStore';

export function Topbar() {
  const authState = useStore(authStore);
  const { mutateAsync: logout } = useLogout();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate({ to: '/login' });
  };

  const user = authState.user;
  const initials = user?.name
    ? user.name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
    : 'U';

  return (
    <header className="h-16 shrink-0 border-b flex items-center justify-between px-4 sticky top-0 bg-background/80 backdrop-blur-md z-10">
      <div className="flex items-center gap-4">
        <SidebarTrigger className="text-muted-foreground hover:text-foreground" />

        {/* System Health Indicators */}
        <div className="hidden md:flex items-center gap-4 text-sm font-medium">
          <div className="flex items-center gap-2 text-muted-foreground">
            <span className="flex h-2 w-2 rounded-full bg-emerald-500 ring-2 ring-emerald-500/20"></span>
            Docker: Online
          </div>
          <div className="h-4 w-px bg-border"></div>
          <div className="flex items-center gap-2 text-muted-foreground">
            <span>
              CPU: <span className="text-foreground">12%</span>
            </span>
          </div>
          <div className="h-4 w-px bg-border"></div>
          <div className="flex items-center gap-2 text-muted-foreground">
            <span>
              RAM: <span className="text-foreground">4.2GB</span>
            </span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4">
        <div className="hidden sm:block bg-primary/10 text-primary border border-primary/20 px-2 py-1 rounded-md text-xs font-medium">
          v0.1.0 Available
        </div>

        {/* Auth User Menu */}
        {user && (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button
                type="button"
                className="flex items-center gap-2 outline-none focus-visible:ring-2 focus-visible:ring-ring rounded-full"
              >
                <Avatar className="h-8 w-8 border border-border">
                  <AvatarImage
                    src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${user.email}`}
                    alt={user.name}
                  />
                  <AvatarFallback className="bg-primary/10 text-primary text-xs font-medium">
                    {initials}
                  </AvatarFallback>
                </Avatar>
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel className="font-normal">
                <div className="flex flex-col space-y-1">
                  <p className="text-sm font-medium leading-none">{user.name}</p>
                  <p className="text-xs leading-none text-muted-foreground">{user.email}</p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="text-muted-foreground hover:text-foreground cursor-pointer">
                <HugeiconsIcon icon={Settings01Icon} className="mr-2 h-4 w-4" />
                <span>Account Settings</span>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={handleLogout}
                className="text-destructive focus:text-destructive cursor-pointer"
              >
                <HugeiconsIcon icon={Logout01Icon} className="mr-2 h-4 w-4" />
                <span>Log out</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </div>
    </header>
  );
}
