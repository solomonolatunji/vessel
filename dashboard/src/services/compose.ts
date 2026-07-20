import { apiClient } from '#/lib/apiClient';

export interface ComposeAnalyzeRequest {
  projectId: string;
  composeContent: string;
}

export interface ComposeAnalyzeResponse {
  appServices: any[];
  databases: any[];
}

class ComposeService {
  async analyze(req: ComposeAnalyzeRequest): Promise<ComposeAnalyzeResponse> {
    const response = await apiClient.post<{ data: ComposeAnalyzeResponse }>(
      '/compose/analyze',
      req
    );
    return response.data;
  }

  async deploy(projectId: string, composeContent: string): Promise<any> {
    const formData = new FormData();
    formData.append('projectId', projectId);

    const blob = new Blob([composeContent], { type: 'text/yaml' });
    formData.append('file', blob, 'docker-compose.yml');

    const response = await apiClient.post<{ data: any }>('/compose/deploy', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  }
}

export const composeService = new ComposeService();
