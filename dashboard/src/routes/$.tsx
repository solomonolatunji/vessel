import { Alert01Icon, ArrowLeft01Icon } from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { createFileRoute, Link } from '@tanstack/react-router';

export const Route = createFileRoute('/$')({
  component: NotFound,
});

function NotFound() {
  return (
    <div className="flex h-full w-full flex-col items-center justify-center p-8 text-center text-zinc-100">
      <div className="flex items-center justify-center h-20 w-20 rounded-full bg-red-500/10 mb-6">
        <HugeiconsIcon icon={Alert01Icon} className="h-10 w-10 text-red-500" />
      </div>
      <h1 className="text-4xl font-bold tracking-tight mb-2">404 - Page Not Found</h1>
      <p className="text-zinc-400 mb-8 max-w-md">
        The page you are looking for doesn't exist, has been moved, or you don't have access to it.
      </p>
      <Link
        to="/"
        className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-zinc-400 disabled:pointer-events-none disabled:opacity-50 bg-zinc-50 text-zinc-900 hover:bg-zinc-50/90 h-10 px-4 py-2"
      >
        <HugeiconsIcon icon={ArrowLeft01Icon} className="mr-2 h-4 w-4" />
        Back to Dashboard
      </Link>
    </div>
  );
}
