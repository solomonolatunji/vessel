import { createFileRoute } from '@tanstack/react-router';
import { EnvironmentCanvas } from '#/features/canvas/environment-canvas';
import { useGetCanvasSummary, useGetEnvironmentCanvas } from '#/hooks/useCanvas';

export const Route = createFileRoute('/_dashboard/projects/$projectId/canvas')({
  component: CanvasRouteComponent,
});

function CanvasRouteComponent() {
  const { projectId } = Route.useParams();
  const { data: summaryRes, isLoading: summaryLoading } = useGetCanvasSummary(projectId);

  const envId = summaryRes?.data?.defaultEnvironment?.id;
  const { data: envRes, isLoading: envLoading } = useGetEnvironmentCanvas(envId || '');

  if (summaryLoading || (envId && envLoading)) {
    return <div className="p-4">Loading canvas...</div>;
  }

  if (!envRes?.data) {
    return <div className="p-4">No environment data found for this project.</div>;
  }

  return (
    <div className="flex h-full flex-col gap-4 p-4">
      <h1 className="font-bold text-2xl">Environment Canvas</h1>
      <div className="flex-1 rounded border bg-white shadow-sm">
        <EnvironmentCanvas envData={envRes.data} />
      </div>
    </div>
  );
}
