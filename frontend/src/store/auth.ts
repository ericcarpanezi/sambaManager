import { create } from "zustand";

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  username: string | null;
  permissions: string[];
  setAuth: (data: Partial<AuthState>) => void;
  clearAuth: () => void;
}

const initialState = {
  accessToken: null,
  refreshToken: null,
  username: null,
  permissions: [],
};

export const useAuthStore = create<AuthState>((set) => ({
  ...initialState,
  setAuth: (data) =>
    set((state) => {
      const next = { ...state, ...data };
      localStorage.setItem("agdm-auth", JSON.stringify(next));
      return next;
    }),
  clearAuth: () => {
    localStorage.removeItem("agdm-auth");
    set(initialState);
  },
}));
