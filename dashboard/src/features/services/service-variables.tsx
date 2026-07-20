import { FileSearch, Loader2, Plus, Trash } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { useCreate, useDelete, useEnvSuggestions, useList } from '#/hooks/useServices';
import type { AppService, EnvExampleVariableSuggestion, Variable } from '#/interfaces/deployment';

export function ServiceVariables({ app }: { app: AppService }) {
  const { data, isLoading } = useList(app.id);
  const [scanEnabled, setScanEnabled] = useState(false);
  const { data: suggestionsData, isLoading: isScanning } = useEnvSuggestions(app.id, scanEnabled);
  const createVar = useCreate();
  const deleteVar = useDelete();
  const [key, setKey] = useState('');
  const [val, setVal] = useState('');

  if (isLoading) return <Loader2 className="animate-spin" />;

  const vars = data?.data || [];
  const suggestions = suggestionsData?.data || [];

  // Filter out suggestions that already exist in vars
  const availableSuggestions = suggestions.filter(
    (s: EnvExampleVariableSuggestion) => !vars.find((v: Variable) => v.key === s.key)
  );

  return (
    <div className="space-y-6">
      <div className="flex gap-2">
        <Input placeholder="Key" value={key} onChange={(e) => setKey(e.target.value)} />
        <Input placeholder="Value" value={val} onChange={(e) => setVal(e.target.value)} />
        <Button
          onClick={() => {
            if (!key) return;
            createVar.mutate({ serviceId: app.id, payload: { key, value: val, isSecret: false } });
            setKey('');
            setVal('');
          }}
        >
          Add
        </Button>
      </div>

      <div className="space-y-2">
        {vars.length === 0 ? (
          <div className="rounded-lg border border-dashed bg-gray-50/50 p-8 text-center">
            <p className="text-gray-500 text-sm">No environment variables configured</p>
          </div>
        ) : (
          vars.map((v: Variable) => (
            <div
              key={v.id}
              className="flex items-center justify-between rounded border bg-white p-3 shadow-sm"
            >
              <div className="font-mono text-sm">
                <span className="font-semibold text-primary">{v.key}</span>
                <span className="mx-2 text-gray-400">=</span>
                <span className="max-w-xs truncate text-gray-600">
                  {v.isSecret ? '********' : v.value}
                </span>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => deleteVar.mutate({ serviceId: app.id, id: v.id })}
                disabled={deleteVar.isPending}
              >
                <Trash className="h-4 w-4 text-red-500" />
              </Button>
            </div>
          ))
        )}
      </div>

      {app.repositoryUrl && (
        <div className="mt-8 rounded-lg border bg-slate-50/50 p-4">
          <div className="mb-4 flex items-center justify-between">
            <div>
              <h3 className="flex items-center gap-2 font-medium text-sm">
                <FileSearch className="h-4 w-4 text-primary" />
                Auto-fill from repository
              </h3>
              <p className="mt-1 text-gray-500 text-xs">
                Scan your repository for .env.example files to prepopulate variables
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setScanEnabled(true)}
              disabled={isScanning || availableSuggestions.length > 0}
            >
              {isScanning ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
              {availableSuggestions.length > 0 ? 'Scanned' : 'Scan Repository'}
            </Button>
          </div>

          {availableSuggestions.length > 0 && (
            <div className="mt-4 space-y-2 border-slate-200 border-t pt-4">
              {availableSuggestions.map((s: EnvExampleVariableSuggestion) => (
                <div
                  key={s.key}
                  className="flex items-center justify-between rounded border border-slate-200 bg-white p-2 text-sm"
                >
                  <div className="font-mono">
                    <span className="font-semibold text-slate-700">{s.key}</span>
                    {s.value && (
                      <>
                        <span className="mx-2 text-slate-400">=</span>
                        <span className="text-slate-500">{s.value}</span>
                      </>
                    )}
                  </div>
                  <Button
                    size="sm"
                    variant="secondary"
                    className="h-7 gap-1 text-xs"
                    onClick={() => {
                      createVar.mutate({
                        serviceId: app.id,
                        payload: { key: s.key, value: s.value, isSecret: false },
                      });
                    }}
                  >
                    <Plus className="h-3 w-3" /> Add
                  </Button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
