export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  dorm_number: DormitoryNumber;
  photos: string[];
  creator_name: string;
  room: string;
  telegram: string;
  status: string
}

export type DormitoryNumber = "1" | "3" | "4";

export interface User {
  id: string;
  username: string;
  name: string;
  room: string;
  dormitory: DormitoryNumber;
  telegram: string;
  role: 'user' | 'moderator';
}