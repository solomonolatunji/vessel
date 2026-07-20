import Editor from '@monaco-editor/react';
import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useGetApp } from '#/hooks/useApps';
import { useGetCode, useSaveCode } from '#/hooks/useServerless';

const TEMPLATES: Record<string, string> = {
  nodejs: `module.exports = async (req, res) => {
  res.send("Hello from Node.js Edge Function!");
};`,
  python: `def handler(request):
    return "Hello from Python Flask Function!", 200`,
  go: `package main

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Go Serverless Function!"))
}
`,
};

export const Route = createFileRoute('/_dashboard/services/$serviceId/serverless')({
  component: ServiceServerlessRoute,
});

function ServiceServerlessRoute() {
  const { serviceId } = Route.useParams();
  const { data: appData, isLoading: appLoading } = useGetApp(serviceId);
  const { data: codeData, isLoading: codeLoading } = useGetCode(serviceId);
  const saveCode = useSaveCode();

  const [runtime, setRuntime] = useState('nodejs');
  const [code, setCode] = useState(TEMPLATES.nodejs);

  useEffect(() => {
    if (codeData?.data) {
      if (codeData.data.codeContent) setCode(codeData.data.codeContent);
      if (codeData.data.runtime) setRuntime(codeData.data.runtime);
    }
  }, [codeData]);

  const handleRuntimeChange = (newRuntime: string) => {
    setRuntime(newRuntime);
    if (!code || Object.values(TEMPLATES).includes(code)) {
      setCode(TEMPLATES[newRuntime] || '');
    }
  };

  if (appLoading || codeLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const app = appData?.data;

  if (!app) {
    return <div>Service not found.</div>;
  }

  const handleSave = async () => {
    try {
      await saveCode.mutateAsync({
        serviceId,
        payload: { codeContent: code, runtime },
      });
      toast.success('Function code saved successfully. You can now deploy this service.');
    } catch (err: any) {
      toast.error(err.message || 'Failed to save code');
    }
  };

  const getLanguage = () => {
    switch (runtime) {
      case 'nodejs':
        return 'javascript';
      case 'python':
        return 'python';
      case 'go':
        return 'go';
      default:
        return 'javascript';
    }
  };

  return (
    <div className="flex h-full flex-col space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-bold text-2xl">Serverless Functions</h1>
        <div className="flex items-center gap-4">
          <Select value={runtime} onValueChange={handleRuntimeChange}>
            <SelectTrigger className="w-32">
              <SelectValue placeholder="Runtime" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="nodejs">Node.js</SelectItem>
              <SelectItem value="python">Python</SelectItem>
              <SelectItem value="go">Go</SelectItem>
            </SelectContent>
          </Select>
          <Button onClick={handleSave} disabled={saveCode.isPending}>
            {saveCode.isPending ? 'Saving...' : 'Save Code'}
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-hidden rounded-md border border-gray-200">
        <Editor
          height="100%"
          language={getLanguage()}
          theme="vs-dark"
          value={code}
          onChange={(value) => setCode(value || '')}
          options={{
            minimap: { enabled: false },
            fontSize: 14,
          }}
        />
      </div>
    </div>
  );
}
