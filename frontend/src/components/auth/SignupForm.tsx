// src/components/auth/SignupForm.tsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
//import { useAuth } from '../../context/AuthContext';
import { useLanguage } from '../../context/LanguageContext';
import { signup } from '../../api/client';
import './auth.css';
import { Link } from 'react-router-dom';

export const SignupForm = () => {
  const [formData, setFormData] = useState({
    email: ''
    //password: '',
    //confirmPassword: ''
  });
  const [error, setError] = useState('');
  //const { login: authLogin } = useAuth();
  const navigate = useNavigate();
  const { t } = useLanguage();
  const [loading, setLoading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    setError('');
    setSuccessMessage('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccessMessage('');

     // Валидация
    if (!formData.email) {
      setError('Введите email');
      return;
    }

    if (!formData.email.includes('@')) {
      setError('Введите корректный email');
      return;
    }

    
    try {
      setLoading(true);
      // Отправляем только email - бэкенд сгенерирует временный пароль
      await signup({ email: formData.email });
      
      // Успешная регистрация - показываем сообщение
      setSuccessMessage('На ваш email отправлен временный пароль. Используйте его для входа.');
      
      // Очищаем форму
      setFormData({ email: '' });
      
      // Автоматически перенаправлять не нужно - пользователь должен проверить почту
      // navigate('/login');
    } catch (err) {
       setError(err instanceof Error ? err.message : 'Ошибка регистрации');
    }finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
    {/*<div className="signup-form">*/}
    <div className="auth-card">
      <h2>{t('auth.register')}</h2>

       {successMessage && (
          <div className="auth-success">
            {successMessage}
            <div className="success-instructions">
              <p>1. Проверьте вашу почту</p>
              <p>2. Найдите письмо с временным паролем</p>
              <p>3. Войдите с этим паролем</p>
              <p>4. После входа вы сможете сменить пароль в настройках профиля</p>
            </div>
          </div>
        )}

      {error && <div className="auth-error">{error}</div>}

      {!successMessage && (
          <>
      <form onSubmit={handleSubmit}  className="auth-form">
         <div className="form-group">
        <label htmlFor="email">{t('auth.email') || 'Email'}</label>
        <input
          type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="example@example.com"
              required
              disabled={loading}
        />
        </div>

        <div className="form-note">
                <p>После регистрации на ваш email будет отправлен временный пароль</p>
        </div>
        
        <button 
          type="submit"
          className="auth-btn"
          disabled={loading}>
            {loading ? t('loading') : t('auth.register')}
        </button>
      </form>

      <div className="auth-links">
          <p>
            {t('auth.haveAccount')} <Link to="/login">{t('auth.login')}</Link>
          </p>
      </div>
      </>
        )}
    </div>
    </div>
  );
};
export default SignupForm;