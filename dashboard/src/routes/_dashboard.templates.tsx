import { createFileRoute } from '@tanstack/react-router';
import { LayoutTemplate, Loader2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { useListOneClickApps } from '#/hooks/useTemplates';

export const Route = createFileRoute('/_dashboard/templates')({
  component: TemplatesPage,
});

function TemplatesPage() {
  const { data: response, isLoading } = useListOneClickApps();
  const templates = Object.values(response?.data || {});

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <LayoutTemplate className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Templates</h1>
            <p className="text-muted-foreground text-sm">
              Deploy pre-configured applications with one click.
            </p>
          </div>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : templates.length === 0 ? (
        <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card/40">
          <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
            <LayoutTemplate className="h-5 w-5 text-primary" />
          </div>
          <h3 className="mt-4 font-bold text-foreground text-lg tracking-tight">No templates</h3>
          <p className="mt-1 max-w-sm text-center text-muted-foreground text-sm">
            There are no one-click templates available at this time.
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {templates.map((template: any) => (
            <Card key={template.id}>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  {template.logo && (
                    <img src={template.logo} alt={template.name} className="h-6 w-6" />
                  )}
                  {template.name}
                </CardTitle>
                <CardDescription>{template.description}</CardDescription>
              </CardHeader>
              <CardContent>
                <Button className="w-full">Deploy</Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
