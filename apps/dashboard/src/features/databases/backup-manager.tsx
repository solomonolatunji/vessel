import { Loader2 } from 'lucide-react';
import { useList } from '#/hooks/useBackups';

export function BackupManager({ database }: { database: any }) {
  const { data, isLoading } = useList();

  if (isLoading) {
    return <Loader2 className="h-6 w-6 animate-spin text-gray-500" />;
  }

  const backups = data?.data || [];

  return (
    <div className="space-y-4 rounded border p-4 shadow-sm">
      <h2 className="font-semibold text-lg">Backups for {database.name}</h2>
      {backups.length === 0 ? (
        <p className="text-gray-500 text-sm">No backups found.</p>
      ) : (
        <div className="space-y-2">
          {backups.map((b: { id: string; status: string; createdAt: string }) => (
            <div key={b.id} className="flex justify-between rounded border p-2 text-sm">
              <span className="font-medium">{b.status}</span>
              <span>{new Date(b.createdAt).toLocaleString()}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
