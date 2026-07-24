import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { BackupManager } from '#/features/databases/backup-manager';
import { DataBrowser } from '#/features/databases/data-browser';
import { DatabaseConnectionCard } from '#/features/databases/database-connection-card';
import { DatabaseImportModal } from '#/features/databases/database-import-modal';
import { DatabaseNetworking } from '#/features/databases/database-networking';
import { RedisKeyBrowser } from '#/features/databases/redis-key-browser';

import { HealthcheckCard } from '#/features/services/healthcheck-card';
import { MaintenanceModeCard } from '#/features/services/maintenance-mode-card';
import { RuntimeModeCard } from '#/features/services/runtime-mode-card';
import { useGetApp } from '#/hooks/useApps';
import { useGetDatabase } from '#/hooks/useDatabases';

export const Route = createFileRoute('/_dashboard/services/$serviceId/')({
  component: ServiceIndexRoute,
});

function ServiceIndexRoute() {
  const { serviceId } = Route.useParams();
  const [importModalOpen, setImportModalOpen] = useState(false);

  const { data: appData, isLoading: appLoading } = useGetApp(serviceId);
  const { data: dbData, isLoading: dbLoading } = useGetDatabase(serviceId);

  if (appLoading || dbLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const app = appData?.data;
  const db = dbData?.data;

  if (db) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="font-bold text-2xl">Database: {db.name}</h1>
          <Button onClick={() => setImportModalOpen(true)}>Import Data</Button>
        </div>
        <DatabaseImportModal
          isOpen={importModalOpen}
          onOpenChange={setImportModalOpen}
          databaseId={db.id}
        />
        <DatabaseConnectionCard database={db} />
        <DatabaseNetworking
          database={{ id: db.id, isPublic: !!db.externalDns, publicEndpoint: db.externalDns }}
        />
        {['postgres', 'postgresql', 'mysql', 'mariadb'].includes(db.engine) ? (
          <DataBrowser databaseId={db.id} />
        ) : db.engine === 'redis' ? (
          <RedisKeyBrowser databaseId={db.id} />
        ) : (
          <div className="rounded-lg border p-8 text-center text-muted-foreground">
            Data browser coming soon for {db.engine}
          </div>
        )}
        <BackupManager database={db} />
      </div>
    );
  }

  if (app) {
    return (
      <div className="space-y-6">
        <h1 className="font-bold text-2xl">Service: {app.name}</h1>
        <RuntimeModeCard serviceId={app.id} />
        <HealthcheckCard serviceId={app.id} />
        <MaintenanceModeCard serviceId={app.id} />
      </div>
    );
  }

  return <div>Service not found.</div>;
}
