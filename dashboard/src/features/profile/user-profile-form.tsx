import { Check, Loader2, Mail, User } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { InputOTP, InputOTPGroup, InputOTPSlot } from '#/components/ui/input-otp';
import { Row, Section } from '#/components/ui/section';
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
    <form onSubmit={handleSave}>
      <Section
        icon={<User className="h-4 w-4" />}
        title="Profile Name"
        action={
          <Button
            type="submit"
            size="sm"
            disabled={isLoading || updateProfile.isPending || name === profile?.data?.name}
          >
            {updateProfile.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Check className="mr-2 h-4 w-4" />
            )}
            Save Changes
          </Button>
        }
      >
        <Row label="Full Name" description="Update your display name.">
          <Input
            id="name"
            placeholder="John Doe"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="h-10 w-full"
          />
        </Row>
      </Section>
    </form>
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
      <form onSubmit={handleRequest}>
        <Section
          icon={<Mail className="h-4 w-4" />}
          title="Email Address"
          action={
            <Button
              type="submit"
              size="sm"
              disabled={isLoading || requestEmailChange.isPending || email === profile?.data?.email}
            >
              {requestEmailChange.isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Check className="mr-2 h-4 w-4" />
              )}
              Save Changes
            </Button>
          }
        >
          <Row
            label="New Email Address"
            description="Change the email address associated with your account."
          >
            <Input
              id="email"
              type="email"
              placeholder="john@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="h-10 w-full"
            />
          </Row>
        </Section>
      </form>

      <Dialog open={otpOpen} onOpenChange={setOtpOpen}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-100 [&>button]:hidden">
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                  <Mail className="h-5 w-5 text-primary" />
                  Verify Email Change
                </DialogTitle>
                <DialogDescription>Code sent to new email</DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>
          <div className="h-px w-full bg-border/50" />
          <div className="flex flex-col items-center gap-4 p-5">
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
          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              variant="ghost"
              onClick={() => setOtpOpen(false)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              onClick={handleVerify}
              disabled={otp.length !== 6 || verifyEmailChange.isPending}
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              {verifyEmailChange.isPending ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <Check className="h-3.5 w-3.5" />
              )}
              Verify & Update
            </Button>
          </div>
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
    <form onSubmit={handleSave}>
      <Section
        icon={
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
            className="h-4 w-4"
          >
            <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
            <path d="M7 11V7a5 5 0 0 1 10 0v4" />
          </svg>
        }
        title="Change Password"
        action={
          <Button
            type="submit"
            size="sm"
            disabled={changePassword.isPending || !oldPassword || !newPassword}
          >
            {changePassword.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Check className="mr-2 h-4 w-4" />
            )}
            Save Changes
          </Button>
        }
      >
        <Row label="Current Password" description="Update the password you use to sign in.">
          <Input
            id="oldPassword"
            type="password"
            value={oldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
            className="h-10 w-full"
            placeholder="Current Password"
          />
        </Row>
        <Row label="New Password" description="Must be at least 8 characters long.">
          <Input
            id="newPassword"
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="h-10 w-full"
            placeholder="New Password"
          />
        </Row>
      </Section>
    </form>
  );
}
