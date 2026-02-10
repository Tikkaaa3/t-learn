export const API_BASE = "http://localhost:8080";

// A wrapper around the native fetch API
export async function apiClient<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  // Get the token from localStorage
  const token = localStorage.getItem("t_learn_token");

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>), // Cast incoming headers
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  // Make the Request
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  // Handle Errors
  if (!response.ok) {
    // Try to parse the error message from the backend JSON
    let errorMessage = `Error ${response.status}: ${response.statusText}`;
    try {
      const errorData = await response.json();
      if (errorData.error) errorMessage = errorData.error;
    } catch {
      // If response isn't JSON, just use the status text
    }
    throw new Error(errorMessage);
  }

  // Return Data (if any)
  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}
