import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { apiClient } from '#/lib/apiClient';

interface BuildSettingsProps {
  serviceId: string;
  initialData?: {
    staticOutputDir?: string;
    installCmd?: string;
    buildCmd?: string;
    startCmd?: string;
  };
}

export function BuildSettings({ serviceId, initialData }: BuildSettingsProps) {
  const queryClient = useQueryClient();
  const [staticOutputDir, setStaticOutputDir] = useState(initialData?.staticOutputDir || '');
  const [installCmd, setInstallCmd] = useState(initialData?.installCmd || '');
  const [buildCmd, setBuildCmd] = useState(initialData?.buildCmd || '');
  const [startCmd, setStartCmd] = useState(initialData?.startCmd || '');
  const [isExpanded, setIsExpanded] = useState(false);

  const updateMutation = useMutation({
    mutationFn: async (data: Record<string, string>) => {
      return apiClient.put(`/services/${serviceId}/build`, data);
    },
    onSuccess: () => {
      toast.success('Build settings updated');
      queryClient.invalidateQueries({ queryKey: ['service', serviceId] });
    },
  });

  const handleSave = () => {
    updateMutation.mutate({
      staticOutputDir,
      installCmd,
      buildCmd,
      startCmd,
    });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Build & Output Settings</CardTitle>
        <CardDescription>Configure how your service is built and served.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="staticOutputDir">Static Output Directory</Label>
          <div className="flex items-start gap-2">
            <div className="flex-1 space-y-2">
              <Input
                id="staticOutputDir"
                placeholder="e.g., dist, build, .output/public"
                value={staticOutputDir}
                onChange={(e) => setStaticOutputDir(e.target.value)}
              />
              <p className="text-sm text-zinc-500">Leave empty for standard dynamic deployments.</p>
            </div>
          </div>
          {staticOutputDir && (
            <div className="mt-2">
              <Badge variant="secondary" className="font-mono text-xs">
                nginx:alpine (port 80) wrapper active
              </Badge>
            </div>
          )}
        </div>

        <div className="rounded-md border border-zinc-200 dark:border-zinc-800">
          <button
            type="button"
            className="flex w-full items-center justify-between p-4 font-medium text-sm hover:bg-zinc-50 dark:hover:bg-zinc-900/50"
            onClick={() => setIsExpanded(!isExpanded)}
          >
            <span>Railpack / Nixpacks Build Overrides</span>
            <span className="text-zinc-500">{isExpanded ? '▼' : '▶'}</span>
          </button>

          {isExpanded && (
            <div className="space-y-4 border-zinc-200 border-t p-4 dark:border-zinc-800">
              <div className="space-y-2">
                <Label htmlFor="installCmd">Install Command (--install-cmd)</Label>
                <Input
                  id="installCmd"
                  placeholder="e.g., npm ci"
                  value={installCmd}
                  onChange={(e) => setInstallCmd(e.target.value)}
                  className="font-mono text-sm"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="buildCmd">Build Command (--build-cmd)</Label>
                <Input
                  id="buildCmd"
                  placeholder="e.g., npm run build"
                  value={buildCmd}
                  onChange={(e) => setBuildCmd(e.target.value)}
                  className="font-mono text-sm"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="startCmd">Start Command (--start-cmd)</Label>
                <Input
                  id="startCmd"
                  placeholder="e.g., npm start"
                  value={startCmd}
                  onChange={(e) => setStartCmd(e.target.value)}
                  className="font-mono text-sm"
                />
              </div>
            </div>
          )}
        </div>

        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={updateMutation.isPending}>
            {updateMutation.isPending ? 'Saving...' : 'Save Settings'}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
