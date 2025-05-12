// src/api/client.ts
import { Smartphone, ApiError } from '@/types';

//const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081/api/v1';
//const API_BASE_URL = 'http://localhost:8081/api/v1'; //полный URL
const API_BASE_URL = '/api/v1'; // Теперь будет проксироваться

export async function apiClient<T>(
  endpoint: string,
  config?: RequestInit
): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    credentials: 'include',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
    ...config,
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => null);
    throw new Error(errorData?.message || `HTTP error! status: ${response.status}`);
  }

  return response.json() as Promise<T>;
}
/*
export async function apiClient<T>(endpoint: string): Promise<T> {
  const response = await fetch(`http://localhost:8081/api/v1${endpoint}`);
  
  // Принудительная проверка типа
  const contentType = response.headers.get('content-type');
  if (!contentType?.includes('application/json')) {
    const text = await response.text();
    throw new Error(`Expected JSON, got ${contentType}: ${text.substring(0, 100)}`);
  }

  return response.json();
}
  */
// Конкретные методы API
export const getSmartphones = (): Promise<Smartphone[]> => 
  apiClient<Smartphone[]>('/smartphones');

export const getSmartphoneById = (id: number): Promise<Smartphone> => 
  apiClient<Smartphone>(`/smartphones/${id}`);

/*export const login = (credentials: {
  email: string;
  password: string;
}): Promise<{ token: string }> =>
  apiClient('/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  });*/