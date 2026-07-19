import { createFileRoute } from '@tanstack/react-router';
import { Code2, LayoutTemplate, Loader2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { useListExampleApps, useListOneClickApps } from '#/hooks/useTemplates';

export const Route = createFileRoute('/_dashboard/templates')({
  component: TemplatesPage,
});

function TemplatesPage() {
  const { data: oneClickResponse, isLoading: oneClickLoading } = useListOneClickApps();
  const { data: examplesResponse, isLoading: examplesLoading } = useListExampleApps();

  const templates = Array.isArray(oneClickResponse) ? oneClickResponse : [];
  const examples = Array.isArray(examplesResponse) ? examplesResponse : [];

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <LayoutTemplate className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Templates & Examples</h1>
            <p className="text-muted-foreground text-sm">
              Deploy pre-configured applications with one click.
            </p>
          </div>
        </div>
      </div>

      <Tabs defaultValue="one-click" className="w-full">
        <TabsList className="mb-6 grid w-full max-w-[400px] grid-cols-2">
          <TabsTrigger value="one-click">One-Click Apps</TabsTrigger>
          <TabsTrigger value="examples">Example Projects</TabsTrigger>
        </TabsList>

        <TabsContent value="one-click" className="space-y-4">
          {oneClickLoading ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : templates.length === 0 ? (
            <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card/40">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
                <LayoutTemplate className="h-5 w-5 text-primary" />
              </div>
              <h3 className="mt-4 font-bold text-foreground text-lg tracking-tight">
                No templates
              </h3>
              <p className="mt-1 max-w-sm text-center text-muted-foreground text-sm">
                There are no one-click templates available at this time.
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
              {templates.map((template) => (
                <Card key={template.id}>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      {(template as any).logo && (
                        <img src={(template as any).logo} alt={template.name} className="h-6 w-6" />
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
        </TabsContent>

        <TabsContent value="examples" className="space-y-4">
          {examplesLoading ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : examples.length === 0 ? (
            <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card/40">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
                <Code2 className="h-5 w-5 text-primary" />
              </div>
              <h3 className="mt-4 font-bold text-foreground text-lg tracking-tight">No examples</h3>
              <p className="mt-1 max-w-sm text-center text-muted-foreground text-sm">
                Could not load examples from TechX/vessl-examples.
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
              {examples.map((example) => (
                <Card key={example.id}>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Code2 className="h-5 w-5" />
                      {example.name}
                    </CardTitle>
                    <CardDescription>{example.description}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <Button
                      className="w-full"
                      variant="outline"
                      onClick={() => window.open(example.repo, '_blank')}
                    >
                      View on GitHub
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
