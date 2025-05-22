import { useEffect, useState } from 'react';
import { getSmartphoneById } from '../../api/client';
import { useParams } from 'react-router-dom';
import { Smartphone } from '../../types';
import './SmartphoneDetail.css'; 

export function SmartphoneDetail() {
  const [phone, setPhone] = useState<Smartphone | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { id } = useParams(); // Получаем ID из URL

  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!id) return;
        const data = await getSmartphoneById(parseInt(id));
        setPhone(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

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
            <p><strong>Memory:</strong> {phone.memory}GB</p>
            <p><strong>RAM:</strong> {phone.ram}GB</p>
            <p><strong>Display:</strong> {phone.display_size}"</p>
            <p className="price"><strong>Price:</strong> {phone.price.toString()}</p> {/*${phone.price.toLocaleString()}*/}
          </div>
          {phone.ratings_count > 0 && (
            <div className="rating">
              Rating: {(phone.ratings_sum / phone.ratings_count).toFixed(1)}/5
              ({phone.ratings_count} reviews)
            </div>
          )}
          <button className="add-to-cart">Add to Cart</button>
        </div>
      </div>
      {phone.description && (
        <div className="description">
          <h3>Description</h3>
          <p>{phone.description}</p>
        </div>
      )}
    </div>
  );
}