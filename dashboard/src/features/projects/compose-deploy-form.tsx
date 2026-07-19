import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';

export function ComposeDeployForm({ projectId }: { projectId: string }) {
  const [composeText, setComposeText] = useState('');

  const handleDeploy = () => {
    console.log('Deploying docker-compose for project:', projectId);
    console.log(composeText);
    // integrate with actual API
    alert('Deploy triggered! Check console for details.');
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Deploy via Docker Compose</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        <textarea
          className="h-64 w-full rounded border p-2 font-mono text-sm"
          placeholder="version: '3'\nservices:\n  web:\n    image: nginx"
          value={composeText}
          onChange={(e) => setComposeText(e.target.value)}
        />
        <Button onClick={handleDeploy}>Deploy Compose</Button>
      </CardContent>
    </Card>
  );
}
