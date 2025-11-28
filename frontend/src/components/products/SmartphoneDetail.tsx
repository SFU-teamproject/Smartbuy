import { ChangeEvent, useEffect, useState } from 'react';
import { addCartItem, addReview, deleteReview, getReviews, getSmartphoneById, updateReview } from '../../api/client';
import { useParams } from 'react-router-dom';
import { Review, ReviewForAdd, ReviewForUpdate, Smartphone } from '../../types';
import { useAuth } from '../../context/AuthContext';
import { Link } from 'react-router-dom';
import './SmartphoneDetail.css';

export function SmartphoneDetail() {
    const [phone, setPhone] = useState<Smartphone | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const { id } = useParams(); // Получаем ID из URL
    const { user, token, refreshCart } = useAuth(); // Хук вызывается в начале компонента
    /* review */
    const [reviews, setReviews] = useState<Review[]>([]);
    const [review, setReview] = useState<ReviewForAdd>({ comment: "", rating: 5 });
    const [myReview, setMyReview] = useState<Review | null>(null);
    const [showReviewFlag, setShowReviewFlag] = useState<Boolean>(false);
    const [selectedRating, setSelectedRating] = useState(5);

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

    useEffect(() => {
        const fetchReview = async () => {
            if (!id) return;
            const data = await getReviews(parseInt(id));
            setReviews(data);
            setShowReviewFlag(data.some(d => d.user_id === user?.id))
            setMyReview(data.filter(d => d.user_id === user?.id)[0])
        };
        fetchReview();
    }, [id, user?.id])

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

    /* review */
    const handleChangeComment = (event: ChangeEvent<HTMLInputElement>) => {
        setReview({ ...review, comment: event.target.value })
    }
    const handleChangeStars = (event: ChangeEvent<HTMLInputElement>) => {
        setReview({ ...review, rating: parseInt(event.target.value) })
        setSelectedRating(parseInt(event.target.value));
    }

    const refreshReviews = async () => {
        if (!id) return;
        const data = await getReviews(parseInt(id));
        setReviews(data);
        setShowReviewFlag(data.some(d => d.user_id === user?.id))
        setMyReview(data.filter(d => d.user_id === user?.id)[0])
    }

    const handleAddReview = async () => {
        if (!id || !token) return;
        try {
            await addReview(parseInt(id), token, review);
            refreshReviews();
        } catch (error) {
            console.error('Failed:', error);
        }
    };
    const handleReviewUpdate = () => {
        setReview({
            comment: myReview?.comment,
            rating: myReview?.rating ? myReview.rating : 0
        })
        setShowReviewFlag(false);
    }
    const handleUpdateReview = async () => {
        if (!id || !token) return;
        try {
            const myreviewforupdate: ReviewForUpdate = {
                id: myReview!.id,
                comment: review?.comment,
                rating: review!.rating

            }
            await updateReview(parseInt(id), token, myreviewforupdate);
            refreshReviews();
        } catch (error) {
            console.error('Failed:', error);
        }
    }
    const handleReviewRemove = async (reviewId: number) => {
        if (!id || !token) return;
        try {
            await deleteReview(parseInt(id), reviewId, token);
            refreshReviews();
        } catch (error) {
            console.error('Failed:', error);
        }
    }
    const dateConvert = (timestamp: Date) => {
        let currentDate = new Date(timestamp);
        return currentDate.toDateString();
    }

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
            {!showReviewFlag ? (
                <div className='create_review'>
                    <h2>Оставьте свой отзыв:</h2>
                    <label htmlFor='review_comment'>Комментарий: </label>
                    <input type="text" id="review_comment" value={review?.comment} onChange={handleChangeComment} />
                    <div className='review_stars'>
                        <span>Оценка: </span>
                        <label htmlFor='review_star_1'>1</label>
                        <input type='radio' name='star' id='review_star_1' value={1} onChange={handleChangeStars} checked={selectedRating === 1} ></input>
                        <label htmlFor='review_star_2'>2</label>
                        <input type='radio' name='star' id='review_star_2' value={2} onChange={handleChangeStars} checked={selectedRating === 2} ></input>
                        <label htmlFor='review_star_3'>3</label>
                        <input type='radio' name='star' id='review_star_3' value={3} onChange={handleChangeStars} checked={selectedRating === 3} ></input>
                        <label htmlFor='review_star_4'>4</label>
                        <input type='radio' name='star' id='review_star_4' value={4} onChange={handleChangeStars} checked={selectedRating === 4} ></input>
                        <label htmlFor='review_star_5'>5</label>
                        <input type='radio' name='star' id='review_star_5' value={5} onChange={handleChangeStars} checked={selectedRating === 5} ></input>
                    </div>
                    {!myReview ?
                        (<button type='submit' onClick={handleAddReview} className='review_button'>Создать</button>) :
                        (<button type='submit' onClick={handleUpdateReview} className='review_button'>Обновить</button>)
                    }
                    <br />
                </div>
            ) : (
                <div className='Update_review'>
                    <h2>Ваш отзыв:</h2>
                    <div className='reviews_container'>
                        {myReview && (
                            <div className='review'>
                                {myReview.user_name ? <div>user name: {(myReview.user_name)}</div> : ''}
                                <div>rating: {(myReview?.rating)} </div>
                                {myReview.comment ? <div>comment: {(myReview.comment)}</div> : ''}
                                <div>created at: {(dateConvert(myReview.created_at))}</div>
                                {myReview.updated_at ? <div>updated at: {(dateConvert(myReview.updated_at))}</div> : ''}
                                <div>
                                    <button className='review_button' onClick={handleReviewUpdate}>обновить</button>
                                    <button type='submit' className='review_button' onClick={e => handleReviewRemove(myReview.id)}>удалить</button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            )}
            {reviews.length > 0 ? <h2>Отзывы:</h2> : <h2>Отзывов нет</h2>}
            {reviews.map((review) => (
                <div className='reviews_container'>
                    <div className='review'>
                        {review.user_name ? <div>user name: {(review.user_name)}</div> : ''}
                        <div>rating: {(review.rating)}</div>
                        {review.comment ? <div>comment: {(review.comment)}</div> : ''}
                        <div>created at: {(dateConvert(review.created_at))}</div>
                        {review.updated_at ? <div>updated at: {(dateConvert(review.updated_at))}</div> : ''}
                    </div>
                    <br />
                    <hr />
                </div>
            ))}
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