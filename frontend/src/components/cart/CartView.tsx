// src/components/cart/CartView.tsx
import { useEffect, useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { getCartItems, updateCartItem, deleteCartItem } from '../../api/client';
import { CartItem } from '../../types';
import './Cart.css'; // Предполагается, что стили вынесены в отдельный файл

export const CartView = () => {
  const { cart, token, refreshCart } = useAuth();
  const [items, setItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!cart?.id || !token) {
      setLoading(false);
      return;
    }

    const loadCart = async () => {
      try {
        const data = await getCartItems(cart.id, token);
        setItems(data);
      } catch (error) {
        console.error('Ошибка загрузки корзины:', error);
      } finally {
        setLoading(false);
      }
    };

    loadCart();
  }, [cart?.id, token]);

  const handleUpdateQuantity = async (itemId: number, newQuantity: number) => {
    if (!cart?.id || !token || newQuantity < 1) return;
    
    try {
      await updateCartItem(cart.id, itemId, newQuantity, token);
      refreshCart(); // Обновляем корзину через контекст
    } catch (error) {
      console.error('Ошибка обновления количества:', error);
    }
  };

  const handleRemoveItem = async (itemId: number) => {
    if (!cart?.id || !token) return;
    
    try {
      await deleteCartItem(cart.id, itemId, token);
      refreshCart(); // Обновляем корзину через контекст
    } catch (error) {
      console.error('Ошибка удаления товара:', error);
    }
  };

  if (loading) return <div>Загрузка корзины...</div>;
  if (!cart) return <div>Корзина не найдена</div>;

  return (
    <div className="cart">
      <h2>Моя корзина #{cart.id}</h2>
      {items.length === 0 ? (
        <p>Ваша корзина пуста</p>
      ) : (
        <ul className="cart-items">
          {items.map(item => (
            <li key={item.id} className="cart-item">
              <div className="item-info">
                <span>Товар ID: {item.smartphone_id}</span>
                <div className="quantity-controls">
                  <button 
                    onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                    disabled={item.quantity <= 1}
                  >
                    -
                  </button>
                  <span>{item.quantity}</span>
                  <button 
                    onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                  >
                    +
                  </button>
                </div>
                <button 
                  onClick={() => handleRemoveItem(item.id)}
                  className="remove-btn"
                >
                  Удалить
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};