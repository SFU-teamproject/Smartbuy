import { Smartphone, ApiError, SignupData, AuthResponse, LoginData, User, CartItem, Cart } from '../types';

const API_BASE_URL = '/api/v1';

export async function apiClient<T>(
  endpoint: string,
  config?: RequestInit & { token?: string } // Добавляем опциональный token в конфиг
): Promise<T> {
  const headers = {
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    ...(config?.token && { 'Authorization': `Bearer ${config.token}` }), // Используем token из конфига
    ...config?.headers
  };

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    credentials: 'include',
    headers,
    ...config
  });

  if (!response.ok) {
    if (response.status === 401) {
      // Если токен невалидный, разлогиниваем
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    const errorData = await response.json().catch(() => null);
    throw new Error(errorData?.message || `HTTP error! status: ${response.status}`);
  }

  return response.json() as Promise<T>;
}


// Конкретные методы API

export const getSmartphones = (): Promise<Smartphone[]> => 
  apiClient<Smartphone[]>('/smartphones');

export const getSmartphoneById = (id: number, token?: string): Promise<Smartphone> => 
  apiClient<Smartphone>(`/smartphones/${id}`, token ? { token } : undefined);

export const getSmartphonesByIds = (ids: number[]): Promise<Smartphone[]> => {
  const idsString = ids.join(',');
  return apiClient<Smartphone[]>(`/smartphones?ids=${idsString}`);
};

export const signup = (data: SignupData): Promise<AuthResponse> => 
  apiClient<AuthResponse>('/signup', {
    method: 'POST',
    body: JSON.stringify(data)
  });

export const login = (data: LoginData): Promise<AuthResponse> => 
  apiClient<AuthResponse>('/login', {
    method: 'POST',
    body: JSON.stringify(data)
  });

// Пример обновленного метода с токеном
export const getUsers = (token: string): Promise<User[]> => 
  apiClient<User[]>('/users', { headers: { Authorization: `Bearer ${token}` }, });

export const getUserById = (id: number, token: string): Promise<User> => 
  apiClient<User>(`/users/${id}`, { headers: { Authorization: `${token}` }, });

export const getCartById = (cartId: number, token: string): Promise<Cart> =>
  apiClient<Cart>(`/carts/${cartId}`, { headers: { Authorization: `Bearer ${token}` }, });

// Получить корзину пользователя
export const getCartByUserId = (userId: number, token: string): Promise<Cart> =>
  apiClient<Cart>(`/carts?user_id=${userId}`, { headers: { Authorization: `Bearer ${token}` }, });

// Получить все корзины (для админа)
export const getAllCarts = (token: string): Promise<Cart[]> =>
  apiClient<Cart[]>('/carts', {
    headers: { Authorization: `Bearer ${token}` }
  });

// Получить items корзины
export const getCartItems = (cartId: number, token: string): Promise<CartItem[]> =>
  apiClient<CartItem[]>(`/carts/${cartId}/items`,  { headers: { Authorization: `Bearer ${token}` }, });

// Добавить item в корзину
export const addCartItem = (cartId: number, smartphoneId: number, token: string): Promise<CartItem> =>
  apiClient<CartItem>(`/carts/${cartId}/items`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify({ smartphone_id: smartphoneId })
  });

// Изменить количество item
export const updateCartItem = (cartId: number, itemId: number, quantity: number, token: string): Promise<CartItem> =>
  apiClient<CartItem>(`/carts/${cartId}/items/${itemId}`, {
    method: 'PATCH',
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify({ quantity })
  });

// Удалить item из корзины
export const deleteCartItem = (cartId: number, itemId: number, token: string): Promise<void> =>
  apiClient<void>(`/carts/${cartId}/items/${itemId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` }
  });