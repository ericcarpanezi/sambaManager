import { create } from "zustand";
const initialState = {
    accessToken: null,
    refreshToken: null,
    username: null,
    permissions: [],
};
export const useAuthStore = create((set) => ({
    ...initialState,
    setAuth: (data) => set((state) => {
        const next = { ...state, ...data };
        localStorage.setItem("agdm-auth", JSON.stringify(next));
        return next;
    }),
    clearAuth: () => {
        localStorage.removeItem("agdm-auth");
        set(initialState);
    },
}));
