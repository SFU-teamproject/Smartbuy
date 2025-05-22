// src/components/Layout.tsx
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Layout.css'; // Стили для шапки

export const Layout = ({ children }: { children: React.ReactNode }) => {
  const { user, logout, isAuthenticated, isAdmin } = useAuth();
  
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
              {isAdmin && <Link to="/users" className="nav-link">Пользователи</Link>}
              <button onClick={logout} className="logout-btn">Выйти ({user?.name})</button>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">Вход</Link>
              <Link to="/signup" className="nav-link">Регистрация</Link>
            </>
          )}
        </nav>
      </header>
      <main  className="app-content">{children}</main>
    </div>
  );
};