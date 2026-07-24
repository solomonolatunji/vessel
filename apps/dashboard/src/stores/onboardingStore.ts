import { create } from 'zustand';

export interface OnboardingState {
  currentStep: number;
  isImportModalOpen: boolean;
  setStep: (step: number) => void;
  nextStep: () => void;
  prevStep: () => void;
  setImportModalOpen: (isOpen: boolean) => void;
}

export const useOnboardingStore = create<OnboardingState>((set) => ({
  currentStep: 1,
  isImportModalOpen: false,
  setStep: (step: number) => {
    set(() => ({
      currentStep: Math.min(Math.max(step, 1), 3),
    }));
  },
  nextStep: () => {
    set((state) => ({
      currentStep: Math.min(state.currentStep + 1, 3),
    }));
  },
  prevStep: () => {
    set((state) => ({
      currentStep: Math.max(state.currentStep - 1, 1),
    }));
  },
  setImportModalOpen: (isOpen: boolean) => {
    set({ isImportModalOpen: isOpen });
  },
}));
