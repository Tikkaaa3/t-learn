import { apiClient } from "./client";

// Define what the backend returns when we log in
export interface LoginResponse {
  token: string;
  user: {
    id: string;
    username: string;
  };
}

export interface ApiKeyResponse {
  api_key: string;
}

export interface RegisterResponse {
  id: string;
  username: string;
  email: string;
}

// The Login Function
export async function loginUser(
  username: string,
  pass: string,
): Promise<LoginResponse> {
  return apiClient<LoginResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ username, password: pass }),
  });
}

export async function generateApiKey(): Promise<ApiKeyResponse> {
  // This endpoint is protected, but apiClient automatically attaches your
  // login token from localStorage, so it will pass the MiddlewareAuth check.
  return apiClient<ApiKeyResponse>("/auth/token", {
    method: "POST",
  });
}

export async function registerUser(
  username: string,
  email: string,
  pass: string,
): Promise<RegisterResponse> {
  return apiClient<RegisterResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify({ username, email, password: pass }),
  });
}
