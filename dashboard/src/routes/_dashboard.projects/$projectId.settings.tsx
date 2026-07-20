import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/projects/$projectId/settings')({
  component: SettingsRouteComponent,
});

function SettingsRouteComponent() {
  return (
    <div className="flex flex-col gap-4 p-4">
      <h1 className="font-bold text-2xl">Project Settings</h1>

      <div className="grid grid-cols-1 gap-4">
        {/* Project-level settings will go here. Domains belong to services. */}
        <p className="text-gray-500 text-sm">Project settings are under construction.</p>
      </div>
    </div>
  );
}
