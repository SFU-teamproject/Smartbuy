import { useEffect, useState } from 'react';
import { getSmartphoneById } from '../../api/client';
import { useParams } from 'react-router-dom';
import { Smartphone } from '../../types';
import { useAuth } from '../../context/AuthContext';
import { addCartItem } from '../../api/client';
import './SmartphoneDetail.css'; 

export function SmartphoneDetail() {
  const [phone, setPhone] = useState<Smartphone | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { id } = useParams(); // Получаем ID из URL
  const { user, token } = useAuth(); // Хук вызывается в начале компонента

  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!id) return;
        const smartphone = await getSmartphoneById(parseInt(id));
        setPhone(smartphone);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  const handleAddToCart = async (smartphoneId: number) => {
      if (!user?.cart?.id || !token) return;
      try {
        await addCartItem(user.cart.id, smartphoneId, token);
        alert('Item added to cart!');
      } catch (error) {
        console.error('Failed to add to cart:', error);
      }
    };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!phone) return <div>Smartphone not found</div>;

  return (
    <div className="smartphone-detail">
      <div className="detail-container">
        <img 
          src={phone.image_path || '/placeholder-phone.jpg'} 
          alt={phone.model} 
          className="detail-image"
        />
        <div className="detail-info">
          <h2>{phone.producer} {phone.model}</h2>
          <div className="specs">
            <p><strong>Память:</strong> {phone.memory}GB</p>
            <p><strong>RAM:</strong> {phone.ram}GB</p>
            <p><strong>Экран:</strong> {phone.display_size}"</p>
            <p className="price"><strong>Цена:</strong> {phone.price.toString()}</p> {/*${phone.price.toLocaleString()}*/}
          </div>
          {phone.ratings_count > 0 && (
            <div className="rating">
              Рейтинг: {(phone.ratings_sum / phone.ratings_count).toFixed(1)}/5
              ({phone.ratings_count} reviews)
            </div>
          )}
          <button className="add-to-cart" onClick={() => handleAddToCart(phone.id)} >Add to Cart</button>
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