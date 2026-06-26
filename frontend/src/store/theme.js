import { create } from "zustand";
export const useThemeStore = create((set) => ({
    theme: "dark",
    toggle: () => set((state) => {
        const nextTheme = state.theme === "dark" ? "light" : "dark";
        document.documentElement.classList.toggle("dark", nextTheme === "dark");
        return { theme: nextTheme };
    }),
}));
