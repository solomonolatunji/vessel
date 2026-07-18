import { Loader2, ShieldCheck, Trash2 } from 'lucide-react';
import QRCode from 'qrcode';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '#/components/ui/dialog';
import { InputOTP, InputOTPGroup, InputOTPSlot } from '#/components/ui/input-otp';
import { Section } from '#/components/ui/section';
import { useDisable2FA, useSetup2FA, useVerify2FA } from '#/hooks/useAuth';
import { useGetProfile } from '#/hooks/useProfile';

function OtpSlots() {
  const indices = [0, 1, 2, 3, 4, 5] as const;
  return (
    <InputOTPGroup className="gap-2">
      {indices.map((i) => (
        <InputOTPSlot
          key={i}
          index={i}
          className="size-12 rounded-xl border border-border/50 bg-background/80 text-base first:rounded-xl first:border last:rounded-xl last:border"
        />
      ))}
    </InputOTPGroup>
  );
}

function QrCodeImage({ uri }: { uri: string }) {
  const [dataUrl, setDataUrl] = useState('');

  useEffect(() => {
    if (!uri) return;
    QRCode.toDataURL(uri, {
      width: 200,
      margin: 1,
      color: { dark: '#000000', light: '#ffffff' },
    })
      .then(setDataUrl)
      .catch(() => setDataUrl(''));
  }, [uri]);

  if (!dataUrl) {
    return (
      <div className="flex h-50 w-50 items-center justify-center rounded-xl border border-border/50 bg-muted/30">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="rounded-xl bg-white p-4 shadow-sm">
      <img src={dataUrl} alt="2FA QR Code" className="h-50 w-50" />
    </div>
  );
}

export function Security2FASetup() {
  const { data: profile, isLoading } = useGetProfile();

  const setup2FA = useSetup2FA();
  const verify2FA = useVerify2FA();
  const disable2FA = useDisable2FA();

  const [qrCodeUri, setQrCodeUri] = useState('');
  const [verifyOpen, setVerifyOpen] = useState(false);
  const [disableOpen, setDisableOpen] = useState(false);
  const [otp, setOtp] = useState('');

  const isEnabled = profile?.data?.totpEnabled;

  const handleEnableClick = () => {
    setup2FA.mutate(undefined, {
      onSuccess: (res) => {
        setQrCodeUri(res.data?.qrCodeUri ?? '');
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
      { passcode: otp },
      {
        onSuccess: () => {
          setVerifyOpen(false);
          setQrCodeUri('');
          toast.success('Two-factor authentication enabled successfully.');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  const handleDisable = () => {
    disable2FA.mutate(
      { passcode: otp },
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
      <Section
        icon={<ShieldCheck className="h-4 w-4" />}
        title="Two-Factor Authentication"
        action={
          isEnabled ? (
            <Button variant="destructive" size="sm" onClick={handleDisableClick}>
              Disable 2FA
            </Button>
          ) : (
            <Button size="sm" onClick={handleEnableClick} disabled={setup2FA.isPending}>
              {setup2FA.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Enable 2FA
            </Button>
          )
        }
      >
        <div className="py-4">
          <p className="text-muted-foreground text-sm leading-relaxed">
            {isEnabled
              ? 'Two-factor authentication is currently enabled. You will need to enter a code from your authenticator app when signing in.'
              : 'Protect your account from unauthorized access by requiring a second authentication method in addition to your password.'}
          </p>
        </div>
      </Section>

      <Dialog open={verifyOpen} onOpenChange={setVerifyOpen}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-[400px] [&>button]:hidden">
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                  <ShieldCheck className="h-5 w-5 text-primary" />
                  Setup Two-Factor Auth
                </DialogTitle>
                <DialogDescription className="mt-1.5 flex items-center gap-1.5 font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                  <ShieldCheck className="h-3 w-3" />
                  Scan QR code & enter code
                </DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  type="button"
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="h-px w-full bg-border/50" />

          <div className="flex flex-col items-center gap-5 p-5">
            <QrCodeImage uri={qrCodeUri} />

            <div className="flex flex-col items-center gap-2">
              <p className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                ENTER VERIFICATION CODE
              </p>
              <InputOTP maxLength={6} value={otp} onChange={setOtp}>
                <OtpSlots />
              </InputOTP>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              variant="ghost"
              onClick={() => setVerifyOpen(false)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              onClick={handleVerify}
              disabled={otp.length !== 6 || verify2FA.isPending}
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              {verify2FA.isPending ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <ShieldCheck className="h-3.5 w-3.5" />
              )}
              Verify & Enable
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={disableOpen} onOpenChange={setDisableOpen}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-[400px] [&>button]:hidden">
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-destructive text-xl tracking-tight">
                  <Trash2 className="h-5 w-5" />
                  Disable Two-Factor Auth
                </DialogTitle>
                <DialogDescription className="mt-1.5 flex items-center gap-1.5 font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                  <Trash2 className="h-3 w-3" />
                  Confirm with authenticator app
                </DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  type="button"
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
            <p className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
              ENTER VERIFICATION CODE
            </p>
            <InputOTP maxLength={6} value={otp} onChange={setOtp}>
              <OtpSlots />
            </InputOTP>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              variant="ghost"
              onClick={() => setDisableOpen(false)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDisable}
              disabled={otp.length !== 6 || disable2FA.isPending}
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              {disable2FA.isPending ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <Trash2 className="h-3.5 w-3.5" />
              )}
              Disable 2FA
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
