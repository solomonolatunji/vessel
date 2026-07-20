import { useCompletion } from '@ai-sdk/react';
import { AlertCircle, Bot, Loader2, Sparkles } from 'lucide-react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog';

interface AIDiagnoseDialogProps {
  logs: string;
}

export function AIDiagnoseDialog({ logs }: AIDiagnoseDialogProps) {
  // Use Vercel AI SDK to stream completion
  const { completion, isLoading, complete, error } = useCompletion({
    api: '/api/ai/diagnose', // This endpoint needs to be implemented in the Go backend later
  });

  const handleDiagnose = () => {
    // Send the last 2000 characters of logs to avoid token limits
    const logExcerpt = logs.slice(-2000);
    complete(logExcerpt);
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button size="sm" variant="outline" className="h-7 gap-2 px-2">
          <Sparkles className="h-3.5 w-3.5 text-blue-500" />
          <span className="text-xs">AI Diagnose</span>
        </Button>
      </DialogTrigger>
      <DialogContent className="border-zinc-800 bg-zinc-950 sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-zinc-100">
            <Bot className="h-5 w-5 text-blue-500" />
            AI Log Diagnosis
          </DialogTitle>
          <DialogDescription className="text-zinc-400">
            Analyzing your logs to find errors and suggest solutions.
          </DialogDescription>
        </DialogHeader>

        <div className="flex max-h-[400px] min-h-[200px] flex-col gap-4 overflow-y-auto py-4">
          {!completion && !isLoading && !error && (
            <div className="flex h-full flex-col items-center justify-center gap-3 text-zinc-500">
              <Sparkles className="h-8 w-8 opacity-50" />
              <p className="text-sm">Click the button below to start diagnosing the logs.</p>
            </div>
          )}

          {isLoading && !completion && (
            <div className="flex h-full items-center justify-center gap-3 text-zinc-500">
              <Loader2 className="h-6 w-6 animate-spin text-blue-500" />
              <p className="text-sm">Analyzing logs with AI...</p>
            </div>
          )}

          {error && (
            <div className="flex items-center gap-3 rounded-md border border-red-500/20 bg-red-400/10 p-3 text-red-400">
              <AlertCircle className="h-5 w-5 shrink-0" />
              <p className="text-sm">
                Failed to diagnose logs. Make sure the AI endpoint is configured.
              </p>
            </div>
          )}

          {completion && (
            <div className="prose prose-invert prose-sm max-w-none text-zinc-300">
              {completion.split('\n').map((line: string, i: number) => (
                <p key={i} className="mb-2">
                  {line}
                </p>
              ))}
            </div>
          )}
        </div>

        <div className="flex justify-end border-zinc-800 border-t pt-4">
          <Button onClick={handleDiagnose} disabled={isLoading || !logs.trim()} className="gap-2">
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Sparkles className="h-4 w-4" />
            )}
            {completion ? 'Re-analyze' : 'Start Analysis'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
