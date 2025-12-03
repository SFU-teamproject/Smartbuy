// src/components/auth/LoginForm.tsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useLanguage } from '../../context/LanguageContext';
import { login } from '../../api/client';
import { Link } from 'react-router-dom';
import './auth.css'; // Импортируем общие стили для форм авторизации

export const LoginForm = () => {
  const [formData, setFormData] = useState({ 
    email: '', // было name
    password: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login: authLogin } = useAuth();
  const { t } = useLanguage();
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    setError('');
  };

    
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

     // Валидация
    if (!formData.email || !formData.password) {
      setError('Заполните все поля');
      return;
    }
    
    if (!formData.email.includes('@')) {
      setError('Введите корректный email');
      return;
    }
    try {
      setLoading(true);
      const response = await login({
        email: formData.email, // было name: formData.name
        password: formData.password
      });
      
      authLogin(response); // Сохраняем данные в контекст authLogin(response.user, response.token);
      navigate('/'); // Перенаправляем на главную
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка входа';
      //setError(err instanceof Error ? err.message : 'Ошибка входа'); //setError(err instanceof Error ? err.message : 'Login failed');
       // Пользовательское сообщение для временных паролей
      if (errorMessage.includes('Invalid credentials') || errorMessage.includes('invalid credentials')) {
        setError('Неверный email или пароль. Используйте временный пароль из письма.');
        } else {
        setError(errorMessage);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
     <div className="auth-container">
      <div className="auth-card">
        <h2>{t('auth.login')}</h2>
        
        {error && (
          <div className="auth-error">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label htmlFor="email">{t('auth.email') || 'Email'}</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="example@mail.com"
              required
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">{t('auth.password')}</label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              placeholder={t('auth.enterPassword') || 'Введите пароль'}
              required
              disabled={loading}
            />
          </div>
       
          <button 
            type="submit" 
            className="auth-btn"
            disabled={loading}
          >
            {loading ? t('loading') : t('auth.login')}
          </button>
        </form>

        <div className="auth-links">
          <p>
            {t('auth.noAccount')} <Link to="/signup">{t('auth.register')}</Link>
          </p>
        </div>
      </div>
    </div>
  );
};
// Добавляем экспорт по умолчанию для совместимости
export default LoginForm;