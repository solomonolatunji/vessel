export interface ContainerHealth {
  status: 'running' | 'stopped' | 'restarting' | 'error' | 'building';
  cpuUsagePercentage: number;
  memoryUsageBytes: number;
  memoryLimitBytes: number;
  uptimeSeconds: number;
}

export interface ProjectConfig {
  id: string;
  name: string;
  repositoryUrl?: string;
  branch?: string;
  dockerfilePath?: string;
  domain?: string;
  envVarsCount: number;
  health: ContainerHealth;
  createdAt: string;
  updatedAt: string;
}

export interface SystemInfo {
  version: string;
  goVersion: string;
  dockerVersion: string;
  caddyVersion: string;
  os: string;
  arch: string;
  totalMemoryMB: number;
  freeMemoryMB: number;
  cpuCores: number;
  updateAvailable: boolean;
  latestVersion?: string;
}
