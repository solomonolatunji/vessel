import { Bell, Database, Lock, Settings as SettingsIcon, GitBranch } from 'lucide-react';
import { useState } from 'react';
import { BackupsList } from './backups-list';
import { GeneralSettings } from './general-settings';
import { GitAppsManager } from './git-apps-manager';
import { NotificationsSettings } from './notifications-settings';
import { OAuthProvidersList } from './oauth-providers-list';

type TabId = 'general' | 'notifications' | 'oauth' | 'backups' | 'git';

type Tab = { id: TabId; label: string; icon: React.ReactNode };

const TABS: Tab[] = [
  {
    id: 'general',
    label: 'General',
    icon: <SettingsIcon className="h-4 w-4" />,
  },
  {
    id: 'notifications',
    label: 'Notifications',
    icon: <Bell className="h-4 w-4" />,
  },
  { id: 'oauth', label: 'OAuth', icon: <Lock className="h-4 w-4" /> },
  { id: 'backups', label: 'Backups', icon: <Database className="h-4 w-4" /> },
  { id: 'git', label: 'Git Apps', icon: <GitBranch className="h-4 w-4" /> },
];

export const SettingsLayout = () => {
  const [activeId, setActiveId] = useState<TabId>('general');

  return (
    <div className="flex min-h-full flex-col">
      <div className="border-border/50 border-b bg-background/50 px-6 pt-6 backdrop-blur-sm">
        <div className="mb-5 flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <SettingsIcon className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Instance Settings</h1>
            <p className="text-muted-foreground text-sm">Manage your Vessl server configuration</p>
          </div>
        </div>

        <nav className="flex gap-1 overflow-x-auto">
          {TABS.map((t) => {
            const isActive = t.id === activeId;
            return (
              <button
                key={t.id}
                onClick={() => setActiveId(t.id)}
                type="button"
                className={[
                  'flex shrink-0 items-center gap-2 rounded-t-lg border border-b-0 px-4 py-2.5 text-sm transition-colors',
                  isActive
                    ? 'border-border/50 bg-card font-medium text-foreground'
                    : 'border-transparent text-muted-foreground hover:text-foreground',
                ].join(' ')}
              >
                {t.icon}
                {t.label}
              </button>
            );
          })}
        </nav>
      </div>

      <div className="flex-1 p-6">
        {activeId === 'general' && <GeneralSettings />}
        {activeId === 'notifications' && <NotificationsSettings />}
        {activeId === 'oauth' && <OAuthProvidersList />}
        {activeId === 'backups' && <BackupsList />}
        {activeId === 'git' && <GitAppsManager />}
      </div>
    </div>
  );
};
