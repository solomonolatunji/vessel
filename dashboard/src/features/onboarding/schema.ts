import { z } from 'zod';

export const setupSchema = z
  .object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    email: z.string().email('Please enter a valid email address'),
    password: z.string().min(8, 'Password must be at least 8 characters long'),
    confirmPassword: z.string(),

    env: z.object({
      jwtSecret: z
        .string()
        .refine((val) => !val || val.length >= 32, {
          message: 'Secret must be at least 32 characters if provided',
        })
        .optional(),
      dataDir: z.string().min(1, 'Data directory is required'),
      dashboardUrl: z.string().url('Must be a valid URL'),
      port: z.number().min(1).max(65535),
    }),

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
