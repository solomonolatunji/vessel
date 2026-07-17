import { Loader2, Mail, User } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { InputOTP, InputOTPGroup, InputOTPSlot } from '#/components/ui/input-otp';
import { Label } from '#/components/ui/label';
import {
  useChangePassword,
  useGetProfile,
  useRequestEmailChange,
  useUpdateProfile,
  useVerifyEmailChange,
} from '#/hooks/useProfile';

export function ProfileNameForm() {
  const { data: profile, isLoading } = useGetProfile();
  const updateProfile = useUpdateProfile();
  const [name, setName] = useState('');

  // Sync state once profile loads
  if (!isLoading && profile?.data && !name && profile.data.name) {
    setName(profile.data.name);
  }

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfile.mutate(
      { name },
      {
        onSuccess: () => toast.success('Profile name updated!'),
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <div className="space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
      <form onSubmit={handleSave} className="space-y-6">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
            <User className="h-5 w-5" />
          </div>
          <div>
            <h2 className="font-semibold text-lg">Profile Name</h2>
            <p className="text-muted-foreground text-sm">Update your display name.</p>
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="name" className="font-medium text-sm">
            Full Name
          </Label>
          <Input
            id="name"
            placeholder="John Doe"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="h-9 max-w-md text-sm"
          />
        </div>

        <div className="flex justify-end border-border/50 border-t pt-4">
          <Button
            type="submit"
            size="sm"
            disabled={isLoading || updateProfile.isPending || name === profile?.data?.name}
          >
            {updateProfile.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Save Name
          </Button>
        </div>
      </form>
    </div>
  );
}

export function ProfileEmailForm() {
  const { data: profile, isLoading } = useGetProfile();
  const requestEmailChange = useRequestEmailChange();
  const verifyEmailChange = useVerifyEmailChange();

  const [email, setEmail] = useState('');
  const [otpOpen, setOtpOpen] = useState(false);
  const [otp, setOtp] = useState('');

  // Sync state once profile loads
  if (!isLoading && profile?.data && !email && profile.data.email) {
    setEmail(profile.data.email);
  }

  const handleRequest = (e: React.FormEvent) => {
    e.preventDefault();
    requestEmailChange.mutate(
      { newEmail: email },
      {
        onSuccess: () => {
          setOtpOpen(true);
          toast.success('Verification code sent to your new email.');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  const handleVerify = (e: React.FormEvent) => {
    e.preventDefault();
    verifyEmailChange.mutate(
      { otp },
      {
        onSuccess: () => {
          setOtpOpen(false);
          setOtp('');
          toast.success('Email updated successfully!');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <>
      <div className="space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
        <form onSubmit={handleRequest} className="space-y-6">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
              <Mail className="h-5 w-5" />
            </div>
            <div>
              <h2 className="font-semibold text-lg">Email Address</h2>
              <p className="text-muted-foreground text-sm">
                Change the email address associated with your account.
              </p>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="email" className="font-medium text-sm">
              New Email Address
            </Label>
            <Input
              id="email"
              type="email"
              placeholder="john@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="h-9 max-w-md text-sm"
            />
          </div>

          <div className="flex justify-end border-border/50 border-t pt-4">
            <Button
              type="submit"
              size="sm"
              disabled={isLoading || requestEmailChange.isPending || email === profile?.data?.email}
            >
              {requestEmailChange.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Update Email
            </Button>
          </div>
        </form>
      </div>

      <Dialog open={otpOpen} onOpenChange={setOtpOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Verify Email Change</DialogTitle>
            <DialogDescription>
              We've sent a 6-digit verification code to your new email. Enter it below.
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-center py-6">
            <InputOTP maxLength={6} value={otp} onChange={setOtp}>
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOtpOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleVerify}
              disabled={otp.length !== 6 || verifyEmailChange.isPending}
            >
              {verifyEmailChange.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Verify & Update
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

export function ProfilePasswordForm() {
  const changePassword = useChangePassword();
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    changePassword.mutate(
      { oldPassword, newPassword },
      {
        onSuccess: () => {
          toast.success('Password updated successfully!');
          setOldPassword('');
          setNewPassword('');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <div className="space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
      <form onSubmit={handleSave} className="space-y-6">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="h-5 w-5"
            >
              <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
              <path d="M7 11V7a5 5 0 0 1 10 0v4" />
            </svg>
          </div>
          <div>
            <h2 className="font-semibold text-lg">Change Password</h2>
            <p className="text-muted-foreground text-sm">Update the password you use to sign in.</p>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="oldPassword" className="font-medium text-sm">
              Current Password
            </Label>
            <Input
              id="oldPassword"
              type="password"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              className="h-9 text-sm"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="newPassword" className="font-medium text-sm">
              New Password
            </Label>
            <Input
              id="newPassword"
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              className="h-9 text-sm"
            />
          </div>
        </div>

        <div className="flex justify-end border-border/50 border-t pt-4">
          <Button
            type="submit"
            size="sm"
            disabled={changePassword.isPending || !oldPassword || !newPassword}
          >
            {changePassword.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Update Password
          </Button>
        </div>
      </form>
    </div>
  );
}
