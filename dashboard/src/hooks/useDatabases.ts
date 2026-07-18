import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { databasesService } from '#/services/databases';

export const useListDatabases = () => {
  return useQuery({
    queryKey: ['databases', 'listDatabases'].filter(Boolean),
    queryFn: () => databasesService.listDatabases(),
  });
};

export const useGetDatabase = (id: string) => {
  return useQuery({
    queryKey: ['databases', 'getDatabase', id].filter(Boolean),
    queryFn: () => databasesService.getDatabase(id),
  });
};

export const useCreateDatabase = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof databasesService.createDatabase>[0] }) =>
      databasesService.createDatabase(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['databases'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useDeleteDatabase = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => databasesService.deleteDatabase(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['databases'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useStartDatabase = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => databasesService.startDatabase(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['databases'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useStopDatabase = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => databasesService.stopDatabase(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['databases'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useRestartDatabase = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => databasesService.restartDatabase(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['databases'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useQueryDatabase = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      payload: Parameters<typeof databasesService.queryDatabase>[1];
    }) => databasesService.queryDatabase(payload.id, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['databases'] });
    },
  });
};

export const useGetSchemas = (id: string) => {
  return useQuery({
    queryKey: ['databases', 'getSchemas', id].filter(Boolean),
    queryFn: () => databasesService.getSchemas(id),
  });
};

export const useGetTableData = (
  id: string,
  table: string,
  params?: { limit?: number; offset?: number }
) => {
  return useQuery({
    queryKey: ['databases', 'getTableData', id, table, params].filter(Boolean),
    queryFn: () => databasesService.getTableData(id, table, params),
  });
};

export const useInsertTableRow = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      table: string;
      payload: Parameters<typeof databasesService.insertTableRow>[2];
    }) => databasesService.insertTableRow(payload.id, payload.table, payload.payload),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: ['databases', 'getTableData', variables.id, variables.table],
      });
    },
  });
};

export const useUpdateTableRow = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      table: string;
      payload: Parameters<typeof databasesService.updateTableRow>[2];
    }) => databasesService.updateTableRow(payload.id, payload.table, payload.payload),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: ['databases', 'getTableData', variables.id, variables.table],
      });
    },
  });
};

export const useDeleteTableRow = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      table: string;
      payload?: Parameters<typeof databasesService.deleteTableRow>[2];
    }) => databasesService.deleteTableRow(payload.id, payload.table, payload.payload),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: ['databases', 'getTableData', variables.id, variables.table],
      });
    },
  });
};

export const useImportDatabase = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      payload: Parameters<typeof databasesService.importDatabase>[1];
    }) => databasesService.importDatabase(payload.id, payload.payload),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({ queryKey: ['databases', 'getSchemas', variables.id] });
      await queryClient.invalidateQueries({
        queryKey: ['databases', 'getTableData', variables.id],
      });
    },
  });
};
