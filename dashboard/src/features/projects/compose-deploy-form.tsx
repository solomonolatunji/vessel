import Editor from '@monaco-editor/react';
import { Database, Play, Search, Server } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { useAnalyzeCompose, useDeployCompose } from '#/hooks/useCompose';
import type { ComposeAnalyzeResponse } from '#/services/compose';

export function ComposeDeployForm({ projectId }: { projectId: string }) {
  const [composeText, setComposeText] = useState(`version: '3'
services:
  web:
    image: nginx
    ports:
      - "80:80"`);

  const [analysisResult, setAnalysisResult] = useState<ComposeAnalyzeResponse | null>(null);

  const analyzeMutation = useAnalyzeCompose();
  const deployMutation = useDeployCompose();

  const handleAnalyze = () => {
    analyzeMutation.mutate(
      { projectId, composeContent: composeText },
      {
        onSuccess: (data) => {
          setAnalysisResult(data);
          toast.success('Compose file analyzed successfully!');
        },
        onError: (err: any) => {
          toast.error(err.response?.data?.message || 'Failed to analyze compose file');
        },
      }
    );
  };

  const handleDeploy = () => {
    deployMutation.mutate(
      { projectId, composeContent: composeText },
      {
        onSuccess: () => {
          toast.success('Compose stack deployed successfully!');
          setAnalysisResult(null);
        },
        onError: (err: any) => {
          toast.error(err.response?.data?.message || 'Failed to deploy compose stack');
        },
      }
    );
  };

  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <Card className="flex h-[70vh] flex-col">
        <CardHeader>
          <CardTitle>Docker Compose</CardTitle>
          <CardDescription>Paste your docker-compose.yml here</CardDescription>
        </CardHeader>
        <CardContent className="flex min-h-0 flex-1 flex-col p-0 pb-4">
          <div className="w-full flex-1 border-y">
            <Editor
              height="100%"
              defaultLanguage="yaml"
              theme="vs-dark"
              value={composeText}
              onChange={(value) => setComposeText(value || '')}
              options={{
                minimap: { enabled: false },
                fontSize: 14,
                lineNumbers: 'on',
                scrollBeyondLastLine: false,
              }}
            />
          </div>
          <div className="flex justify-end gap-2 p-4 pb-0">
            <Button
              variant="outline"
              onClick={handleAnalyze}
              disabled={analyzeMutation.isPending || !composeText.trim()}
            >
              <Search className="mr-2 h-4 w-4" />
              {analyzeMutation.isPending ? 'Analyzing...' : 'Analyze'}
            </Button>
            <Button
              onClick={handleDeploy}
              disabled={deployMutation.isPending || !composeText.trim() || !analysisResult}
            >
              <Play className="mr-2 h-4 w-4" />
              {deployMutation.isPending ? 'Deploying...' : 'Deploy'}
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card className="h-[70vh] overflow-y-auto bg-muted/30">
        <CardHeader>
          <CardTitle>Preview</CardTitle>
          <CardDescription>
            {analysisResult
              ? 'These resources will be created in your project'
              : 'Click Analyze to preview the resources that will be created'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {!analysisResult && !analyzeMutation.isPending && (
            <div className="flex h-32 items-center justify-center rounded-lg border border-dashed text-muted-foreground">
              No preview available
            </div>
          )}

          {analyzeMutation.isPending && (
            <div className="flex h-32 animate-pulse items-center justify-center rounded-lg border border-dashed text-muted-foreground">
              Analyzing resources...
            </div>
          )}

          {analysisResult && (
            <div className="flex flex-col gap-6">
              {analysisResult.appServices?.length > 0 && (
                <div className="flex flex-col gap-3">
                  <h3 className="flex items-center gap-2 font-semibold">
                    <Server className="h-4 w-4" />
                    App Services ({analysisResult.appServices.length})
                  </h3>
                  {analysisResult.appServices.map((svc: any, idx: number) => (
                    <div key={idx} className="rounded-lg border bg-background p-4 shadow-sm">
                      <div className="font-medium text-lg">{svc.name}</div>
                      <div className="mt-1 grid grid-cols-2 gap-2 text-muted-foreground text-sm">
                        <div>
                          <span className="font-medium text-foreground">Image:</span>{' '}
                          {svc.imageRef || 'Will be built'}
                        </div>
                        <div>
                          <span className="font-medium text-foreground">Runtime:</span>{' '}
                          {svc.runtimeMode}
                        </div>
                        {svc.buildEngine && (
                          <div>
                            <span className="font-medium text-foreground">Build Engine:</span>{' '}
                            {svc.buildEngine}
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {analysisResult.databases?.length > 0 && (
                <div className="flex flex-col gap-3">
                  <h3 className="flex items-center gap-2 font-semibold">
                    <Database className="h-4 w-4" />
                    Databases ({analysisResult.databases.length})
                  </h3>
                  {analysisResult.databases.map((db: any, idx: number) => (
                    <div key={idx} className="rounded-lg border bg-background p-4 shadow-sm">
                      <div className="font-medium text-lg capitalize">{db.name}</div>
                      <div className="mt-1 grid grid-cols-2 gap-2 text-muted-foreground text-sm">
                        <div>
                          <span className="font-medium text-foreground">Engine:</span> {db.engine}
                        </div>
                        <div>
                          <span className="font-medium text-foreground">Version:</span> {db.version}
                        </div>
                        <div>
                          <span className="font-medium text-foreground">Port:</span> {db.port}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {analysisResult.appServices?.length === 0 &&
                analysisResult.databases?.length === 0 && (
                  <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-4 text-sm text-yellow-600 dark:border-yellow-900 dark:bg-yellow-950/20 dark:text-yellow-400">
                    No recognizable services found. Codedock looks for valid image names or build
                    directives to generate AppServices and Databases.
                  </div>
                )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
