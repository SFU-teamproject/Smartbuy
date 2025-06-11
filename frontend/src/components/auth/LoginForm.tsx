// src/components/auth/LoginForm.tsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { login } from '../../api/client';

export const LoginForm = () => {
  const [formData, setFormData] = useState({ name: '', password: '' });
  const [error, setError] = useState('');
  const { login: authLogin } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const response = await login(formData);
      authLogin(response); // Сохраняем данные в контекст
      navigate('/'); // Перенаправляем на главную
    } catch (err) {
      setError('Неверные учетные данные или ошибка сервера'); //setError(err instanceof Error ? err.message : 'Login failed');
    }
  };

  return (
    <div className="login-form">
    <h2>Вход</h2>
    {error && <div className="error-message">{error}</div>}
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={formData.name}
        onChange={(e) => setFormData({...formData, name: e.target.value})}
        placeholder="Логин"
        required
      />
      <input
        type="password"
        value={formData.password}
        onChange={(e) => setFormData({...formData, password: e.target.value})}
        placeholder="Пароль"
        required
      />
      {error && <div className="error">{error}</div>}
      <button type="submit">Войти</button>
    </form>
    </div>
  );
};