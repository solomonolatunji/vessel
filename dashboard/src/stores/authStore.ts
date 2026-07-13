import { Store } from '@tanstack/store';
import type { User } from '#/interfaces';

interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
}

const isBrowser = typeof window !== 'undefined';

const getInitialAuth = (): AuthState => {
  if (!isBrowser) {
    return { token: null, user: null, isAuthenticated: false };
  }

  const storedToken = localStorage.getItem('vessl_auth_token');
  const storedUser = localStorage.getItem('vessl_auth_user');
  let parsedUser = null;
  if (storedUser) {
    try {
      parsedUser = JSON.parse(storedUser);
    } catch (e) {
      console.error('Failed to parse stored user:', e);
      localStorage.removeItem('vessl_auth_user');
    }
  }

  return {
    token: storedToken,
    user: parsedUser,
    isAuthenticated: !!storedToken,
  };
};

export const authStore = new Store<AuthState>(getInitialAuth());

if (isBrowser) {
  authStore.subscribe(() => {
    const state = authStore.state;
    if (state.token) {
      localStorage.setItem('vessl_auth_token', state.token);
    } else {
      localStorage.removeItem('vessl_auth_token');
    }

    if (state.user) {
      localStorage.setItem('vessl_auth_user', JSON.stringify(state.user));
    } else {
      localStorage.removeItem('vessl_auth_user');
    }
  });
}

export const authActions = {
  setAuth: (token: string, user: User) => {
    authStore.setState((state) => ({
      ...state,
      token,
      user,
      isAuthenticated: true,
    }));
  },
  logout: () => {
    authStore.setState((state) => ({
      ...state,
      token: null,
      user: null,
      isAuthenticated: false,
    }));
  },
};

import { useSelector } from '@tanstack/react-store';

export const useAuthState = () => {
  return useSelector(authStore, (state) => state);
};
