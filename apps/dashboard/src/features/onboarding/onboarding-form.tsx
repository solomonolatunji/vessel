import { zodResolver } from '@hookform/resolvers/zod';
import {
  ArrowLeft,
  ArrowRight,
  Check,
  Database,
  Globe,
  Server,
  Settings,
  Shield,
} from 'lucide-react';
import { FormProvider, useForm } from 'react-hook-form';
import { Button } from '#/components/ui/button';
import { useSetup } from '#/hooks/useAuth';
import { useOnboardingStore } from '#/stores/onboardingStore';
import { ImportModal, type SetupSchema, StepDomain, StepOwner, StepRuntime, setupSchema } from '.';

export const OnboardingForm = ({ cwd }: { cwd?: string }) => {
  const { mutateAsync: setupUser, isPending } = useSetup();
  const currentStep = useOnboardingStore((state) => state.currentStep);
  const isImportModalOpen = useOnboardingStore((state) => state.isImportModalOpen);

  const methods = useForm<SetupSchema>({
    resolver: zodResolver(setupSchema),
    defaultValues: {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
      env: {
        jwtSecret: '',
        dataDir: cwd ? `${cwd}/data` : './data',
        dashboardUrl: 'http://localhost:3000',
        port: 8080,
      },
      dashboardDomain: '',
      defaultWildcardDomain: '',
    },
  });

  const { handleSubmit, trigger } = methods;

  const nextStep = async () => {
    let fieldsToValidate: (keyof SetupSchema)[] = [];
    switch (currentStep) {
      case 1:
        fieldsToValidate = ['name', 'email', 'password', 'confirmPassword'];
        break;
      case 2:
        fieldsToValidate = ['env'];
        break;
      case 3:
        fieldsToValidate = ['dashboardDomain', 'defaultWildcardDomain'];
        break;
    }

    const isValid = await trigger(fieldsToValidate);
    if (isValid) {
      if (typeof document !== 'undefined' && document.activeElement instanceof HTMLElement) {
        document.activeElement.blur();
      }
      useOnboardingStore.getState().nextStep();
    }
  };

  const prevStep = () => useOnboardingStore.getState().prevStep();

  const onSubmit = async (data: SetupSchema) => {
    if (currentStep !== 3) {
      return nextStep();
    }
    if (typeof document !== 'undefined' && document.activeElement instanceof HTMLElement) {
      document.activeElement.blur();
    }
    try {
      await setupUser(data);
    } catch (error) {
      console.error('Setup failed:', error);
    }
  };

  const steps = [
    {
      num: 1,
      label: 'Owner',
      title: 'Create owner access',
      description:
        'This first account controls the instance, deployment actions, stored secrets, backups, and future user management.',
      icon: Shield,
    },
    {
      num: 2,
      label: 'Runtime',
      title: 'Runtime environment',
      description: 'Configure the deployment runtime settings. These are written to .env.local.',
      icon: Settings,
    },
    {
      num: 3,
      label: 'Root Domain',
      title: 'Root domain',
      description: 'Configure the default domains for your applications.',
      icon: Globe,
    },
  ];

  return (
    <div className="mx-auto w-full max-w-5xl pt-8 pb-6">
      <div className="mb-6 flex items-center justify-between border-border/50 border-b pb-6">
        <div className="flex items-center gap-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <Server className="h-6 w-6" />
          </div>
          <div>
            <p className="mb-1 font-semibold text-muted-foreground text-xs uppercase tracking-wider">
              FIRST RUN
            </p>
            <h1 className="font-bold text-3xl tracking-tight">Set up Codedock</h1>
          </div>
        </div>
        <div className="font-medium text-muted-foreground text-sm uppercase tracking-widest">
          STEP {currentStep} OF 3
        </div>
      </div>

      <div className="relative mb-8 flex justify-between">
        <div className="absolute top-1/2 left-0 h-px w-full -translate-y-1/2 bg-border/50" />
        <div
          className="absolute top-1/2 left-0 h-px -translate-y-1/2 bg-primary transition-all duration-300"
          style={{ width: `${((currentStep - 1) / 2) * 100}%` }}
        />
        {steps.map((step) => {
          const isComplete = step.num < currentStep;
          const isCurrent = step.num === currentStep;
          const StepIcon = step.icon;

          return (
            <div
              key={step.num}
              className="relative z-10 flex items-center gap-3 bg-background px-4 py-1"
            >
              <div
                className={`flex h-8 w-8 items-center justify-center rounded-full border-2 transition-colors ${
                  isComplete
                    ? 'border-primary bg-primary text-primary-foreground'
                    : isCurrent
                      ? 'border-primary text-primary'
                      : 'border-muted-foreground text-muted-foreground'
                }`}
              >
                {isComplete ? <Check className="h-4 w-4" /> : <StepIcon className="h-4 w-4" />}
              </div>
              <div className="flex flex-col">
                <span className="font-semibold text-[10px] text-muted-foreground">0{step.num}</span>
                <span
                  className={`font-bold text-xs uppercase tracking-widest ${isCurrent ? 'text-foreground' : 'text-muted-foreground'}`}
                >
                  {step.label}
                </span>
              </div>
            </div>
          );
        })}
      </div>

      <FormProvider {...methods}>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="mb-6 min-h-87.5 rounded-xl border border-border/50 bg-card/40 p-6 shadow-xl backdrop-blur-xl">
            <div className="mb-6">
              <p className="mb-3 font-bold text-primary text-xs uppercase tracking-widest">
                STEP 0{currentStep}
              </p>
              <h2 className="mb-4 font-bold text-2xl">{steps[currentStep - 1].title}</h2>
              <p className="max-w-2xl text-muted-foreground text-sm leading-relaxed">
                {steps[currentStep - 1].description}
              </p>
            </div>
            {currentStep === 1 && <StepOwner />}
            {currentStep === 2 && <StepRuntime />}
            {currentStep === 3 && <StepDomain />}
          </div>

          <div className="flex items-center justify-between pb-6">
            <Button
              type="button"
              variant="outline"
              onClick={prevStep}
              disabled={currentStep === 1}
              className="flex h-11 items-center gap-2 rounded-xl px-6 font-semibold text-muted-foreground text-xs uppercase tracking-widest hover:text-foreground"
            >
              <ArrowLeft className="h-4 w-4" /> BACK
            </Button>

            {currentStep < 3 ? (
              <Button
                type="button"
                onClick={nextStep}
                className="flex h-11 items-center gap-2 rounded-xl border-primary/20 bg-primary/10 px-6 font-semibold text-primary text-xs uppercase tracking-widest hover:bg-primary/20 hover:text-primary"
              >
                CONTINUE <ArrowRight className="h-4 w-4" />
              </Button>
            ) : (
              <Button
                type="submit"
                disabled={isPending}
                className="flex h-11 items-center gap-2 rounded-xl bg-primary px-6 font-semibold text-primary-foreground text-xs uppercase tracking-widest hover:bg-primary/90"
              >
                {isPending ? 'SAVING...' : 'COMPLETE SETUP'}
                {!isPending && <Check className="h-4 w-4" />}
              </Button>
            )}
          </div>
        </form>
      </FormProvider>

      <div className="mt-16 flex justify-center border-border/50 border-t pt-10">
        <Button
          variant="outline"
          onClick={() => useOnboardingStore.getState().setImportModalOpen(true)}
          className="flex h-11 items-center gap-2 rounded-xl bg-background px-6 font-semibold text-muted-foreground text-xs uppercase tracking-widest transition-all duration-300 hover:border-primary/50 hover:text-foreground"
        >
          <Database className="h-4 w-4" />
          IMPORT EXISTING CODEDOCK
        </Button>
        <ImportModal
          open={isImportModalOpen}
          onOpenChange={useOnboardingStore.getState().setImportModalOpen}
        />
      </div>
    </div>
  );
};
