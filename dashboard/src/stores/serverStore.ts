import { useSelector } from '@tanstack/react-store';
import { Store } from '@tanstack/store';

export interface CloudServer {
  id: string;
  name: string;
}

interface ServerState {
  activeServer: CloudServer | null;
  servers: CloudServer[];
}

const isBrowser = typeof window !== 'undefined';

const getInitialServer = (): ServerState => {
  if (!isBrowser) return { activeServer: null, servers: [] };

  try {
    const stored = localStorage.getItem('vessl_active_server');
    const activeServer = stored ? (JSON.parse(stored) as CloudServer) : null;
    return { activeServer, servers: [] };
  } catch {
    return { activeServer: null, servers: [] };
  }
};

export const serverStore = new Store<ServerState>(getInitialServer());

if (isBrowser) {
  serverStore.subscribe(() => {
    const { activeServer } = serverStore.state;
    if (activeServer) {
      localStorage.setItem('vessl_active_server', JSON.stringify(activeServer));
      localStorage.setItem('vessl_active_server_id', String(activeServer.id));
    } else {
      localStorage.removeItem('vessl_active_server');
      localStorage.removeItem('vessl_active_server_id');
    }
  });
}

export const serverActions = {
  setServers: (servers: CloudServer[]) => {
    serverStore.setState((s) => {
      const stillActive = servers.find((w) => String(w.id) === String(s.activeServer?.id));
      return {
        ...s,
        servers,
        activeServer: stillActive ?? servers[0] ?? null,
      };
    });
  },
  switchServer: (server: CloudServer) => {
    serverStore.setState((s) => ({ ...s, activeServer: server }));
    // When switching servers, we should clear the active workspace so it re-fetches
    // correctly for the new server
    localStorage.removeItem('vessl_active_workspace');
    window.location.reload();
  },
};

export const useServerState = () => {
  return useSelector(serverStore, (state) => state);
};
