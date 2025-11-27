import { useEffect, useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { getOrders, cancelOrder } from '../../api/client';
import { Order } from '../../types';
import './OrderHistory.css';

export const OrderHistory = () => {
  const { user, token } = useAuth();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchOrders = async () => {
      if (!token){
        setError('Необходима авторизация');
        setLoading(false);
        return;
      };
      
      try {
        setLoading(true);
        setError('');
        const ordersData = await getOrders(token);
        setOrders(ordersData);
      } catch (err) {
        console.error('Error fetching orders:', err);
        setError(err instanceof Error ? err.message : 'Ошибка загрузки заказов');
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, [token]);

  const handleCancelOrder = async (orderId: number) => {
    if (!token) return;
    
    if (!window.confirm('Вы уверены, что хотите отменить заказ?')) return;
    
    try {
      await cancelOrder(orderId, token);
      setOrders(orders.map(order => 
        order.id === orderId 
          ? { ...order, status: 'cancelled' }// ? { ...order, status: 'cancelled', updated_at: new Date().toISOString() }
          : order
      ));
    } catch (err) {
      alert('Не удалось отменить заказ');
    }
  };

  const getStatusInfo = (status: string) => {
    const statusInfo = {
      pending: { text: 'Ожидает обработки', color: '#f39c12' },
      processing: { text: 'В обработке', color: '#3498db' },
      shipped: { text: 'Отправлен', color: '#9b59b6' },
      delivered: { text: 'Доставлен', color: '#27ae60' },
      cancelled: { text: 'Отменен', color: '#e74c3c' }
    };
    
    return statusInfo[status as keyof typeof statusInfo] || { text: status, color: '#95a5a6' };
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (loading) return(
      <div className="order-history">
        <div className="loading">Загрузка истории заказов...</div>
      </div>
    );
  if (error) {
    return (
      <div className="order-history">
        <div className="error">
          <h3>Произошла ошибка</h3>
          <p>{error}</p>
          <button onClick={() => window.location.reload()} className="retry-btn">
            Попробовать снова
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="order-history">
      <div className="order-history-header">
        <h2>История заказов</h2>
        <p>Здесь вы можете просмотреть историю ваших заказов</p>
      </div>

      {orders.length === 0 ? (
        <div className="empty-orders">
          <div className="empty-icon">📦</div>
          <h3>Заказов пока нет</h3>
          <p>Совершите свой первый заказ, и он появится здесь</p>
        </div>
      ) : (
        <div className="orders-list">
          {orders.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()).map(order => (
            <div key={order.id} className="order-card">
              <div className="order-header">
                <div className="order-info">
                  <h3>Заказ #{order.id}</h3>
                  <span className="order-date">{formatDate(order.created_at)}</span>
                </div>
                <div className="order-status">
                  <span 
                    className="status-badge"
                    style={{ backgroundColor: getStatusInfo(order.status).color }}
                  >
                    {getStatusInfo(order.status).text}
                  </span>
                  <span className="order-total">
                    {order.total_amount.toLocaleString('ru-RU')} ₽
                  </span>
                </div>
              </div>

              <div className="order-items">
                <h4>Товары:</h4>
                {order.items.map(item => (
                  <div key={item.id} className="order-item">
                    <div className="item-info">
                      <span className="item-name">
                        {item.smartphone?.producer} {item.smartphone?.model}
                      </span>
                      <span className="item-quantity">× {item.quantity}</span>
                    </div>
                    <span className="item-price">
                      {item.price.toLocaleString('ru-RU')} ₽
                    </span>
                  </div>
                ))}
              </div>

              <div className="order-actions">
                {(order.status === 'pending' || order.status === 'processing') && (
                  <button 
                    onClick={() => handleCancelOrder(order.id)}
                    className="cancel-btn"
                  >
                    Отменить заказ
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};