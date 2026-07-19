export interface SystemStats {
  cpu: CPUStats;
  disk: DiskStats;
  memory: MemoryStats;
  loadAvg: number[];
  processes: number;
  uptimeSeconds: number;
  docker: DockerStats;
  backups?: { size: string };
  packageCache?: { size: string };
  systemLogs?: { size: string };
}

export interface DockerStats {
  images: DockerLayerStat;
  containers: DockerLayerStat;
  volumes: DockerLayerStat;
  buildCache: DockerLayerStat;
  reclaimableGb: number;
}

export interface DockerLayerStat {
  active: string;
  totalCount: string;
  size: string;
  reclaimable: string;
}

export interface CPUStats {
  cores: number;
  percent: number;
}

export interface DiskStats {
  freeGb: number;
  percent: number;
  totalGb: number;
  usedGb: number;
}

export interface MemoryStats {
  freeMb: number;
  percent: number;
  totalMb: number;
  usedMb: number;
}
