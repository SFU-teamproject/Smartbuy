// src/types.ts
export interface Smartphone {
  id: number;
  model: string;
  producer: string;
  memory: number;
  ram: number;
  display_size: number;
  price: number;
  ratings_sum: number;
  ratings_count: number;
  image_path?: string;
  description?: string;
}

export interface ApiError {
  message: string;
  status?: number;
}