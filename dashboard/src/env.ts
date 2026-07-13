import { createEnv } from '@t3-oss/env-core';
import { z } from 'zod';

export const env = createEnv({
  server: {
    SERVER_URL: z.string().url().optional(),
  },

  clientPrefix: 'VITE_',

  client: {
    VITE_APP_TITLE: z.string().min(1).optional(),
    VITE_API_URL: z.string().url().optional().default('http://localhost:8080/api'),
    VITE_IS_CLOUD: z
      .enum(['true', 'false'])
      .default('false')
      .transform((v) => v === 'true'),
  },

  runtimeEnv: import.meta.env,

  emptyStringAsUndefined: true,
});
