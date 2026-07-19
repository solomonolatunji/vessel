-- Add missing columns to backup_configs and s3_destinations

ALTER TABLE backup_configs ADD COLUMN description TEXT;
ALTER TABLE backup_configs ADD COLUMN db_user TEXT;
ALTER TABLE backup_configs ADD COLUMN db_password TEXT;
ALTER TABLE backup_configs ADD COLUMN backup_enabled INTEGER DEFAULT 1;
ALTER TABLE backup_configs ADD COLUMN s3_enabled INTEGER DEFAULT 0;
ALTER TABLE backup_configs ADD COLUMN disable_local INTEGER DEFAULT 0;
ALTER TABLE backup_configs ADD COLUMN timezone TEXT DEFAULT 'UTC';
ALTER TABLE backup_configs ADD COLUMN timeout INTEGER DEFAULT 3600;
ALTER TABLE backup_configs ADD COLUMN max_backups INTEGER DEFAULT 0;
ALTER TABLE backup_configs ADD COLUMN max_storage_gb INTEGER DEFAULT 0;

ALTER TABLE s3_destinations ADD COLUMN description TEXT DEFAULT '';
ALTER TABLE s3_destinations ADD COLUMN provider TEXT DEFAULT 's3';
