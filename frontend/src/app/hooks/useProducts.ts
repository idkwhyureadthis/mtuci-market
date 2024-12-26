import { useState, useMemo, useCallback } from 'react';
import { Product, DormitoryNumber } from '../types';

const ITEMS_PER_PAGE = 20;

export function useProducts(products: Product[]) {
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedDormitory, setSelectedDormitory] = useState<DormitoryNumber | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  const handlePageChange = useCallback((page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, []);

  const filteredProducts = useMemo(() => {
    return products.filter((product) => {
      const matchesDormitory = selectedDormitory ? product.dorm_number === selectedDormitory : true;
      const matchesSearch = product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                           product.description.toLowerCase().includes(searchQuery.toLowerCase());
      return matchesDormitory && matchesSearch;
    });
  }, [products, selectedDormitory, searchQuery]);

  const totalPages = Math.ceil(filteredProducts.length / ITEMS_PER_PAGE);
  
  const paginatedProducts = useMemo(() => {
    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    return filteredProducts.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  }, [filteredProducts, currentPage]);

  return {
    products: paginatedProducts,
    currentPage,
    totalPages,
    setCurrentPage: handlePageChange,
    selectedDormitory,
    setSelectedDormitory,
    searchQuery,
    setSearchQuery,
  };
}