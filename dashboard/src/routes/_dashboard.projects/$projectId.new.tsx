import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { Code2, Container, Database, GitBranch, LayoutTemplate, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { CreateDatabaseModal } from '#/features/databases/create-database-modal';
import { useListExampleApps, useListOneClickApps } from '#/hooks/useTemplates';

export const Route = createFileRoute('/_dashboard/projects/$projectId/new')({
  component: NewResourcePage,
});

function NewResourcePage() {
  const { projectId } = Route.useParams();
  const navigate = useNavigate();
  const [dbModalOpen, setDbModalOpen] = useState(false);

  const { data: oneClickResponse, isLoading: oneClickLoading } = useListOneClickApps();
  const { data: examplesResponse, isLoading: examplesLoading } = useListExampleApps();

  const templates = Array.isArray(oneClickResponse) ? oneClickResponse : [];
  const examples = Array.isArray(examplesResponse) ? examplesResponse : [];

  return (
    <div className="space-y-6 p-4">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <LayoutTemplate className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Add New Resource</h1>
            <p className="text-muted-foreground text-sm">
              Deploy a new application, database, or service to this project.
            </p>
          </div>
        </div>
      </div>

      <Tabs defaultValue="resources" className="w-full">
        <TabsList className="mb-6 grid w-full max-w-150 grid-cols-3">
          <TabsTrigger value="resources">Resources</TabsTrigger>
          <TabsTrigger value="one-click">One-Click Apps</TabsTrigger>
          <TabsTrigger value="examples">Example Projects</TabsTrigger>
        </TabsList>

        <TabsContent value="resources" className="space-y-4">
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
            <Card
              className="cursor-pointer transition-colors hover:border-primary/50"
              onClick={() => navigate({ to: '/sources' })}
            >
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <GitBranch className="h-5 w-5" />
                  GitHub Repository
                </CardTitle>
                <CardDescription>
                  Deploy source code from a public or private GitHub repository.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card
              className="cursor-pointer transition-colors hover:border-primary/50"
              onClick={() => setDbModalOpen(true)}
            >
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-5 w-5" />
                  Database
                </CardTitle>
                <CardDescription>
                  Provision a PostgreSQL, MySQL, Redis, or other database.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card
              className="cursor-pointer transition-colors hover:border-primary/50"
              onClick={() => alert('Docker Image deployment flow coming soon!')}
            >
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Container className="h-5 w-5" />
                  Docker Image
                </CardTitle>
                <CardDescription>
                  Deploy a pre-built Docker image from any public or private registry.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card
              className="cursor-pointer transition-colors hover:border-primary/50"
              onClick={() =>
                navigate({
                  to: '/projects/$projectId/compose',
                  params: { projectId },
                })
              }
            >
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Code2 className="h-5 w-5" />
                  Docker Compose
                </CardTitle>
                <CardDescription>
                  Deploy multiple services defined in a docker-compose.yml file.
                </CardDescription>
              </CardHeader>
            </Card>
          </div>
        </TabsContent>

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

      <CreateDatabaseModal
        isOpen={dbModalOpen}
        onOpenChange={setDbModalOpen}
        projectId={projectId}
      />
    </div>
  );
}
