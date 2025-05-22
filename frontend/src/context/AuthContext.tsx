// src/context/AuthContext.tsx
/*import { createContext, useContext, useState, ReactNode } from 'react';
import { AuthResponse, User, Cart } from '../types';
import { getCartItems, getUserById } from '../api/client';
import { useEffect } from 'react';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (data: AuthResponse) => void;
  logout: () => void;
  isAuthenticated: boolean;
  isAdmin: boolean;
  updateCart: (cart: Cart) => void;
  cart: Cart | null;
  refreshCart: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
  const [cart, setCart] = useState<Cart | null>(null);

  const login = (data: AuthResponse) => {
    localStorage.setItem('token', data.token);
    setUser(data.user);
    setToken(data.token);
    setCart(data.user.cart || null);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
    setToken(null);
    setCart(null);
  };

  const updateCart = (cart: Cart) => {
    setUser(prev => prev ? {...prev, cart} : null);
    setCart(cart);
  };

  const refreshCart = async () => {
    if (!user?.cart?.id || !token) return;
    try {
      const data = await getCartItems(user.cart.id, token);
      const updatedCart = { ...user.cart, items: data };
      updateCart(updatedCart);
    } catch (error) {
      console.error('Failed to refresh cart:', error);
    }
  };

  const isAuthenticated = !!token;
  const isAdmin = user?.role === 'admin';

  useEffect(() => {
    const storedToken = localStorage.getItem('token');
  if (storedToken) {
    // Дополнительно можно проверить валидность токена
    setToken(storedToken);
    // Можно добавить запрос для получения данных пользователя
  }
}, []);

  return (
    <AuthContext.Provider 
      value={{
        user,
        token,
        login,
        logout,
        isAuthenticated,
        isAdmin,
        updateCart,
        cart,
        refreshCart
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
*/// src/context/AuthContext.tsx
import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { AuthResponse, User, Cart } from '../types';
import { getCartItems, getUserById, getCartByUserId } from '../api/client';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (data: AuthResponse) => void;
  logout: () => void;
  isAuthenticated: boolean;
  isAdmin: boolean;
  updateCart: (cart: Cart) => void;
  cart: Cart | null;
  refreshCart: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
  const [cart, setCart] = useState<Cart | null>(null);

  // Инициализация пользователя при загрузке приложения
  useEffect(() => {
    const initializeAuth = async () => {
      const storedToken = localStorage.getItem('token');
      if (storedToken) {
        try {
          // Получаем данные пользователя по токену
          const userData = await getUserById(1, storedToken); // Нужно реализовать endpoint /me на бэкенде
          const userCart = await getCartByUserId(userData.id, storedToken);
          
          setToken(storedToken);
          setUser(userData);
          setCart(userCart);
        } catch (error) {
          console.error('Auth initialization error:', error);
          localStorage.removeItem('token');
        }
      }
    };

    initializeAuth();
  }, []);

  const login = async (data: AuthResponse) => {
    localStorage.setItem('token', data.token);
    setToken(data.token);
    setUser(data.user);
    setCart(data.user.cart || null);
    
    // Получаем актуальные данные корзины
    if (data.user.cart?.id) {
      try {
        const cartItems = await getCartItems(data.user.cart.id, data.token);
        setCart({ ...data.user.cart, items: cartItems });
      } catch (error) {
        console.error('Failed to load cart:', error);
      }
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
    setToken(null);
    setCart(null);
  };

  const updateCart = (cart: Cart) => {
    setUser(prev => prev ? { ...prev, cart } : null);
    setCart(cart);
  };

  const refreshCart = async () => {
    if (!user?.cart?.id || !token) return;
    try {
      const data = await getCartItems(user.cart.id, token);
      const updatedCart = { ...user.cart, items: data };
      updateCart(updatedCart);
    } catch (error) {
      console.error('Failed to refresh cart:', error);
    }
  };

  const isAuthenticated = !!token;
  const isAdmin = user?.role === 'admin';

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        login,
        logout,
        isAuthenticated,
        isAdmin,
        updateCart,
        cart,
        refreshCart
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};