// src/components/cart/CartView.tsx
import { useEffect, useState, useCallback } from 'react';
import { useAuth } from '../../context/AuthContext';
import { getCartItems, updateCartItem, deleteCartItem, getSmartphoneById } from '../../api/client';
import { CartItem, Smartphone } from '../../types';
import './Cart.css'; 
import { motion, AnimatePresence } from 'framer-motion'; // Для анимаций

export const CartView = () => {
  const { cart, token, refreshCart } = useAuth();
  const [items, setItems] = useState<(CartItem & { smartphone?: Smartphone })[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  
  // Загрузка данных корзины и товаров
  useEffect(() => {
    if (!cart?.id || !token) {
      setLoading(false);
      return;
    }

    const loadCartData = async () => {
      try {
        const cartItems = await getCartItems(cart.id, token);
        
        // Загружаем информацию о каждом товаре
        const itemsWithProducts = await Promise.all(
          cartItems.map(async (item) => {
            try {
              const smartphone = await getSmartphoneById(item.smartphone_id, token);
              return { ...item, smartphone };
            } catch {
              return item; // Если не удалось загрузить данные товара
            }
          })
        );

        setItems(itemsWithProducts);
      } catch (error) {
        console.error('Ошибка загрузки корзины:', error);
      } finally {
        setLoading(false);
      }
    };

    loadCartData();
  }, [cart?.items, token]);

  // Подсчет общей суммы
  useEffect(() => {
    const newTotal = items.reduce((sum, item) => {
      const price = item.smartphone?.price || 0;
      return sum + (price * item.quantity);
    }, 0);
    setTotal(newTotal);
  }, [items]);

  const handleUpdateQuantity = useCallback(async (itemId: number, newQuantity: number) => {
    if (!cart?.id || !token || newQuantity < 1) return;
    
    try {
      await updateCartItem(cart.id, itemId, newQuantity, token);
      refreshCart();
    } catch (error) {
      console.error('Ошибка обновления количества:', error);
    }
  }, [cart?.id, token, refreshCart]);

  const handleRemoveItem = useCallback(async (itemId: number) => {
    if (!cart?.id || !token) return;
    
    try {
      await deleteCartItem(cart.id, itemId, token);
      refreshCart();
    } catch (error) {
      console.error('Ошибка удаления товара:', error);
    }
  }, [cart?.id, token, refreshCart]);

  if (loading) return <div className="loading">Загрузка корзины...</div>;
  if (!cart) return <div className="empty">Корзина не найдена</div>;

  return (
    <div className="cart">
      <h2>Моя корзина</h2>
      
      <AnimatePresence>
        {items.length === 0 ? (
          <motion.p 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="empty-message"
          >
            Ваша корзина пуста
          </motion.p>
        ) : (
          <>
            <ul className="cart-items">
              {items.map(item => (
                <motion.li
                  key={item.id}
                  className="cart-item"
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, x: -100 }}
                  transition={{ duration: 0.2 }}
                >
                  <div className="item-info">
                    <div className="product-info">
                      <h3>{item.smartphone?.producer} {item.smartphone?.model || `Товар #${item.smartphone_id}`}</h3>
                      <p className="price">{item.smartphone?.price ? `${item.smartphone.price.toLocaleString('ru-RU')}` : 'Цена не указана'}</p>
                    </div>
                    
                    <div className="quantity-controls">
                      <button 
                        onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                        disabled={item.quantity <= 1}
                      >
                        -
                      </button>
                      <motion.span
                        key={`quantity-${item.id}-${item.quantity}`}
                        initial={{ scale: 1.2 }}
                        animate={{ scale: 1 }}
                        className="quantity"
                      >
                        {item.quantity}
                      </motion.span>
                      <button 
                        onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                      >
                        +
                      </button>
                    </div>
                    
                    <div className="item-total">
                      {item.smartphone?.price ? 
                        `${(item.smartphone.price * item.quantity).toLocaleString('ru-RU')} ₽` : 
                        '—'}
                    </div>
                    
                    <button 
                      onClick={() => handleRemoveItem(item.id)}
                      className="remove-btn"
                    >
                      Удалить
                    </button>
                  </div>
                </motion.li>
              ))}
            </ul>
            
            <motion.div 
              className="cart-summary"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.1 }}
            >
              <h3>Итого: {total.toLocaleString()} ₽</h3>
              <button className="checkout-btn">Оформить заказ</button>
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </div>
  );
};