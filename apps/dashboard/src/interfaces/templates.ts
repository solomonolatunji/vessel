export interface OneClickDeployRequest {
  appId: string;
  projectId: string;
  name: string;
}

export interface OneClickDeployResponse {
  success: boolean;
  message: string;
  serviceId?: string;
}

export interface ComposeDeployResponse {
  success: boolean;
  message: string;
  serviceIds?: string[];
}

export interface ArchiveDeployResponse {
  success: boolean;
  message: string;
  serviceId?: string;
}
