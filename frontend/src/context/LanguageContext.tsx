import React, { createContext, useContext, useEffect, useState } from 'react';

type Language = 'ru' | 'en';

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

// Словари переводов
const translations = {
  ru: {
    // Общие
    'app.title': 'Smartbuy',
    'loading': 'Загрузка...',
    'error': 'Ошибка',
    'notFound': 'Не найдено',
    'save': 'Сохранить',
    'cancel': 'Отмена',
    'delete': 'Удалить',
    'edit': 'Редактировать',
    'yes': 'Да',
    'no': 'Нет',
    
    // Навигация
    'nav.products': 'Товары',
    'nav.cart': 'Корзина',
    'nav.orders': 'Заказы',
    'nav.users': 'Пользователи',
    'nav.login': 'Вход',
    'nav.signup': 'Регистрация',
    'nav.logout': 'Выйти',
    'nav.continueShopping': 'Продолжить покупки',
    'nav.backToShopping': '← Вернуться к покупкам',
    
    // Авторизация
    'auth.login': 'Войти',
    'auth.register': 'Зарегистрироваться',
    'auth.name': 'Имя',
    'auth.password': 'Пароль',
    'auth.confirmPassword': 'Подтвердите пароль',
    'auth.noAccount': 'Нет аккаунта?',
    'auth.haveAccount': 'Уже есть аккаунт?',
    'auth.email': 'Email',
    'auth.enterEmail': 'Введите email',
    'auth.enterPassword': 'Введите пароль',
    
    // Товары
    'products.ourSmartphones': 'Наши смартфоны',
    'products.memory': 'Память',
    'products.ram': 'Оперативная память',
    'products.screen': 'Диагональ экрана',
    'products.price': 'Цена',
    'products.rating': 'Оценка',
    'products.noRatings': 'Нет оценок',
    'products.addToCart': 'Добавить в корзину',
    'products.inCart': 'Уже в корзине',
    'products.popular': 'Популярные товары',
    'products.showPopular': 'Показать популярные',
    'products.showAll': 'Показать все',
    'products.showAllprod': 'Показать все товары',
    'products.found': 'Найдено товаров',
    'products.notFound': 'Товары не найдены',
    'products.changeFilters': 'Попробуйте изменить параметры фильтрации',
    'products.resetFilters': 'Показать все товары',
    'products.producer': 'Производитель',
    'products.specs': 'Характеристики',
    'products.description': 'Описание',
    'products.addRewiew': 'Оставьте свой отзыв',
    'products.comment': 'Комментарий',
    'products.create': 'Создать',
    'products.update': 'Обновить',
    'products.youReview': 'Ваш отзыв',
    'products.reviews': 'Отзывы',

    
    // Фильтры
    'filters.title': 'Каталог товаров',
    'filters.active': 'Активные фильтры',
    'filters.clearAll': 'Очистить все',
    'filters.price': 'Фильтры по цене',
    'filters.price.budget': 'До 20 000 ₽',
    'filters.price.mid': '20 000 - 40 000 ₽',
    'filters.price.premium': 'От 40 000 ₽',
    'filters.price.all': 'Все товары',
    
    // Корзина
    'cart.myCart': 'Моя корзина',
    'cart.empty': 'Ваша корзина пуста',
    'cart.total': 'Итого',
    'cart.checkout': 'Оформить заказ',
    'cart.remove': 'Удалить',
    'cart.unitPrice': 'Цена за шт.',
    'cart.items': 'товары',
    'cart.freeShipping': 'Бесплатно',
    
    // Заказы
    'orders.history': 'История заказов',
    'orders.viewHistory': 'Здесь вы можете просмотреть историю ваших заказов',
    'orders.noOrders': 'Заказов пока нет',
    'orders.firstOrder': 'Совершите свой первый заказ, и он появится здесь',
    'orders.goShopping': 'Перейти к покупкам',
    'orders.order': 'Заказ',
    'orders.status': 'Статус',
    'orders.date': 'Дата',
    'orders.total': 'Сумма',
    'orders.items': 'Товары',
    'orders.cancel': 'Отменить заказ',
    'orders.cancelConfirm': 'Вы уверены, что хотите отменить заказ?',
    'orders.status.pending': 'Ожидает обработки',
    'orders.status.processing': 'В обработке',
    'orders.status.shipped': 'Отправлен',
    'orders.status.delivered': 'Доставлен',
    'orders.status.cancelled': 'Отменен',
    
    // Футер
    'footer.description': 'Лучший выбор смартфонов по доступным ценам',
    'footer.contacts': 'Контакты',
    'footer.info': 'Информация',
    'footer.social': 'Мы в соцсетях',
    'footer.about': 'О компании',
    'footer.delivery': 'Доставка и оплата',
    'footer.guarantee': 'Гарантия',
    'footer.privacy': 'Политика конфиденциальности',
    'footer.rights': 'Все права защищены',
  },
  en: {
    // Common
    'app.title': 'Smartbuy',
    'loading': 'Loading...',
    'error': 'Error',
    'notFound': 'Not found',
    'save': 'Save',
    'cancel': 'Cancel',
    'delete': 'Delete',
    'edit': 'Edit',
    'yes': 'Yes',
    'no': 'No',
    
    // Navigation
    'nav.products': 'Products',
    'nav.cart': 'Cart',
    'nav.orders': 'Orders',
    'nav.users': 'Users',
    'nav.login': 'Login',
    'nav.signup': 'Sign Up',
    'nav.logout': 'Logout',
    'nav.continueShopping': 'Continue Shopping',
    'nav.backToShopping': '← Back to Shopping',
    
    // Auth
    'auth.login': 'Login',
    'auth.register': 'Register',
    'auth.name': 'Name',
    'auth.password': 'Password',
    'auth.confirmPassword': 'Confirm Password',
    'auth.noAccount': 'No account?',
    'auth.haveAccount': 'Already have an account?',
    'auth.email': 'Email',
    'auth.enterEmail': 'Enter email',
    'auth.enterPassword': 'Enter password',
    
    // Products
    'products.ourSmartphones': 'Our Smartphones',
    'products.memory': 'Memory',
    'products.ram': 'RAM',
    'products.screen': 'Screen',
    'products.price': 'Price',
    'products.rating': 'Rating',
    'products.noRatings': 'No ratings',
    'products.addToCart': 'Add to Cart',
    'products.inCart': 'In Cart',
    'products.popular': 'Popular Products',
    'products.showPopular': 'Show Popular',
    'products.showAll': 'Show All',
    'products.found': 'Products found',
    'products.notFound': 'Products not found',
    'products.changeFilters': 'Try changing filter parameters',
    'products.resetFilters': 'Show all products',
    'products.producer': 'Producer',
    'products.specs': 'Specifications',
    'products.description': 'Description',
    'products.addRewiew': 'Leave your review',
    'products.comment': 'Comment',
    'products.create': 'Create',
    'products.update': 'Update',
    'products.youReview': 'Your review',
    'products.reviews': 'Reviews',

    
    // Filters
    'filters.title': 'Product Catalog',
    'filters.active': 'Active filters',
    'filters.clearAll': 'Clear all',
    'filters.price': 'Price',
    'filters.price.budget': 'Up to 20,000 ₽',
    'filters.price.mid': '20,000 - 40,000 ₽',
    'filters.price.premium': 'From 40,000 ₽',
    'filters.price.all': 'All products',
    
    // Cart
    'cart.myCart': 'My Cart',
    'cart.empty': 'Your cart is empty',
    'cart.total': 'Total',
    'cart.checkout': 'Place an order',
    'cart.remove': 'Remove',
    'cart.unitPrice': 'Unit price',
    'cart.items': 'items',
    'cart.freeShipping': 'Free',
    
    // Orders
    'orders.history': 'Order History',
    'orders.viewHistory': 'Here you can view your order history',
    'orders.noOrders': 'No orders yet',
    'orders.firstOrder': 'Make your first order and it will appear here',
    'orders.goShopping': 'Go to shopping',
    'orders.order': 'Order',
    'orders.status': 'Status',
    'orders.date': 'Date',
    'orders.total': 'Total',
    'orders.items': 'Items',
    'orders.cancel': 'Cancel order',
    'orders.cancelConfirm': 'Are you sure you want to cancel the order?',
    'orders.status.pending': 'Pending',
    'orders.status.processing': 'Processing',
    'orders.status.shipped': 'Shipped',
    'orders.status.delivered': 'Delivered',
    'orders.status.cancelled': 'Cancelled',
    
    // Footer
    'footer.description': 'The best choice of smartphones at affordable prices',
    'footer.contacts': 'Contacts',
    'footer.info': 'Information',
    'footer.social': 'Follow us',
    'footer.about': 'About company',
    'footer.delivery': 'Delivery and payment',
    'footer.guarantee': 'Guarantee',
    'footer.privacy': 'Privacy policy',
    'footer.rights': 'All rights reserved',
  }
};

export const LanguageProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [language, setLanguage] = useState<Language>(() => {
    const savedLanguage = localStorage.getItem('language');
    return (savedLanguage as Language) || 'ru';
  });

  useEffect(() => {
    localStorage.setItem('language', language);
  }, [language]);

  const t = (key: string): string => {
    return translations[language][key as keyof typeof translations[typeof language]] || key;
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
};

export const useLanguage = () => {
  const context = useContext(LanguageContext);
  if (context === undefined) {
    throw new Error('useLanguage must be used within a LanguageProvider');
  }
  return context;
};