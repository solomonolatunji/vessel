import { z } from 'zod';

export const setupSchema = z
  .object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    email: z.string().email('Please enter a valid email address'),
    password: z.string().min(8, 'Password must be at least 8 characters long'),
    confirmPassword: z.string(),

    githubAppId: z.string().optional(),
    githubClientId: z.string().optional(),
    githubClientSecret: z.string().optional(),
    githubPrivateKey: z.string().optional(),
    githubWebhookSecret: z.string().optional(),

    dashboardDomain: z.string().optional(),
    defaultWildcardDomain: z.string().optional(),

    s3Skip: z.boolean().optional(),
    s3AccountId: z.string().optional(),
    s3Bucket: z.string().optional(),
    s3AccessKeyId: z.string().optional(),
    s3SecretAccessKey: z.string().optional(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ['confirmPassword'],
  });

export type SetupSchema = z.infer<typeof setupSchema>;
