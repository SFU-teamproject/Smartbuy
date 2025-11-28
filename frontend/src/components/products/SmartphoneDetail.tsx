import { useEffect, useState } from 'react';
import { getSmartphoneById } from '../../api/client';
import { useParams } from 'react-router-dom';
import { Smartphone } from '../../types';
import { useAuth } from '../../context/AuthContext';
import { addCartItem } from '../../api/client';
import { Link } from 'react-router-dom';
import './SmartphoneDetail.css'; 

export function SmartphoneDetail() {
  const [phone, setPhone] = useState<Smartphone | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { id } = useParams(); // Получаем ID из URL
  const { user, token, refreshCart } = useAuth(); // Хук вызывается в начале компонента

  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!id) return;
        const smartphone = await getSmartphoneById(parseInt(id));
        setPhone(smartphone);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Не удалось загрузить информацию о товаре');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  const handleAddToCart = async (smartphoneId: number) => {
     if (!user?.cart?.id || !token) {
      alert('Необходимо авторизоваться для добавления в корзину');
      return;
    }
      try {
        await addCartItem(user.cart.id, smartphoneId, token);
        alert('Item added to cart!');
        refreshCart();
      } catch (error) {
        console.error('Failed to add to cart:', error);
      }
    };

  const inBucket = (smartphoneId: number) => {
    const items = user?.cart?.items;
    if (items) {
      if (items.find(item => item.smartphone_id === smartphoneId)) {
        return true;
      } 
    }
    return false;
  };

   if (loading) return <div className="loading">Загрузка...</div>;
  if (error) return <div className="error">Ошибка: {error}</div>;
  if (!phone) return <div className="not-found">Товар не найден</div>;

  return (
    <div className="smartphone-detail">
      {/* Хлебные крошки */}
     {/*<nav className="breadcrumbs">
        <Link to="/" className="breadcrumb-link">Главная</Link>
        <span className="breadcrumb-separator">/</span>
        <Link to="/" className="breadcrumb-link">Смартфоны</Link>
        <span className="breadcrumb-separator">/</span>
        <span className="breadcrumb-current">{phone.producer} {phone.model}</span>
      </nav> */}

      <div className="detail-container">
        <div className="image-section">
          <img 
            src={phone.image_path || '/placeholder-phone.jpg'} 
            alt={phone.model} 
            className="detail-image"
          />
        </div>
        <div className="detail-info">
          <div className="product-header">
            <h1 className="product-title">{phone.producer} {phone.model}</h1>
            {phone.ratings_count > 0 && (
              <div className="rating-badge">
                <span className="rating-stars">⭐</span>
                <span className="rating-value">
                  {(phone.ratings_sum / phone.ratings_count).toFixed(1)}
                </span>
                <span className="rating-count">({phone.ratings_count})</span>
              </div>
            )}
          </div>

          <div className="price-section">
            <span className="price">{phone.price.toLocaleString('ru-RU')} ₽</span>
          </div>

          <div className="specs-grid">
            <div className="spec-item">
              <span className="spec-label">Память</span>
              <span className="spec-value">{phone.memory} GB</span>
            </div>
            <div className="spec-item">
              <span className="spec-label">Оперативная память</span>
              <span className="spec-value">{phone.ram} GB</span>
            </div>
            <div className="spec-item">
              <span className="spec-label">Диагональ экрана</span>
              <span className="spec-value">{phone.display_size}"</span>
            </div>
            <div className="spec-item">
              <span className="spec-label">Производитель</span>
              <span className="spec-value">{phone.producer}</span>
            </div>
          </div>

          <div className="action-section">
            <button 
              className={`add-to-cart ${inBucket(phone.id) ? 'in-cart' : ''}`}
              onClick={() => handleAddToCart(phone.id)} 
              disabled={inBucket(phone.id)}
            >
              {inBucket(phone.id) ? (
                <>
                  <span className="cart-icon">✓</span>
                  Уже в корзине
                </>
              ) : (
                <>
                  <span className="cart-icon">🛒</span>
                  Добавить в корзину
                </>
              )}
            </button>
            
          </div>
        </div>
      </div>

      {phone.description && (
        <div className="description-section">
          <h2 className="description-title">Описание</h2>
          <div className="description-content">
            <p>{phone.description}</p>
          </div>
        </div>
      )}
    </div>
  );
}
       {/*  
        <div className="detail-info">
           <div className="product-header">
            <h2>{phone.producer} {phone.model}</h2>
          <div className="specs">
            <p><strong>Память:</strong> {phone.memory}GB</p>
            <p><strong>RAM:</strong> {phone.ram}GB</p>
            <p><strong>Экран:</strong> {phone.display_size}"</p>
            <p className="price"><strong>Цена:</strong> {phone.price.toString()}</p> {/*${phone.price.toLocaleString()}*/} /*
          </div>
          {phone.ratings_count > 0 && (
            <div className="rating">
              Рейтинг: {(phone.ratings_sum / phone.ratings_count).toFixed(1)}/5
              ({phone.ratings_count} reviews)
            </div>
          )}
          <button className="add-to-cart" onClick={() => handleAddToCart(phone.id)} disabled={inBucket(phone.id)}>{inBucket(phone.id) ? "Уже в корзине" : "Добавить в корзину"}</button>
        </div>
      </div>
      {phone.description && (
        <div className="description">
          <h3>Описание</h3>
          <p>{phone.description}</p>
        </div>
      )}
    </div>
  );
}
}*/