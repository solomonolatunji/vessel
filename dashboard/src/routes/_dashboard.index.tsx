import { createFileRoute } from '@tanstack/react-router';
import { Box, LogOut, Plus, SearchIcon } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { CreateProjectModal } from '#/features/projects/create-project-modal';

export const Route = createFileRoute('/_dashboard/')({
  component: DashboardIndex,
});

function DashboardIndex() {
  const [createOpen, setCreateOpen] = useState(false);

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Box className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Projects</h1>
            <p className="text-muted-foreground text-sm">Manage your Vessl projects.</p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            className="gap-2"
            onClick={() => {
              const event = new KeyboardEvent('keydown', { key: 'k', metaKey: true });
              document.dispatchEvent(event);
            }}
          >
            <SearchIcon className="h-4 w-4" />
            SEARCH
          </Button>

          <Button onClick={() => setCreateOpen(true)} className="gap-2">
            <Plus className="h-4 w-4" />
            NEW PROJECT
          </Button>

          <Button variant="destructive" size="icon">
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <CreateProjectModal open={createOpen} onOpenChange={setCreateOpen} />
    </div>
  );
}
