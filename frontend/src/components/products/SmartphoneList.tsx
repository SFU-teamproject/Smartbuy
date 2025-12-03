// src/components/products/SmartphoneList.tsx
import { useEffect, useState } from 'react';
import { getSmartphones, getSmartphonesByIds } from '../../api/client';
import { Smartphone } from '../../types';
import './SmartphoneList.css'; 
import { Link, useSearchParams } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { addCartItem } from '../../api/client';
import { useLanguage } from '../../context/LanguageContext';
import Pagination from './Pagination/Pagination';

export function SmartphoneList() {
  // Все хуки вызываются в начале, до любых условий
  const [smartphones, setSmartphones] = useState<Smartphone[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [filteredIds, setFilteredIds] = useState<number[] | null>(null);
  const { user, token, refreshCart } = useAuth(); // Хук вызывается в начале компонента
  const [searchParams, setSearchParams] = useSearchParams();
  const { t } = useLanguage();

  // Получаем параметры фильтрации из URL
  const producerFilter = searchParams.get('producer');
  const priceFilter = searchParams.get('price');

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

  // Функция для фильтрации смартфонов по производителю и цене
  const getFilteredSmartphones = () => {
    let filtered = smartphones;

    // Фильтрация по производителю
    if (producerFilter) {
      filtered = filtered.filter(phone => 
        phone.producer.toLowerCase().includes(producerFilter.toLowerCase())
      );
    }

    // Фильтрация по цене
    if (priceFilter) {
      switch (priceFilter) {
        case 'budget':
          filtered = filtered.filter(phone => phone.price < 20000);
          break;
        case 'mid':
          filtered = filtered.filter(phone => phone.price >= 20000 && phone.price <= 40000);
          break;
        case 'premium':
          filtered = filtered.filter(phone => phone.price > 40000);
          break;
        default:
          break;
      }
    }

    return filtered;
  };

  const handleShowPopular = () => {
    setFilteredIds([1, 3, 4]);
     // Сбрасываем другие фильтры при активации популярных
    setSearchParams({});
  };

  const handleResetFilter = () => {
    setFilteredIds(null);
    setSearchParams({});
  };

  const handleAddToCart = async (smartphoneId: number) => {
    if (!user?.cart?.id || !token) return;
    
    try {
      await addCartItem(user.cart.id, smartphoneId, token);
      alert('Товар добавлен в корзину!');
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
  
   // Получаем отфильтрованные смартфоны
  const filteredSmartphones = getFilteredSmartphones();


  const ITEMS_PER_PAGE = 8;
  const [currentPage, setCurrentPage] = useState<number>(1);
  //const totalPages = Math.ceil(smartphones.length / ITEMS_PER_PAGE);
  const totalPages = Math.ceil(filteredSmartphones.length / ITEMS_PER_PAGE); // Исправлено: используем filteredSmartphones
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  //const currentSmartphones = smartphones.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  const currentSmartphones = filteredSmartphones.slice(startIndex, startIndex + ITEMS_PER_PAGE); // Исправлено: используем filteredSmartphones
  const pageNumbers = Array.from({ length: totalPages }, (_, index) => index + 1);

  // Сбрасываем страницу при изменении фильтров
  useEffect(() => {
    setCurrentPage(1);
  }, [producerFilter, priceFilter, filteredIds]);

  const handlePrevPage = () => {
    setCurrentPage((prev) => Math.max(prev - 1, 1));
  };

  const handleNextPage = () => {
    setCurrentPage((prev) => Math.min(prev + 1, totalPages));
  };

  const handlePageClick = (pageNumber: number) => {
    setCurrentPage(pageNumber);
  };

  // Функция для отображения текущих фильтров
  const getActiveFiltersInfo = () => {
  const filters = [];
  if (producerFilter) filters.push(`Производитель: ${producerFilter}`);
  if (priceFilter) {
    const priceLabels: Record<string, string> = {
      budget: 'До 20 000 ₽',
      mid: '20 000 - 40 000 ₽', 
      premium: 'От 40 000 ₽'
    };
    filters.push(`Цена: ${priceLabels[priceFilter] || priceFilter}`);
  }
  if (filteredIds) filters.push('Популярные товары');
  
  return filters;
};

  const activeFilters = getActiveFiltersInfo();

   // Условный рендеринг после всех хуков
  if (loading) return <div className="loading">Загрузка...</div>;
  if (error) return <div className="error">Ошибка: {error}</div>;

  return (
    <div className="smartphone-list">
      <h2>{t('products.ourSmartphones')}</h2>
       {/* <div className="filter-controls">
        <button onClick={handleShowPopular} className="filter-button">
          Show Popular (IDs: 1, 3, 4)
        </button>
        {filteredIds && (
          <button onClick={handleResetFilter} className="filter-button">
            Show All
          </button>
        )}
      </div>

      Информация о активных фильтрах
      {activeFilters.length > 0 && (
        <div className="active-filters">
          <h3></h3>
          <div className="filters-list">
            {activeFilters.map((filter, index) => (
              <span key={index} className="filter-tag">
                {filter}
              </span>
            ))}
            <button onClick={handleResetFilter} className="clear-filters">
             X
            </button>
          </div>
        </div>
      )}  */}
      {/* {/* Информация о активных фильтрах расширенная*/}{/*
      {activeFilters.length > 0 && (
        <div className="active-filters">
          <div className="filters-header">
            <span className="filters-title">Активные фильтры:</span>
            <button onClick={handleResetFilter} className="clear-all-filters">
              Очистить все
            </button>
          </div>
          <div className="filters-list">
            {activeFilters.map((filter, index) => (
              <div key={index} className="filter-tag">
                <span className="filter-text">{filter}</span>
                <button 
                  onClick={() => {
                    // Удаляем конкретный фильтр
                    if (filter.includes('Производитель')) {
                      searchParams.delete('producer');
                    } else if (filter.includes('Цена')) {
                      searchParams.delete('price');
                    } else if (filter.includes('Популярные')) {
                      setFilteredIds(null);
                    }
                    setSearchParams(searchParams);
                  }}
                  className="remove-filter"
                >
                  ×
                </button>
              </div>
            ))}
          </div>
        </div>
      )}*/}

      <div className="filter-controls">
        <button onClick={handleShowPopular} className="filter-button">
          Показать популярные
        </button>
        {filteredIds && (
          <button onClick={handleResetFilter} className="filter-button">
            Показать все
          </button>
        )}
      </div>

      <div className="products-info">
        <p>{t('products.found')}: {filteredSmartphones.length}</p>
      </div>
         {/* Основной контент */}
       <div className="products-grid">
        {currentSmartphones.map((phone) => (
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
              <div className="product-specs">
                <p>{t('products.memory')}: {phone.memory}GB </p>
                <p>{t('products.ram')}: {phone.ram}GB </p>
               </div>
              <p className="price">{phone.price.toLocaleString('ru-RU')} ₽</p>
              <div className="rating">
                {t('products.rating')}: {phone.ratings_count > 0 
                  ? (phone.ratings_sum / phone.ratings_count).toFixed(1) 
                  : t('products.noRatings')}
              </div>
              <button 
                onClick={() => handleAddToCart(phone.id)} 
                className="add-to-cart"
                disabled={inBucket(phone.id)}
              >
                {inBucket(phone.id) ? t('products.inCart') : t('products.addToCart')}
              </button>
            </div>
          </div>
        ))}
      </div>

      {filteredSmartphones.length > 0 && (
        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={handlePageClick}
          onPrevPage={handlePrevPage}
          onNextPage={handleNextPage}
        />
      )}
      
      {filteredSmartphones.length === 0 && !loading && (
        <div className="no-products">
          <h3>{t('products.notFound')}</h3>
          <button onClick={handleResetFilter} className="filter-button">
            Показать все товары
          </button>
        </div>
      )}
    </div>
  );
}