import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { workspacesService } from '#/services/workspaces';
import { workspaceActions } from '#/stores/workspaceStore';

export const useListWorkspaces = () => {
  return useQuery({
    queryKey: ['workspaces', 'listWorkspaces'].filter(Boolean),
    queryFn: async () => {
      const res = await workspacesService.listWorkspaces();
      if (res?.data) {
        workspaceActions.setWorkspaces(res.data);
      }
      return res;
    },
  });
};

export const useGetWorkspace = (id: string) => {
  return useQuery({
    queryKey: ['workspaces', 'getWorkspace', id].filter(Boolean),
    queryFn: () => workspacesService.getWorkspace(id),
  });
};

export const useCreateWorkspace = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof workspacesService.createWorkspace>[0] }) =>
      workspacesService.createWorkspace(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useUpdateWorkspace = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      payload: Parameters<typeof workspacesService.updateWorkspace>[1];
    }) => workspacesService.updateWorkspace(payload.id, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useDeleteWorkspace = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => workspacesService.deleteWorkspace(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useListMembers = (id: string) => {
  return useQuery({
    queryKey: ['workspaces', 'listMembers', id].filter(Boolean),
    queryFn: () => workspacesService.listMembers(id),
  });
};

export const useInviteMember = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      payload: Parameters<typeof workspacesService.inviteMember>[1];
    }) => workspacesService.inviteMember(payload.id, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useRemoveMember = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; userId: string }) =>
      workspacesService.removeMember(payload.id, payload.userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useGetInvite = (token: string) => {
  return useQuery({
    queryKey: ['workspaces', 'getInvite', token].filter(Boolean),
    queryFn: () => workspacesService.getInvite(token),
  });
};

export const useAcceptInvite = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { token: string }) => workspacesService.acceptInvite(payload.token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useListTrustedDomains = (id: string) => {
  return useQuery({
    queryKey: ['workspaces', 'listTrustedDomains', id].filter(Boolean),
    queryFn: () => workspacesService.listTrustedDomains(id),
  });
};

export const useCreateTrustedDomain = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      payload: Parameters<typeof workspacesService.createTrustedDomain>[1];
    }) => workspacesService.createTrustedDomain(payload.id, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useDeleteTrustedDomain = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => workspacesService.deleteTrustedDomain(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useListSSHKeys = (id: string) => {
  return useQuery({
    queryKey: ['workspaces', 'listSSHKeys', id].filter(Boolean),
    queryFn: () => workspacesService.listSSHKeys(id),
  });
};

export const useCreateSSHKey = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      payload: Parameters<typeof workspacesService.createSSHKey>[1];
    }) => workspacesService.createSSHKey(payload.id, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useDeleteSSHKey = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => workspacesService.deleteSSHKey(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
  });
};

export const useListAuditLogs = (id: string) => {
  return useQuery({
    queryKey: ['workspaces', 'listAuditLogs', id].filter(Boolean),
    queryFn: () => workspacesService.listAuditLogs(id),
  });
};
