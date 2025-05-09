// src/components/SmartphoneList.tsx
import { useEffect, useState } from 'react';
import { getSmartphones } from '../api/client';
import { Smartphone } from '../types';
import './SmartphoneList.css'; 


export function SmartphoneList() {
  const [smartphones, setSmartphones] = useState<Smartphone[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        console.log('Fetching smartphones...'); // Логирование
        const data = await getSmartphones();
        console.log('Received data:', data); // Проверка данных
        setSmartphones(data);
      } catch (err) {
        console.error('Fetch error:', err);
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">Error: {error}</div>;

  return (
    <div className="smartphone-list">
      <h2>Our Smartphones</h2>
      <div className="products-grid">
        {smartphones.map((phone) => (
          <div key={phone.id} className="product-card">
            <img 
              src={phone.image_path || '/placeholder-phone.jpg'} 
              alt={phone.model}
              className="product-image"
            />
            <div className="product-info">
              <h3>{phone.producer} {phone.model}</h3>
              <p>Memory: {phone.memory}GB</p>
              <p>RAM: {phone.ram}GB</p>
              <p className="price">${phone.price.toLocaleString()}</p>
              <div className="rating">
                Rating: {phone.ratings_count > 0 
                  ? (phone.ratings_sum / phone.ratings_count).toFixed(1) 
                  : 'No ratings'}
              </div>
              <button className="add-to-cart">Add to Cart</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}