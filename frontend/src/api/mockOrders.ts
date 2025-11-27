import { Order, Smartphone } from '../types';

// Моковые смартфоны для заказов
const mockSmartphones: Smartphone[] = [
  {
    id: 1,
    model: "iPhone 15 Pro",
    producer: "Apple",
    memory: 256,
    ram: 8,
    display_size: 6.1,
    price: 99990,
    ratings_sum: 45,
    ratings_count: 9,
    image_path: "/images/iphone15pro.jpg",
    description: "Флагманский смартфон Apple"
  },
  {
    id: 2,
    model: "Galaxy S24",
    producer: "Samsung",
    memory: 256,
    ram: 12,
    display_size: 6.2,
    price: 79990,
    ratings_sum: 38,
    ratings_count: 8,
    image_path: "/images/galaxys24.jpg",
    description: "Флагманский смартфон Samsung"
  },
  {
    id: 3,
    model: "Redmi Note 13",
    producer: "Xiaomi",
    memory: 128,
    ram: 6,
    display_size: 6.5,
    price: 24990,
    ratings_sum: 28,
    ratings_count: 6,
    image_path: "/images/redminote13.jpg",
    description: "Популярный смартфон от Xiaomi"
  }
];

// Моковые заказы
export const mockOrders: Order[] = [
  {
    id: 1,
    user_id: 1,
    status: 'delivered',
    total_amount: 99990,
    created_at: '2024-01-15T10:30:00Z',
    updated_at: '2024-01-18T14:20:00Z',
    items: [
      {
        id: 1,
        order_id: 1,
        smartphone_id: 1,
        quantity: 1,
        price: 99990,
        smartphone: mockSmartphones[0]
      }
    ]
  },
  {
    id: 2,
    user_id: 1,
    status: 'shipped',
    total_amount: 159980,
    created_at: '2024-01-20T16:45:00Z',
    updated_at: '2024-01-21T09:15:00Z',
    items: [
      {
        id: 2,
        order_id: 2,
        smartphone_id: 2,
        quantity: 2,
        price: 79990,
        smartphone: mockSmartphones[1]
      }
    ]
  },
  {
    id: 3,
    user_id: 1,
    status: 'processing',
    total_amount: 24990,
    created_at: '2024-01-25T11:20:00Z',
    updated_at: '2024-01-25T11:20:00Z',
    items: [
      {
        id: 3,
        order_id: 3,
        smartphone_id: 3,
        quantity: 1,
        price: 24990,
        smartphone: mockSmartphones[2]
      }
    ]
  },
  {
    id: 4,
    user_id: 1,
    status: 'cancelled',
    total_amount: 79990,
    created_at: '2024-01-10T09:15:00Z',
    updated_at: '2024-01-11T10:30:00Z',
    items: [
      {
        id: 4,
        order_id: 4,
        smartphone_id: 2,
        quantity: 1,
        price: 79990,
        smartphone: mockSmartphones[1]
      }
    ]
  }
];