import { Copy, Eye, EyeOff } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { Database } from '#/interfaces/database';

interface Props {
  database: Database;
}

export function DatabaseConnectionCard({ database }: Props) {
  const [showPassword, setShowPassword] = useState(false);

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    toast.success(`${label} copied to clipboard`);
  };

  const getConnectionUrl = () => {
    const { engine, username, password, internalDns, port, databaseName } = database;
    let scheme = engine;
    if (engine === 'postgresql') scheme = 'postgres';
    if (engine === 'mongodb') scheme = 'mongodb';
    if (engine === 'redis' || engine === 'dragonfly' || engine === 'keydb') scheme = 'redis';

    // Simplistic connection string builder
    if (engine === 'redis' || engine === 'dragonfly' || engine === 'keydb') {
      return `${scheme}://${password ? `:${password}@` : ''}${internalDns}:${port}`;
    }

    return `${scheme}://${username}:${password}@${internalDns}:${port}/${databaseName}`;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Connection Details</CardTitle>
        <CardDescription>
          Use these credentials to connect to your database instance internally.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label>Connection URL</Label>
          <div className="flex space-x-2">
            <Input
              type={showPassword ? 'text' : 'password'}
              value={getConnectionUrl()}
              readOnly
              className="bg-muted/50 font-mono text-sm"
            />
            <Button variant="outline" size="icon" onClick={() => setShowPassword(!showPassword)}>
              {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </Button>
            <Button
              variant="outline"
              size="icon"
              onClick={() => copyToClipboard(getConnectionUrl(), 'Connection URL')}
            >
              <Copy className="h-4 w-4" />
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-4 pt-2 md:grid-cols-2">
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Host (Internal)</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">{database.internalDns}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.internalDns, 'Host')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Port</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">{database.port}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.port.toString(), 'Port')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Username</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">{database.username}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.username, 'Username')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Password</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">
                {showPassword ? database.password : '••••••••'}
              </span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.password, 'Password')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Database Name</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">{database.databaseName}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.databaseName, 'Database Name')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
