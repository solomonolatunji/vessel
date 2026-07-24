import { createFileRoute } from '@tanstack/react-router';
import { NotFoundComponent } from '#/components/ui/not-found';

export const Route = createFileRoute('/$')({
  component: NotFoundComponent,
});
