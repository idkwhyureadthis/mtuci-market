import { Building2, MapPin, User2Icon } from 'lucide-react';
import { Product } from '../types';
import { ImageGallery } from './ImageGallery';

interface ProductCardProps {
  product: Product;
  place: string;
}

export function ProductCard({ product, place }: ProductCardProps) {
  return (
    <div className={"bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow resize-none h-96 "} id={product.id}>
      <ImageGallery photos={product.photos} title={product.name} />
      <div className="p-4">
        <h3 className="text-lg font-semibold mb-2">{product.name}</h3>
        <p className="text-gray-600 text-sm mb-3 line-clamp-2">{product.description}</p>
        <div className="flex items-center justify-between mb-3">
          <span className="text-xl font-bold text-blue-600">{(product.price)} ₽</span>
          <div className="flex items-center text-gray-500 text-sm">
            <User2Icon size={16}className="mr-"/><p className='mr-2'>{product.creator_name}</p>
            <Building2 size={16} className="mr-1" />
            <span>№{product.dorm_number}</span>
            <MapPin size={16} className="ml-2 mr-1" />
            <span>к. {product.room}</span>
          </div>
        </div>
        {place == "mainpage" ? <a href={'https://t.me/' + product.telegram}><button className="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition-colors">
          Связаться с продавцом
        </button>
        </a>
        : ""
        } 
      </div>
    </div>
  );
}