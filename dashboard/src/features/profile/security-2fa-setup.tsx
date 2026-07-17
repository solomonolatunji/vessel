import { Loader2 } from 'lucide-react';
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
import { InputOTP, InputOTPGroup, InputOTPSlot } from '#/components/ui/input-otp';
import { useDisable2FA, useSetup2FA, useVerify2FA } from '#/hooks/useAuth';
import { useGetProfile } from '#/hooks/useProfile';

export function Security2FASetup() {
  const { data: profile, isLoading } = useGetProfile();

  const setup2FA = useSetup2FA();
  const verify2FA = useVerify2FA();
  const disable2FA = useDisable2FA();

  const [setupData, setSetupData] = useState<{ qrCodeUrl: string; secret: string } | null>(null);
  const [verifyOpen, setVerifyOpen] = useState(false);
  const [disableOpen, setDisableOpen] = useState(false);
  const [otp, setOtp] = useState('');

  const isEnabled = profile?.data?.totpEnabled;

  const handleEnableClick = () => {
    setup2FA.mutate(undefined, {
      onSuccess: (res) => {
        setSetupData(res.data);
        setOtp('');
        setVerifyOpen(true);
      },
      onError: (err) => toast.error(err.message),
    });
  };

  const handleDisableClick = () => {
    setOtp('');
    setDisableOpen(true);
  };

  const handleVerify = () => {
    verify2FA.mutate(
      { token: otp },
      {
        onSuccess: () => {
          setVerifyOpen(false);
          setSetupData(null);
          toast.success('Two-factor authentication enabled successfully.');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  const handleDisable = () => {
    disable2FA.mutate(
      { token: otp },
      {
        onSuccess: () => {
          setDisableOpen(false);
          toast.success('Two-factor authentication disabled.');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  if (isLoading) {
    return (
      <div className="flex justify-center rounded-2xl border border-border/50 bg-card/40 p-10">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <>
      <div className="space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
        <div className="flex items-center justify-between">
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
              <h2 className="font-semibold text-lg">Two-Factor Authentication</h2>
              <p className="text-muted-foreground text-sm">
                Add an extra layer of security to your account.
              </p>
            </div>
          </div>
          {isEnabled ? (
            <Button
              variant="destructive"
              onClick={handleDisableClick}
              className="h-11 font-bold text-xs uppercase tracking-wider"
            >
              DISABLE 2FA
            </Button>
          ) : (
            <Button
              onClick={handleEnableClick}
              disabled={setup2FA.isPending}
              className="h-11 bg-primary font-bold text-primary-foreground text-xs uppercase tracking-wider"
            >
              {setup2FA.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              ENABLE 2FA
            </Button>
          )}
        </div>

        <div className="space-y-4">
          <p className="text-muted-foreground text-sm leading-relaxed">
            {isEnabled
              ? 'Two-factor authentication is currently enabled. You will need to enter a code from your authenticator app when signing in.'
              : 'Protect your account from unauthorized access by requiring a second authentication method in addition to your password.'}
          </p>
        </div>
      </div>

      <Dialog open={verifyOpen} onOpenChange={setVerifyOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Setup Two-Factor Authentication</DialogTitle>
            <DialogDescription>
              Scan the QR code below with your authenticator app (like Google Authenticator or
              Authy), then enter the 6-digit code.
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col items-center space-y-6 py-4">
            {setupData && (
              <div className="rounded-lg bg-white p-4">
                <img src={setupData.qrCodeUrl} alt="2FA QR Code" className="h-48 w-48" />
              </div>
            )}
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
            <Button variant="outline" onClick={() => setVerifyOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleVerify} disabled={otp.length !== 6 || verify2FA.isPending}>
              {verify2FA.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Verify & Enable
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={disableOpen} onOpenChange={setDisableOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Disable Two-Factor Authentication</DialogTitle>
            <DialogDescription>
              Enter a 6-digit code from your authenticator app to confirm you want to disable 2FA.
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
            <Button variant="outline" onClick={() => setDisableOpen(false)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDisable}
              disabled={otp.length !== 6 || disable2FA.isPending}
            >
              {disable2FA.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Disable
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
