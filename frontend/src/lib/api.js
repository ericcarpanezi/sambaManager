import axios from "axios";
const baseURL = import.meta.env.VITE_API_BASE_URL ?? "/api/v1";
export const api = axios.create({
    baseURL,
    timeout: 10000,
});
api.interceptors.request.use((config) => {
    const raw = localStorage.getItem("agdm-auth");
    if (raw) {
        const parsed = JSON.parse(raw);
        if (parsed.accessToken) {
            config.headers.Authorization = `Bearer ${parsed.accessToken}`;
            config.headers["X-CSRF-Token"] = "agdm-static-csrf";
        }
    }
    return config;
});
