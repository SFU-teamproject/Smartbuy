// src/components/auth/SignupForm.tsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { signup } from '../../api/client';

export const SignupForm = () => {
  const [formData, setFormData] = useState({ name: '', password: '' });
  const [error, setError] = useState('');
  const { login: authLogin } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    
    try {
      const response = await signup(formData);
      //authLogin(response); // Автоматически входим после регистрации
      navigate('/login');
    } catch (err) {
      setError('Пользователь уже существует или ошибка сервера');
    }
  };

  return (
    <div className="signup-form">
      <h2>Регистрация</h2>
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
        <button type="submit">Зарегистрироваться</button>
      </form>
    </div>
  );
};