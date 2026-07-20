import { PostHogProvider as PHProvider } from '@posthog/react';
import posthog from 'posthog-js';
import type React from 'react';
import { env } from '#/env';

const posthogKey = env.VITE_PUBLIC_POSTHOG_KEY;
const posthogHost = env.VITE_PUBLIC_POSTHOG_HOST || 'https://us.i.posthog.com';

if (posthogKey) {
  posthog.init(posthogKey, {
    api_host: posthogHost,
    autocapture: true,
    capture_pageview: true,
    capture_pageleave: true,
  });
}

export function PostHogProvider({ children }: { children: React.ReactNode }) {
  if (!posthogKey) {
    return <>{children}</>;
  }

  return <PHProvider client={posthog}>{children}</PHProvider>;
}
