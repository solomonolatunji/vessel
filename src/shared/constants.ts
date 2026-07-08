export const API_ROUTES = {
  SYSTEM_INFO: '/api/system/info',
  SYSTEM_UPDATE: '/api/system/update',
  PROJECTS: '/api/projects',
  DEPLOY: '/api/projects/:id/deploy',
  LOGS_WS: '/ws/logs/:id',
  TERMINAL_WS: '/ws/terminal/:id',
} as const;

export const WEBSOCKET_EVENTS = {
  BUILD_LOG: 'deploy:build_log',
  BUILD_STATUS: 'deploy:status',
  CONTAINER_LOG: 'container:log',
  STATS_CPU_MEM: 'stats:cpu_mem',
} as const;
