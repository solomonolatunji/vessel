import { create } from "zustand";
import { settingsService } from "#/services/settings";

interface SystemState {
  registrationEnabled: boolean;
  siteName: string;
  emailEnabled: boolean;
  isLoaded: boolean;
  isLoading: boolean;
  fetchPublicSettings: () => Promise<void>;
}

export const useSystemStore = create<SystemState>((set) => ({
  registrationEnabled: false,
  siteName: "Vessl",
  emailEnabled: false,
  isLoaded: false,
  isLoading: false,

  fetchPublicSettings: async () => {
    set({ isLoading: true });
    try {
      const response = await settingsService.getPublicSettings();
      set({
        registrationEnabled: response.data.registrationEnabled,
        siteName: response.data.siteName || "Vessl",
        emailEnabled: response.data.emailEnabled,
        isLoaded: true,
      });
    } catch (error) {
      console.error("Failed to fetch public settings:", error);
      // Fallback defaults or just keep current state
      set({ isLoaded: true });
    } finally {
      set({ isLoading: false });
    }
  },
}));
