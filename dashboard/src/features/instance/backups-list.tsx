import { format } from 'date-fns';
import { Calendar, Check, Database, Download, History, Play, Trash2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Row, Section } from '#/components/ui/section';
import { Switch } from '#/components/ui/switch';
import {
  useCreate,
  useDelete,
  useDeleteRecord,
  useList,
  useListRecords,
  useTrigger,
  useUpdate,
} from '#/hooks/useBackups';

export function BackupsList() {
  const { data: configsData, isLoading } = useList();
  const configs = configsData?.data || [];
  const config = configs[0];

  const createBackup = useCreate();
  const updateBackup = useUpdate();
  const triggerBackup = useTrigger();
  const deleteBackup = useDelete();
  const deleteRecord = useDeleteRecord();

  const { data: recordsData, isLoading: isLoadingRecords } = useListRecords(config?.id || '');
  const records = recordsData?.data || [];

  const [name, setName] = useState('vessl-db');
  const [description, setDescription] = useState('Vessl database');
  const [dbUser, setDbUser] = useState('vessl');
  const [dbPassword, setDbPassword] = useState('********');
  const [backupEnabled, setBackupEnabled] = useState(true);
  const [s3Enabled, setS3Enabled] = useState(false);
  const [disableLocal, setDisableLocal] = useState(false);
  const [schedule, setSchedule] = useState('0 0 * * *');
  const [timezone, setTimezone] = useState('UTC');
  const [timeout, setTimeoutVal] = useState('3600');
  const [retentionDays, setRetentionDays] = useState('7');
  const [maxBackups, setMaxBackups] = useState('0');
  const [maxStorage, setMaxStorage] = useState('0');

  useEffect(() => {
    if (config) {
      setName(config.name);
      setDescription(config.description);
      setDbUser(config.dbUser);
      if (config.dbPassword) setDbPassword(config.dbPassword);
      setBackupEnabled(config.backupEnabled);
      setS3Enabled(config.s3Enabled);
      setDisableLocal(config.disableLocal);
      setSchedule(config.schedule);
      setTimezone(config.timezone);
      setTimeoutVal(config.timeout.toString());
      setRetentionDays(config.retentionDays.toString());
      setMaxBackups(config.maxBackups.toString());
      setMaxStorage(config.maxStorageGb.toString());
    }
  }, [config]);

  const handleSave = async (e?: React.FormEvent) => {
    e?.preventDefault();
    try {
      const payload = {
        projectId: 'global',
        name,
        description,
        dbUser,
        dbPassword: dbPassword === '********' ? '' : dbPassword,
        backupEnabled,
        s3Enabled,
        disableLocal,
        schedule,
        timezone,
        timeout: parseInt(timeout, 10),
        retentionDays: parseInt(retentionDays, 10),
        maxBackups: parseInt(maxBackups, 10),
        maxStorageGb: parseInt(maxStorage, 10),
      };

      if (config) {
        await updateBackup.mutateAsync({ id: config.id, payload });
      } else {
        await createBackup.mutateAsync({ payload });
      }
      toast.success('Backup configuration saved');
    } catch {
      toast.error('Failed to save backup configuration');
    }
  };

  const handleTrigger = async () => {
    if (!config) {
      toast.error('Please save the configuration first');
      return;
    }
    try {
      await triggerBackup.mutateAsync({ id: config.id });
      toast.success('Backup triggered successfully');
    } catch {
      toast.error('Failed to trigger backup');
    }
  };

  if (isLoading) {
    return <div className="p-6 text-muted-foreground">Loading backups...</div>;
  }

  return (
    <div className="space-y-6 pb-12">
      <div className="mb-5 flex items-center justify-between">
        <div>
          <h1 className="font-bold text-xl">System Backups</h1>
          <p className="text-muted-foreground text-sm">
            Backup configuration for the Vessl instance database.
          </p>
        </div>

        <div className="flex shrink-0 items-center gap-3">
          <Button
            variant="outline"
            size="sm"
            onClick={handleTrigger}
            disabled={triggerBackup.isPending || !config}
          >
            <Play className="mr-2 h-4 w-4" />
            {triggerBackup.isPending ? 'Triggering...' : 'Backup Now'}
          </Button>
          <Button
            size="sm"
            onClick={handleSave}
            disabled={createBackup.isPending || deleteBackup.isPending}
          >
            <Check className="mr-2 h-4 w-4" />
            {createBackup.isPending || deleteBackup.isPending ? 'Saving...' : 'Save Changes'}
          </Button>
        </div>
      </div>

      <Section icon={<Database className="h-4 w-4" />} title="Database Configuration">
        <Row label="UUID" description="The unique identifier for this backup configuration.">
          <Input disabled value={config?.id || 'Not configured yet (Save to generate)'} />
        </Row>
        <Row label="Name" description="A friendly name for this configuration.">
          <Input value={name} onChange={(e) => setName(e.target.value)} />
        </Row>
        <Row label="Description" description="Optional description of the database.">
          <Input value={description} onChange={(e) => setDescription(e.target.value)} />
        </Row>
        <Row label="Database User" description="The username used to connect to the database.">
          <Input value={dbUser} onChange={(e) => setDbUser(e.target.value)} />
        </Row>
        <Row label="Database Password" description="The password for the database user.">
          <Input
            type="password"
            value={dbPassword}
            onChange={(e) => setDbPassword(e.target.value)}
          />
        </Row>
      </Section>

      <Section icon={<Calendar className="h-4 w-4" />} title="Scheduled Backup">
        <Row label="Backup Enabled" description="Enable or disable scheduled backups globally.">
          <div className="flex items-center gap-2">
            <Switch checked={backupEnabled} onCheckedChange={setBackupEnabled} />
          </div>
        </Row>
        <Row label="S3 Enabled" description="Upload backups to the configured S3 destination.">
          <div className="flex items-center gap-2">
            <Switch checked={s3Enabled} onCheckedChange={setS3Enabled} />
          </div>
        </Row>
        <Row label="Disable Local Backup" description="Do not store backups on the local disk.">
          <div className="flex items-center gap-2">
            <Switch checked={disableLocal} onCheckedChange={setDisableLocal} />
          </div>
        </Row>
        <Row label="Frequency" description="Cron expression for the backup schedule.">
          <Input value={schedule} onChange={(e) => setSchedule(e.target.value)} />
        </Row>
        <Row label="Timezone" description="The timezone used for the cron expression.">
          <Input value={timezone} onChange={(e) => setTimezone(e.target.value)} disabled />
        </Row>
        <Row label="Timeout (seconds)" description="Maximum execution time before failing.">
          <Input value={timeout} onChange={(e) => setTimeoutVal(e.target.value)} disabled />
        </Row>
      </Section>

      <Section icon={<Trash2 className="h-4 w-4" />} title="Retention Settings">
        <div className="py-4 pb-6">
          <ul className="list-disc space-y-1 pl-5 text-muted-foreground text-sm">
            <li>Setting a value to 0 means unlimited retention.</li>
            <li>
              The retention rules work independently - whichever limit is reached first will trigger
              cleanup.
            </li>
          </ul>
        </div>
        <Row label="Number of backups to keep">
          <Input value={maxBackups} onChange={(e) => setMaxBackups(e.target.value)} disabled />
        </Row>
        <Row label="Days to keep backups">
          <Input value={retentionDays} onChange={(e) => setRetentionDays(e.target.value)} />
        </Row>
        <Row label="Maximum storage (GB)">
          <Input value={maxStorage} onChange={(e) => setMaxStorage(e.target.value)} disabled />
        </Row>
      </Section>

      <Section icon={<History className="h-4 w-4" />} title={`Executions (${records.length})`}>
        <div className="flex flex-col gap-4 py-4">
          <div className="flex flex-col gap-4">
            {isLoadingRecords ? (
              <div className="py-8 text-center text-muted-foreground">Loading executions...</div>
            ) : records.length === 0 ? (
              <div className="py-8 text-center text-muted-foreground">No executions yet.</div>
            ) : (
              records.map((record) => (
                <div
                  key={record.id}
                  className="flex flex-col gap-3 rounded-lg border border-border/50 bg-background/50 p-4"
                >
                  <Badge
                    variant="outline"
                    className={
                      record.status === 'completed'
                        ? 'w-fit border-green-500/20 bg-green-500/10 text-green-500'
                        : record.status === 'failed'
                          ? 'w-fit border-red-500/20 bg-red-500/10 text-red-500'
                          : 'w-fit border-yellow-500/20 bg-yellow-500/10 text-yellow-500'
                    }
                  >
                    {record.status === 'completed'
                      ? 'Success'
                      : record.status === 'failed'
                        ? 'Failed'
                        : 'Running'}
                  </Badge>

                  <div className="text-muted-foreground text-sm leading-relaxed">
                    {record.startedAt
                      ? format(new Date(record.startedAt), 'MMM d, HH:mm')
                      : 'Unknown time'}{' '}
                    • Database: vessl • Size: {(record.fileSizeBytes / 1024 / 1024).toFixed(2)} MB
                    <br />
                    Location: {record.filePath}
                  </div>

                  <div className="flex items-center gap-2 text-sm">
                    <span className="text-muted-foreground">Backup Availability:</span>
                    <Badge
                      variant="outline"
                      className="gap-1 border-green-500/20 bg-green-500/10 text-green-500"
                    >
                      <Check className="h-3 w-3" /> Local Storage
                    </Badge>
                  </div>

                  <div className="mt-2 flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      asChild
                      disabled={!record.s3Url && !record.filePath}
                    >
                      {record.s3Url ? (
                        <a href={record.s3Url} target="_blank" rel="noreferrer">
                          <Download className="mr-2 h-4 w-4" />
                          Download S3
                        </a>
                      ) : record.filePath ? (
                        <a
                          href={`${import.meta.env.VITE_API_URL}/backups/${config?.id}/records/${record.id}/download`}
                          target="_blank"
                          rel="noreferrer"
                        >
                          <Download className="mr-2 h-4 w-4" />
                          Download Local
                        </a>
                      ) : (
                        <span>
                          <Download className="mr-2 h-4 w-4" />
                          Download
                        </span>
                      )}
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() =>
                        deleteRecord.mutate({ id: config?.id || '', recordId: record.id })
                      }
                    >
                      <Trash2 className="mr-2 h-4 w-4" />
                      {deleteRecord.isPending ? 'Deleting...' : 'Delete'}
                    </Button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </Section>
    </div>
  );
}
