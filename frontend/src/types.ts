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
  // Для login бэкенд ожидает email и password
  email: string;      // Добавляем email
  password: string;
  name?: string;      // Делаем опциональным
}

export interface SignupData {
  // Для signup бэкенд ожидает только email
  email: string;
}

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

/* merge */

export interface Review {
  id: number;
  smartphone_id: number;
  user_id: number;
  user_name?: string; //Поле user_name присутствует только при GET всех отзывов по айди смартфона, в остальных отсутствует.
  rating: number;
  comment?: string; //Поле comment может отсутствовать
  created_at: Date;
  updated_at?: Date;
}

export interface ReviewForAdd {
  rating: number;
  comment?: string; //Поле comment может отсутствовать
}

export interface ReviewForUpdate {
  id: number;
  rating: number;
  comment?: string; //Поле comment может отсутствовать
}

export interface Payment {
  phone: string;
  delivery_type: DeliveryType;
  address?: Address;
  payment_type?: PaymentType;
}

export interface Address {
  city: string;
  street: string;
  house: number;
}

export enum DeliveryType {
  delivery_self = "delivery_self",
  delivery_courier = "delivery_courier"
}

export enum PaymentType {
  payment_cart = "payment_cart",
  payment_courier = "payment_courier",
  payment_self = "payment_self"
}