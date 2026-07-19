import type { BaseResponse } from '#/interfaces/base';
import type {
  AddMemberRequest,
  CreateTokenRequest,
  CreateWebhookRequest,
  ProjectMember,
  ProjectToken,
} from '#/interfaces/project';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export interface ServiceWebhook {
  id: string;
  serviceId: string;
  url: string;
  eventTypes: string[];
  includePrEnvironments: boolean;
  createdAt: string;
  updatedAt: string;
}

export const projectSettingsService = {
  listWebhooks: async (serviceId: string): Promise<BaseResponse<ServiceWebhook[]>> => {
    try {
      return await apiClient.get<BaseResponse<ServiceWebhook[]>>(`/apps/${serviceId}/webhooks`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createWebhook: async (
    serviceId: string,
    payload: CreateWebhookRequest
  ): Promise<BaseResponse<ServiceWebhook>> => {
    try {
      return await apiClient.post<BaseResponse<ServiceWebhook>>(
        `/apps/${serviceId}/webhooks`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteWebhook: async (serviceId: string, id: string): Promise<void> => {
    try {
      await apiClient.delete(`/apps/${serviceId}/webhooks/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listTokens: async (projectId: string): Promise<BaseResponse<ProjectToken[]>> => {
    try {
      return await apiClient.get<BaseResponse<ProjectToken[]>>(`/projects/${projectId}/tokens`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createToken: async (
    projectId: string,
    payload: CreateTokenRequest
  ): Promise<BaseResponse<{ token: string; projectToken: ProjectToken }>> => {
    try {
      return await apiClient.post<BaseResponse<{ token: string; projectToken: ProjectToken }>>(
        `/projects/${projectId}/tokens`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteToken: async (projectId: string, id: string): Promise<void> => {
    try {
      await apiClient.delete(`/projects/${projectId}/tokens/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listMembers: async (projectId: string): Promise<BaseResponse<ProjectMember[]>> => {
    try {
      return await apiClient.get<BaseResponse<ProjectMember[]>>(`/projects/${projectId}/members`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  addMember: async (
    projectId: string,
    payload: AddMemberRequest
  ): Promise<BaseResponse<ProjectMember>> => {
    try {
      return await apiClient.post<BaseResponse<ProjectMember>>(
        `/projects/${projectId}/members`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  removeMember: async (projectId: string, memberId: string): Promise<void> => {
    try {
      await apiClient.delete(`/projects/${projectId}/members/${memberId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
