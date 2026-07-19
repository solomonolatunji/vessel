import { Loader2, Trash } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { useCreate, useDelete, useList } from '#/hooks/useServices';

export function ServiceVariables({ app }: { app: any }) {
  const { data, isLoading } = useList(app.id);
  const createVar = useCreate();
  const deleteVar = useDelete();
  const [key, setKey] = useState('');
  const [val, setVal] = useState('');

  if (isLoading) return <Loader2 className="animate-spin" />;

  const vars = data?.data || [];

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <Input placeholder="Key" value={key} onChange={(e) => setKey(e.target.value)} />
        <Input placeholder="Value" value={val} onChange={(e) => setVal(e.target.value)} />
        <Button
          onClick={() =>
            createVar.mutate({ serviceId: app.id, payload: { key, value: val, isSecret: false } })
          }
        >
          Add
        </Button>
      </div>
      <div className="space-y-2">
        {vars.map((v) => (
          <div key={v.id} className="flex justify-between rounded border p-2">
            <span>
              {v.key} = {v.value}
            </span>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => deleteVar.mutate({ serviceId: app.id, id: v.id })}
            >
              <Trash className="h-4 w-4 text-red-500" />
            </Button>
          </div>
        ))}
      </div>
    </div>
  );
}
