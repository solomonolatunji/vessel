import { createFileRoute } from '@tanstack/react-router';
import { User } from 'lucide-react';
import { AccessTokensList } from '#/features/profile/access-tokens-list';
import { Security2FASetup } from '#/features/profile/security-2fa-setup';
import {
  ProfileEmailForm,
  ProfileNameForm,
  ProfilePasswordForm,
} from '#/features/profile/user-profile-form';

export const Route = createFileRoute('/_dashboard/profile')({
  component: ProfilePage,
});

function ProfilePage() {
  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <User className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Profile & Security</h1>
            <p className="text-muted-foreground text-sm">
              Manage your personal profile and security preferences.
            </p>
          </div>
        </div>
      </div>

      <div className="grid gap-6">
        <ProfileNameForm />
        <ProfileEmailForm />
        <ProfilePasswordForm />
        <Security2FASetup />
        <AccessTokensList />
      </div>
    </div>
  );
}
