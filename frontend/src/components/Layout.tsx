// src/components/Layout.tsx
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../context/ThemeContext';
import { useLanguage } from '../context/LanguageContext';
import { Footer } from './Footer';
import { CatalogSidebar } from './catalog/CatalogSidebar';
import './Layout.css'; // Стили для шапки

export const Layout = ({ children }: { children: React.ReactNode }) => {
  const { user, logout, isAuthenticated, isAdmin } = useAuth();
  const { theme, toggleTheme } = useTheme();
   const { language, setLanguage, t } = useLanguage();

  const handleLanguageChange = () => {
    setLanguage(language === 'ru' ? 'en' : 'ru');
  };
  
  return (
    <div className="app-container">
      <header className="app-header">
      <h1 className="app-title">
          <Link to="/">{t('app.title')}</Link>
        </h1>
        <nav className="main-nav">
          {isAuthenticated ? (
            <>
              <Link to="/" className="nav-link">{t('nav.products')}</Link>
              <Link to="/cart" className="nav-link">{t('nav.cart')}</Link>
               <Link to="/orders" className="nav-link">{t('nav.orders')}</Link>
              {isAdmin && <Link to="/users" className="nav-link">{t('nav.users')}</Link>}

              <button onClick={handleLanguageChange} className="language-toggle">
                {language === 'ru' ? 'EN' : 'RU'}
              </button>

              <button onClick={toggleTheme} className="theme-toggle">
                {theme === 'light' ? '🌙' : '☀️'}
              </button>

              <button onClick={logout} className="logout-btn">{t('nav.logout')} ({user?.name})</button>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">{t('nav.login')}</Link>
              <Link to="/signup" className="nav-link">{t('nav.signup')}</Link>

              <button onClick={toggleTheme} className="theme-toggle">
                {theme === 'light' ? '🌙' : '☀️'}
              </button>
            </>
          )}
        </nav>
      </header>
         <div className="main-layout">
        {isAuthenticated && (
          <CatalogSidebar />
        )}
        <main className="app-content">
          {children}
        </main>
      </div>
      
      <Footer />
    </div>
  );
};