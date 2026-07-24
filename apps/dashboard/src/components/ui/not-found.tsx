import { Link } from '@tanstack/react-router';
import { LayoutDashboard } from 'lucide-react';
import { Button } from './button';

export function NotFoundComponent() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background p-4 text-foreground">
      <div className="flex max-w-md flex-col items-center space-y-6 text-center">
        <div className="font-bold text-8xl text-muted-foreground/30 tracking-tighter">404</div>
        <div className="space-y-2">
          <h1 className="font-bold text-2xl tracking-tight sm:text-3xl">Page not found</h1>
          <p className="text-muted-foreground">
            Sorry, we couldn't find the page you're looking for. The link might be broken, or the
            page may have been removed.
          </p>
        </div>

        <div className="flex gap-4 pt-4">
          <Button asChild variant="default">
            <Link to="/">
              <LayoutDashboard className="mr-2 h-4 w-4" />
              Go to Dashboard
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
