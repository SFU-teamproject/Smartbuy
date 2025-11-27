// src/types.ts
export interface Smartphone {
  id: number;
  model: string;
  producer: string;
  memory: number;
  ram: number;
  display_size: number;
  price: number;
  ratings_sum: number;
  ratings_count: number;
  image_path?: string;
  description?: string;
}

export interface ApiError {
  message: string;
  status?: number;
}

export interface Cart {
  id: number;
  user_id?: number;
  created_at: string;
  updated_at: string;
  items?: CartItem[];
}

export interface CartItem {
  id: number;
  smartphone_id: number;
  quantity: number;
  smartphone?: Smartphone; // Опционально, если бэкенд возвращает полные данные
}

export interface User {
  id: number;
  name: string;
  role: 'admin' | 'user';
  created_at: string;
  cart?: Cart;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface LoginData {
  name: string;
  password: string;
}

export interface SignupData extends LoginData {}

export interface Order {
  id: number;
  user_id: number;
  status: 'pending' | 'processing' | 'shipped' | 'delivered' | 'cancelled';
  total_amount: number;
  created_at: string;
  updated_at: string;
  items: OrderItem[];
}

export interface OrderItem {
  id: number;
  order_id: number;
  smartphone_id: number;
  quantity: number;
  price: number;
  smartphone?: Smartphone;
}

export interface CreateOrderData {
  user_id: number;
  items: {
    smartphone_id: number;
    quantity: number;
    price: number;
  }[];
}