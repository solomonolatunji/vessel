import { createFileRoute } from '@tanstack/react-router';
import { AISettings } from '#/features/instance/ai-settings';

export const Route = createFileRoute('/_dashboard/ai')({
  component: () => <AISettings />,
});
