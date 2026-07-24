import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { Switch } from '#/components/ui/switch';
import type { CreateDatabaseResponse, DatabaseEngine } from '#/interfaces/database';
import { apiClient } from '#/lib/apiClient';

const engines: { label: string; value: DatabaseEngine }[] = [
  { label: 'PostgreSQL', value: 'postgres' },
  { label: 'TimescaleDB', value: 'timescaledb' },
  { label: 'MySQL', value: 'mysql' },
  { label: 'MariaDB', value: 'mariadb' },
  { label: 'ClickHouse', value: 'clickhouse' },
  { label: 'MongoDB', value: 'mongodb' },
  { label: 'Redis', value: 'redis' },
  { label: 'Dragonfly', value: 'dragonfly' },
  { label: 'KeyDB', value: 'keydb' },
  { label: 'Kafka', value: 'kafka' },
  { label: 'RabbitMQ', value: 'rabbitmq' },
  { label: 'NATS', value: 'nats' },
];

const schema = z.object({
  name: z.string().min(1, 'Name is required'),
  engine: z.string().min(1, 'Engine is required'),
  version: z.string().min(1, 'Version is required'),
  port: z.number().min(1, 'Port is required'),
  username: z.string(),
  password: z.string(),
  databaseName: z.string(),
  volumePath: z.string(),
  customArgs: z.string(),
  logicalReplication: z.boolean(),
});

type FormData = z.infer<typeof schema>;

interface Props {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  projectId?: string;
  environmentId?: string;
}

export function CreateDatabaseModal({
  isOpen,
  onOpenChange,
  projectId = '',
  environmentId = '',
}: Props) {
  const queryClient = useQueryClient();
  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      engine: 'postgres',
      version: 'latest',
      port: 5432,
      username: 'postgres',
      password: '',
      databaseName: 'codedock',
      volumePath: '/var/lib/postgresql/data',
      customArgs: '',
      logicalReplication: false,
    },
  });

  const engineValue = watch('engine');
  const logicalReplicationValue = watch('logicalReplication');

  const createMutation = useMutation({
    mutationFn: (data: FormData) => {
      return apiClient.post<CreateDatabaseResponse>('/databases', {
        ...data,
        projectId,
        environmentId,
      });
    },
    onSuccess: () => {
      toast.success('Database created successfully');
      queryClient.invalidateQueries({ queryKey: ['databases'] });
      onOpenChange(false);
    },
  });

  const onSubmit = (data: FormData) => {
    createMutation.mutate(data);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[80vh] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Provision Database</DialogTitle>
          <DialogDescription>Select a database engine to spin up a new instance.</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="mt-4 space-y-4">
          <div className="space-y-2">
            <Label>Database Engine</Label>
            <Select
              value={engineValue}
              onValueChange={(val) => setValue('engine', val as DatabaseEngine)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select an engine" />
              </SelectTrigger>
              <SelectContent>
                {engines.map((engine) => (
                  <SelectItem key={engine.value} value={engine.value}>
                    {engine.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {errors.engine && <p className="text-red-500 text-sm">{errors.engine.message}</p>}
          </div>

          <div className="space-y-2">
            <Label>Name</Label>
            <Input {...register('name')} placeholder="my-postgres-db" />
            {errors.name && <p className="text-red-500 text-sm">{errors.name.message}</p>}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Version Tag</Label>
              <Input {...register('version')} placeholder="latest" />
            </div>
            <div className="space-y-2">
              <Label>Internal Port</Label>
              <Input type="number" {...register('port', { valueAsNumber: true })} />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Username</Label>
              <Input {...register('username')} />
            </div>
            <div className="space-y-2">
              <Label>Password</Label>
              <Input type="password" {...register('password')} />
            </div>
          </div>

          <div className="space-y-2">
            <Label>Initial Database Name</Label>
            <Input {...register('databaseName')} />
          </div>

          <div className="space-y-2">
            <Label>Volume Path</Label>
            <Input {...register('volumePath')} />
          </div>

          {(engineValue === 'postgres' || engineValue === 'timescaledb') && (
            <div className="mt-4 flex items-center justify-between rounded-md border p-3">
              <div className="space-y-0.5">
                <Label>Logical Replication (CDC)</Label>
                <p className="text-muted-foreground text-xs">
                  Enables max_replication_slots=10 and WAL retention for Change Data Capture tools.
                </p>
              </div>
              <Switch
                checked={logicalReplicationValue}
                onCheckedChange={(val) => setValue('logicalReplication', val)}
              />
            </div>
          )}

          <div className="flex justify-end space-x-2 pt-4">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Provisioning...' : 'Deploy Database'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
