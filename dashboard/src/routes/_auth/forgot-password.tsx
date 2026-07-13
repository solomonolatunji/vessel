import { createFileRoute } from '@tanstack/react-router';
import { ForgotPasswordForm } from '#/features/auth/forgot-password-form';

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPasswordPage,
});

function ForgotPasswordPage() {
  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col space-y-2 text-center mb-8">
        <h1 className="text-3xl font-semibold tracking-tight text-foreground">
          Reset your password
        </h1>
        <p className="text-sm text-muted-foreground">Enter your email to receive a reset link.</p>
      </div>

      <ForgotPasswordForm />
    </div>
  );
}
