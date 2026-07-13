import { useSelector } from '@tanstack/react-store';
import { Store } from '@tanstack/store';
import type { Workspace } from '#/interfaces/workspace';

interface WorkspaceState {
  activeWorkspace: Workspace | null;
  workspaces: Workspace[];
}

const isBrowser = typeof window !== 'undefined';

const getInitialWorkspace = (): WorkspaceState => {
  if (!isBrowser) return { activeWorkspace: null, workspaces: [] };

  try {
    const stored = localStorage.getItem('vessl_active_workspace');
    const activeWorkspace = stored ? (JSON.parse(stored) as Workspace) : null;
    return { activeWorkspace, workspaces: [] };
  } catch {
    return { activeWorkspace: null, workspaces: [] };
  }
};

export const workspaceStore = new Store<WorkspaceState>(getInitialWorkspace());

if (isBrowser) {
  workspaceStore.subscribe(() => {
    const { activeWorkspace } = workspaceStore.state;
    if (activeWorkspace) {
      localStorage.setItem('vessl_active_workspace', JSON.stringify(activeWorkspace));
    } else {
      localStorage.removeItem('vessl_active_workspace');
    }
  });
}

export const workspaceActions = {
  setWorkspaces: (workspaces: Workspace[]) => {
    workspaceStore.setState((s) => {
      const stillActive = workspaces.find((w) => w.id === s.activeWorkspace?.id);
      return {
        ...s,
        workspaces,
        activeWorkspace: stillActive ?? workspaces[0] ?? null,
      };
    });
  },
  switchWorkspace: (workspace: Workspace) => {
    workspaceStore.setState((s) => ({ ...s, activeWorkspace: workspace }));
  },
};

export const useWorkspaceState = () => {
  return useSelector(workspaceStore, (state) => state);
};
