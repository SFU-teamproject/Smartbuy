import { Link } from 'react-router-dom';
import './CatalogSidebar.css';

export const CatalogSidebar = () => {
  const producers = [
    { name: 'Apple', slug: 'apple' },
    { name: 'Samsung', slug: 'samsung' },
    { name: 'Xiaomi', slug: 'xiaomi' },
    { name: 'Huawei', slug: 'huawei' },
    { name: 'Google', slug: 'google' },
    { name: 'POCO', slug: 'poco' },
    { name: 'Honor', slug: 'honor' }
  ];

  return (
    <aside className="catalog-sidebar">
      <h3 className="catalog-title">Каталог товаров</h3>
      <nav className="producers-nav">
        <ul className="producers-list">
          {producers.map((producer) => (
            <li key={producer.slug} className="producer-item">
              <Link 
                to={`/?producer=${producer.slug}`}
                className="producer-link"
              >
                <span className="producer-name">{producer.name}</span>
                <span className="producer-arrow">›</span>
              </Link>
            </li>
          ))}
        </ul>
      </nav>

      {/* Дополнительные фильтры */}
      <div className="additional-filters">
        <h4>Фильтры по цене</h4>
        <div className="filter-group">
          <Link to="/?price=budget" className="filter-link">
            До 20 000 ₽
          </Link>
          <Link to="/?price=mid" className="filter-link">
            20 000 - 40 000 ₽
          </Link>
          <Link to="/?price=premium" className="filter-link">
            От 40 000 ₽
          </Link>
          <Link to="/" className="filter-link">
            Все товары
          </Link>
        </div>
      </div>
    </aside>
  );
};

export default CatalogSidebar;