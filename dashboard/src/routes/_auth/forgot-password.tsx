import { createFileRoute } from '@tanstack/react-router';
import { ForgotPasswordForm } from '#/features/auth/forgot-password-form';

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPasswordPage,
});

import { AlertCircle } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '#/components/ui/alert';
import { useGetPublicSettings } from '#/hooks/useSettings';

function ForgotPasswordPage() {
  const { data, isLoading } = useGetPublicSettings();
  const emailEnabled = data?.data?.emailEnabled;

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col space-y-2 text-center mb-8">
        <h1 className="text-3xl font-semibold tracking-tight text-foreground">
          Reset your password
        </h1>
        <p className="text-sm text-muted-foreground">Enter your email to receive a reset link.</p>
      </div>

      {!isLoading && emailEnabled === false ? (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Email not configured</AlertTitle>
          <AlertDescription>
            Your team is yet to set or enable email. Please contact your administrator.
          </AlertDescription>
        </Alert>
      ) : (
        <ForgotPasswordForm />
      )}
    </div>
  );
}
