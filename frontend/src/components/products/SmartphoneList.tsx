// src/components/products/SmartphoneList.tsx
import { useEffect, useState } from 'react';
import { getSmartphones, getSmartphonesByIds } from '../../api/client';
import { Smartphone } from '../../types';
import './SmartphoneList.css'; 
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { addCartItem } from '../../api/client';

export function SmartphoneList() {
  // Все хуки вызываются в начале, до любых условий
  const [smartphones, setSmartphones] = useState<Smartphone[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [filteredIds, setFilteredIds] = useState<number[] | null>(null);
  const { user, token } = useAuth(); // Хук вызывается в начале компонента

  useEffect(() => {
    const fetchData = async () => {
      try {
        console.log('Fetching smartphones...');
        const data = filteredIds 
          ? await getSmartphonesByIds(filteredIds)
          : await getSmartphones();
        console.log('Received data:', data);
        setSmartphones(data);
      } catch (err) {
        console.error('Fetch error:', err);
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [filteredIds]);

  const handleShowPopular = () => {
    setFilteredIds([1, 3, 4]);
  };

  const handleResetFilter = () => {
    setFilteredIds(null);
  };

  const handleAddToCart = async (smartphoneId: number) => {
    if (!user?.cart?.id || !token) return;
    
    try {
      await addCartItem(user.cart.id, smartphoneId, token);
      alert('Item added to cart!');
    } catch (error) {
      console.error('Failed to add to cart:', error);
    }
  };

  // Условный рендеринг после всех хуков
  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">Error: {error}</div>;

  return (
    <div className="smartphone-list">
      <h2>Our Smartphones</h2>
      <div className="filter-controls">
        <button onClick={handleShowPopular} className="filter-button">
          Show Popular (IDs: 1, 3, 4)
        </button>
        {filteredIds && (
          <button onClick={handleResetFilter} className="filter-button">
            Show All
          </button>
        )}
      </div>

      <div className="products-grid">
        {smartphones.map((phone) => (
          <div key={phone.id} className="product-card">
            <img 
              src={phone.image_path || '/placeholder-phone.jpg'} 
              alt={phone.model}
              className="product-image"
            />
            <div className="product-info">
              <Link to={`/smartphones/${phone.id}`} className="product-link">
                <h3>{phone.producer} {phone.model}</h3>
              </Link>
              <p>Memory: {phone.memory}GB</p>
              <p>RAM: {phone.ram}GB</p>
              <p className="price">{phone.price.toLocaleString('ru-RU')} </p>
              <div className="rating">
                Rating: {phone.ratings_count > 0 
                  ? (phone.ratings_sum / phone.ratings_count).toFixed(1) 
                  : 'No ratings'}
              </div>
              <button 
                onClick={() => handleAddToCart(phone.id)} 
                className="add-to-cart"
              >
                Добавить в корзину
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}