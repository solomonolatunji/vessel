import { AlertCircle, Bot, Code2, Loader2, Terminal } from 'lucide-react';
import { useState } from 'react';
import { Alert, AlertDescription, AlertTitle } from '#/components/ui/alert';
import { Button } from '#/components/ui/button';
import { useExplainFailure } from '#/hooks/useDeployments';

interface DeploymentFailureAiProps {
  deploymentId: string;
}

export function DeploymentFailureAi({ deploymentId }: DeploymentFailureAiProps) {
  const [enabled, setEnabled] = useState(false);
  const { data, isLoading, error } = useExplainFailure(deploymentId, enabled);

  if (!enabled) {
    return (
      <div className="mt-4 rounded-lg border border-red-200 bg-red-50/50 p-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-red-100">
            <AlertCircle className="h-5 w-5 text-red-600" />
          </div>
          <div className="flex-1">
            <h3 className="font-medium text-red-900">Deployment Failed</h3>
            <p className="text-red-700 text-sm">
              This deployment encountered an error during the build or startup phase.
            </p>
          </div>
          <Button onClick={() => setEnabled(true)} variant="outline" className="gap-2 bg-white">
            <Bot className="h-4 w-4" />
            Analyze with AI
          </Button>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="mt-4 flex flex-col items-center justify-center space-y-3 rounded-lg border bg-gray-50/50 p-6 text-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <div>
          <h3 className="font-medium">AI is analyzing the failure...</h3>
          <p className="text-gray-500 text-sm">Reading build logs and checking configurations</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive" className="mt-4">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>AI Analysis Failed</AlertTitle>
        <AlertDescription>
          {error instanceof Error ? error.message : 'An unknown error occurred during analysis.'}
        </AlertDescription>
      </Alert>
    );
  }

  const explanation = data?.data;

  if (!explanation) return null;

  return (
    <div className="mt-4 overflow-hidden rounded-lg border bg-white shadow-sm">
      <div className="flex items-center gap-2 border-b bg-gray-50/80 px-4 py-3">
        <Bot className="h-5 w-5 text-primary" />
        <h3 className="font-medium">AI Failure Analysis</h3>
        <span
          className={`ml-auto rounded-full px-2 py-0.5 text-xs ${
            explanation.confidence === 'high'
              ? 'bg-green-100 text-green-700'
              : explanation.confidence === 'medium'
                ? 'bg-yellow-100 text-yellow-700'
                : 'bg-gray-100 text-gray-700'
          }`}
        >
          {explanation.confidence} confidence
        </span>
      </div>

      <div className="space-y-6 p-4">
        <div>
          <h4 className="mb-1 flex items-center gap-2 font-semibold text-gray-900 text-sm">
            Summary
          </h4>
          <p className="text-gray-700 text-sm">{explanation.summary}</p>
        </div>

        <div>
          <h4 className="mb-1 font-semibold text-gray-900 text-sm">Root Cause</h4>
          <p className="text-gray-700 text-sm">{explanation.cause}</p>
        </div>

        <div className="rounded-md border border-blue-100 bg-blue-50/50 p-4">
          <h4 className="mb-1 font-semibold text-blue-900 text-sm">Suggested Fix</h4>
          <p className="mb-3 text-blue-800 text-sm">{explanation.suggestedFix}</p>

          {explanation.commands && explanation.commands.length > 0 && (
            <div className="mt-4 space-y-2">
              <h5 className="font-medium text-blue-900/70 text-xs uppercase tracking-wider">
                Commands to run
              </h5>
              {explanation.commands.map((cmd, i) => (
                <div
                  key={i}
                  className="flex items-center gap-2 rounded bg-slate-900 p-2 text-slate-50"
                >
                  <Terminal className="h-4 w-4 text-slate-400" />
                  <code className="font-mono text-sm">{cmd}</code>
                </div>
              ))}
            </div>
          )}
        </div>

        {explanation.relatedLogLines && explanation.relatedLogLines.length > 0 && (
          <div>
            <h4 className="mb-2 flex items-center gap-2 font-semibold text-gray-900 text-sm">
              <Code2 className="h-4 w-4 text-gray-500" />
              Related Logs
            </h4>
            <div className="space-y-1 overflow-x-auto rounded border bg-slate-50 p-3">
              {explanation.relatedLogLines.map((line, i) => (
                <div
                  key={i}
                  className="border-slate-300 border-l-2 pl-2 font-mono text-slate-600 text-xs"
                >
                  {line}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
