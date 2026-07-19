import { useList } from '#/hooks/useBackups';
import { Loader2 } from 'lucide-react';

export function BackupManager({ database }: { database: any }) {
  const { data, isLoading } = useList(database.projectId);

  if (isLoading) {
    return <Loader2 className="animate-spin w-6 h-6 text-gray-500" />;
  }

  const backups = data?.data || [];

  return (
    <div className="rounded border p-4 shadow-sm space-y-4">
      <h2 className="font-semibold text-lg">Backups for {database.name}</h2>
      {backups.length === 0 ? (
        <p className="text-sm text-gray-500">No backups found.</p>
      ) : (
        <div className="space-y-2">
          {backups.map((b: any) => (
            <div key={b.id} className="p-2 border rounded text-sm flex justify-between">
              <span className="font-medium">{b.status}</span>
              <span>{new Date(b.createdAt).toLocaleString()}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
