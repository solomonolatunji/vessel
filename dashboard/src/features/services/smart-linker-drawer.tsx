import { useQuery } from '@tanstack/react-query';
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
import { apiClient } from '#/lib/apiClient';

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

  // Fetch available databases or services to link variables from
  const { data: linkableResources, isLoading } = useQuery({
    queryKey: ['project-resources', projectId],
    queryFn: () => apiClient.get<any[]>(`/projects/${projectId}/resources`),
  });

  return (
    <Sheet open={isOpen} onOpenChange={setIsOpen}>
      <SheetTrigger asChild>
        <Button variant="outline" size="sm">
          Smart Linker
        </Button>
      </SheetTrigger>
      <SheetContent className="w-[400px] overflow-y-auto sm:w-[540px]">
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
            {isLoading ? (
              <p className="text-sm text-zinc-500">Loading resources...</p>
            ) : linkableResources?.length ? (
              linkableResources.map((res) => (
                <div key={res.id} className="rounded-md border p-3">
                  <div className="mb-2 font-medium">
                    {res.name} ({res.type})
                  </div>
                  <div className="space-y-2">
                    {res.variables?.map((v: { key: string; value: string }) => (
                      <div
                        key={v.key}
                        className="flex items-center justify-between rounded bg-zinc-50 p-2 text-sm dark:bg-zinc-900"
                      >
                        <span className="font-mono text-xs">{v.key}</span>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 text-xs"
                          onClick={() => {
                            onLinkVariable(v.key, `\${${res.name}.${v.key}}`);
                            setIsOpen(false);
                          }}
                        >
                          Link
                        </Button>
                      </div>
                    ))}
                    {(!res.variables || res.variables.length === 0) && (
                      <p className="text-xs text-zinc-500">No exportable variables found.</p>
                    )}
                  </div>
                </div>
              ))
            ) : (
              <p className="text-sm text-zinc-500">No linkable resources found in this project.</p>
            )}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
