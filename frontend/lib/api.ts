import axios from "axios";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

/**
 * Pre-configured Axios instance that sends cookies with every request.
 * All API calls in this app must use this client.
 */
const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 15000,
});

// Response interceptor: surface error messages cleanly.
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError(error)) {
      const message =
        error.response?.data?.message ??
        error.response?.data?.error ??
        error.message ??
        "An unexpected error occurred";
      return Promise.reject(new Error(message));
    }
    return Promise.reject(error);
  }
);

export default api;
