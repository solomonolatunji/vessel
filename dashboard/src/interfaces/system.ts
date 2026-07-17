export interface RailwayImportRequest {
  token: string;
  projectId: string;
  excludeRailwayVars: boolean;
  recreateDatabases: boolean;
  importData: boolean;
}

export interface RailwayProject {
  id: string;
  name: string;
  description: string;
  environments: RailwayEnvironmentConnection;
  services: RailwayServiceConnection;
}

export interface RailwayEnvironmentConnection {
  edges: RailwayEnvironmentEdge[];
}

export interface RailwayEnvironmentEdge {
  node: RailwayEnvironment;
}

export interface RailwayEnvironment {
  id: string;
  name: string;
}

export interface RailwayServiceConnection {
  edges: RailwayServiceEdge[];
}

export interface RailwayServiceEdge {
  node: RailwayService;
}

export interface RailwayService {
  id: string;
  name: string;
  source: RailwayServiceSource;
}

export interface RailwayServiceSource {
  image: string;
  repo: string;
}

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
