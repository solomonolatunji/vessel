import { Database, LayoutTemplate, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '#/components/ui/sheet';
import { useListByProject } from '#/hooks/useApps';
import { useGetDatabases } from '#/hooks/useDatabases';

interface SmartLinkerDrawerProps {
  projectId: string;
  onLinkVariable: (key: string, value: string) => void;
  exampleEnvValues?: Record<string, string>;
}

export function SmartLinkerDrawer({
  projectId,
  onLinkVariable,
  exampleEnvValues,
}: SmartLinkerDrawerProps) {
  const [isOpen, setIsOpen] = useState(false);

  const { data: dbData, isLoading: isLoadingDbs } = useGetDatabases(projectId);
  const { data: appData, isLoading: isLoadingApps } = useListByProject(projectId);

  const databases = dbData?.data || [];
  const apps = appData?.data || [];

  return (
    <Sheet open={isOpen} onOpenChange={setIsOpen}>
      <SheetTrigger asChild>
        <Button variant="outline" size="sm">
          Smart Linker
        </Button>
      </SheetTrigger>
      <SheetContent className="w-100 overflow-y-auto sm:w-135">
        <SheetHeader>
          <SheetTitle>Smart Variable Linker</SheetTitle>
          <SheetDescription>
            Auto-link environment variables from databases or other services in this project.
          </SheetDescription>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {exampleEnvValues && Object.keys(exampleEnvValues).length > 0 && (
            <div className="space-y-3">
              <h3 className="font-medium text-sm">.env.example Suggestions</h3>
              <div className="flex flex-wrap gap-2">
                {Object.entries(exampleEnvValues).map(([key, value]) => (
                  <Badge
                    key={key}
                    variant="secondary"
                    className="cursor-pointer hover:bg-zinc-200 dark:hover:bg-zinc-800"
                    onClick={() => {
                      onLinkVariable(key, value || '');
                      setIsOpen(false);
                    }}
                  >
                    {key}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          <div className="space-y-4">
            <h3 className="font-medium text-sm">Available Resources</h3>
            {isLoadingDbs || isLoadingApps ? (
              <div className="flex items-center text-sm text-zinc-500">
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Loading resources...
              </div>
            ) : (
              <div className="grid gap-3">
                {databases.length === 0 && apps.length === 0 && (
                  <p className="text-sm text-zinc-500">No resources found in this project.</p>
                )}
                {databases.map((db) => (
                  <div
                    key={db.id}
                    className="flex items-center justify-between rounded-lg border p-3"
                  >
                    <div className="flex items-center gap-3">
                      <Database className="h-4 w-4 text-muted-foreground" />
                      <div>
                        <p className="font-medium text-sm">{db.name}</p>
                        <p className="text-muted-foreground text-xs">
                          {db.internalDns || 'Pending'}:{db.port}
                        </p>
                      </div>
                    </div>
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={() => {
                        const host = db.internalDns || 'localhost';
                        const port = db.port;
                        const url = `${db.engine}://${db.username}:${db.password}@${host}:${port}/${db.databaseName}`;
                        onLinkVariable('DATABASE_URL', url);
                        setIsOpen(false);
                      }}
                    >
                      Link URL
                    </Button>
                  </div>
                ))}
                {apps.map((app) => (
                  <div
                    key={app.id}
                    className="flex items-center justify-between rounded-lg border p-3"
                  >
                    <div className="flex items-center gap-3">
                      <LayoutTemplate className="h-4 w-4 text-muted-foreground" />
                      <div>
                        <p className="font-medium text-sm">{app.name}</p>
                        <p className="text-muted-foreground text-xs">Service</p>
                      </div>
                    </div>
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={() => {
                        onLinkVariable(
                          `${app.name.toUpperCase().replace(/-/g, '_')}_URL`,
                          `http://${app.name}`
                        );
                        setIsOpen(false);
                      }}
                    >
                      Link URL
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
