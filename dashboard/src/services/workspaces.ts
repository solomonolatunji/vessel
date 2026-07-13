import type { BaseResponse } from '#/interfaces/base';
import type {
  AuditLog,
  CreateSSHKeyRequest,
  CreateTrustedDomainRequest,
  CreateWorkspaceRequest,
  CreateWorkspaceResponse,
  GetWorkspaceResponse,
  InviteMemberRequest,
  ListWorkspacesResponse,
  SSHKey,
  TrustedDomain,
  UpdateWorkspaceRequest,
  UpdateWorkspaceResponse,
  WorkspaceInvite,
  WorkspaceMember,
} from '#/interfaces/workspace';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const workspacesService = {
  /**
   * Retrieves all workspaces for the currently authenticated context.
   */
  listWorkspaces: async (): Promise<ListWorkspacesResponse> => {
    try {
      return await apiClient.get<ListWorkspacesResponse>('/workspaces');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getWorkspace: async (id: string): Promise<GetWorkspaceResponse> => {
    try {
      return await apiClient.get<GetWorkspaceResponse>(`/workspaces/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createWorkspace: async (payload: CreateWorkspaceRequest): Promise<CreateWorkspaceResponse> => {
    try {
      return await apiClient.post<CreateWorkspaceResponse>('/workspaces', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateWorkspace: async (
    id: string,
    payload: UpdateWorkspaceRequest
  ): Promise<UpdateWorkspaceResponse> => {
    try {
      return await apiClient.put<UpdateWorkspaceResponse>(`/workspaces/${id}`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteWorkspace: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/workspaces/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listMembers: async (id: string): Promise<BaseResponse<WorkspaceMember[]>> => {
    try {
      return await apiClient.get<BaseResponse<WorkspaceMember[]>>(`/workspaces/${id}/members`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  inviteMember: async (
    id: string,
    payload: InviteMemberRequest
  ): Promise<BaseResponse<WorkspaceInvite>> => {
    try {
      return await apiClient.post<BaseResponse<WorkspaceInvite>>(
        `/workspaces/${id}/invite`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  removeMember: async (id: string, userId: string): Promise<void> => {
    try {
      await apiClient.delete(`/workspaces/${id}/members/${userId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getInvite: async (token: string): Promise<BaseResponse<WorkspaceInvite>> => {
    try {
      return await apiClient.get<BaseResponse<WorkspaceInvite>>(`/workspace-invites/${token}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  acceptInvite: async (token: string): Promise<BaseResponse<void>> => {
    try {
      return await apiClient.post<BaseResponse<void>>(`/workspace-invites/${token}/accept`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listTrustedDomains: async (id: string): Promise<BaseResponse<TrustedDomain[]>> => {
    try {
      return await apiClient.get<BaseResponse<TrustedDomain[]>>(
        `/workspaces/${id}/trusted-domains`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createTrustedDomain: async (
    id: string,
    payload: CreateTrustedDomainRequest
  ): Promise<BaseResponse<TrustedDomain>> => {
    try {
      return await apiClient.post<BaseResponse<TrustedDomain>>(
        `/workspaces/${id}/trusted-domains`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteTrustedDomain: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/trusted-domains/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listSSHKeys: async (id: string): Promise<BaseResponse<SSHKey[]>> => {
    try {
      return await apiClient.get<BaseResponse<SSHKey[]>>(`/workspaces/${id}/ssh-keys`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createSSHKey: async (id: string, payload: CreateSSHKeyRequest): Promise<BaseResponse<SSHKey>> => {
    try {
      return await apiClient.post<BaseResponse<SSHKey>>(`/workspaces/${id}/ssh-keys`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteSSHKey: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/ssh-keys/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listAuditLogs: async (id: string): Promise<BaseResponse<AuditLog[]>> => {
    try {
      return await apiClient.get<BaseResponse<AuditLog[]>>(`/workspaces/${id}/audit-logs`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
