import {
    Smartphone,
    ApiError,
    SignupData,
    AuthResponse,
    LoginData,
    User,
    CartItem,
    Cart,
    Order,
    CreateOrderData,
    Review, ReviewForAdd, ReviewForUpdate
} from '../types';
import { mockOrders } from './mockOrders'; // Импортируем моковые данные для странички заказов

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
    apiClient<CartItem[]>(`/carts/${cartId}/items`, { headers: { Authorization: `Bearer ${token}` }, });

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

// api методы для заказов
export const createOrder = (data: CreateOrderData, token: string): Promise<Order> =>
/* apiClient<Order>('/orders', {
   method: 'POST',
   headers: { Authorization: `Bearer ${token}` },
   body: JSON.stringify(data)
 });*/
// Для демонстрации создаем моковый заказ
{
    const newOrder: Order = {
        id: Math.max(...mockOrders.map(o => o.id)) + 1,
        user_id: data.user_id,
        status: 'pending',
        total_amount: data.items.reduce((sum, item) => sum + (item.price * item.quantity), 0),
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        items: data.items.map((item, index) => ({
            id: Math.max(...mockOrders.flatMap(o => o.items.map(i => i.id))) + index + 1,
            order_id: Math.max(...mockOrders.map(o => o.id)) + 1,
            smartphone_id: item.smartphone_id,
            quantity: item.quantity,
            price: item.price,
            // Здесь можно добавить smartphone данные, если нужно
        }))
    };

    // В реальном приложении здесь был бы API вызов:
    // return apiClient<Order>('/orders', {
    //   method: 'POST',
    //   headers: { Authorization: `Bearer ${token}` },
    //   body: JSON.stringify(data)
    // });

    // Для демо возвращаем моковый заказ
    return new Promise((resolve) => {
        setTimeout(() => resolve(newOrder), 500); // Имитация задержки сети
    });
};


export const getOrders = (token: string): Promise<Order[]> =>
 /* apiClient<Order[]>('/orders', {
    headers: { Authorization: `Bearer ${token}` }
  });*/ {
    // В реальном приложении:
    // return apiClient<Order[]>('/orders', {
    //   headers: { Authorization: `Bearer ${token}` }
    // });

    // Для демо возвращаем моковые заказы
    return new Promise((resolve) => {
        setTimeout(() => resolve(mockOrders), 500);
    });
};


export const getOrderById = (orderId: number, token: string): Promise<Order> =>
 /* apiClient<Order>(`/orders/${orderId}`, {
    headers: { Authorization: `Bearer ${token}` }
  });*/ {
    // В реальном приложении:
    // return apiClient<Order>(`/orders/${orderId}`, {
    //   headers: { Authorization: `Bearer ${token}` }
    // });

    // Для демо находим заказ в моковых данных
    const order = mockOrders.find(o => o.id === orderId);
    if (!order) {
        throw new Error('Order not found');
    }

    return new Promise((resolve) => {
        setTimeout(() => resolve(order), 300);
    });
};

export const cancelOrder = (orderId: number, token: string): Promise<void> =>
  /*apiClient<void>(`/orders/${orderId}/cancel`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` }
  });*/ {
    // В реальном приложении:
    // return apiClient<void>(`/orders/${orderId}/cancel`, {
    //   method: 'POST',
    //   headers: { Authorization: `Bearer ${token}` }
    // });

    // Для демо просто возвращаем успех
    return new Promise((resolve) => {
        setTimeout(() => resolve(), 500);
    });
};

// Отзывы
// Все отзывы к смартфону
export const getReviews = (smartphoneId: number): Promise<Review[]> =>
    apiClient<Review[]>(`/smartphones/${smartphoneId}/reviews`);

// Добавить отзыв к смарфтону
export const addReview = (smartphoneId: number, token: string, review: ReviewForAdd): Promise<Review> =>
    apiClient<Review>(`/smartphones/${smartphoneId}/reviews`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify(review) // добавить условие и комментарий
    });

// Обновить отзыв к смарфону
export const updateReview = (smartphoneId: number, token: string, review: ReviewForUpdate): Promise<Review> =>
    apiClient<Review>(`/smartphones/${smartphoneId}/reviews/${review.id}`, {
        method: 'PATCH',
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify(review) // добавить условие и комментарий
    });

// Удалить отзыв к смартфону
export const deleteReview = (smartphoneId: number, reviewId: number, token: string): Promise<void> =>
    apiClient<void>(`/smartphones/${smartphoneId}/reviews/${reviewId}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
    });