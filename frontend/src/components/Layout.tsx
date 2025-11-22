// src/components/Layout.tsx
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../context/ThemeContext';
import { Footer } from './Footer';
import { CatalogSidebar } from './catalog/CatalogSidebar';
import './Layout.css'; // Стили для шапки

export const Layout = ({ children }: { children: React.ReactNode }) => {
  const { user, logout, isAuthenticated, isAdmin } = useAuth();
  const { theme, toggleTheme } = useTheme();
  
  return (
    <div className="app-container">
      <header className="app-header">
      <h1 className="app-title">
          <Link to="/">Smartbuy</Link>
        </h1>
        <nav className="main-nav">
          {isAuthenticated ? (
            <>
              <Link to="/" className="nav-link">Товары</Link>
              <Link to="/cart" className="nav-link">Корзина</Link>
               <Link to="/orders" className="nav-link">Заказы</Link>
              {isAdmin && <Link to="/users" className="nav-link">Пользователи</Link>}

              <button onClick={toggleTheme} className="theme-toggle">
                {theme === 'light' ? '🌙' : '☀️'}
              </button>

              <button onClick={logout} className="logout-btn">Выйти ({user?.name})</button>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">Вход</Link>
              <Link to="/signup" className="nav-link">Регистрация</Link>

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