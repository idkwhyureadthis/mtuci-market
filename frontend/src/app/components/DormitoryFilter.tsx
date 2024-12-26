import React from 'react';
import { Building2 } from 'lucide-react';
import { DormitoryNumber } from '../types';

interface DormitoryFilterProps {
  selectedDormitory: DormitoryNumber | null;
  onSelect: (dormitory: DormitoryNumber | null) => void;
}

export function DormitoryFilter({ selectedDormitory, onSelect }: DormitoryFilterProps) {
  const dormitories: DormitoryNumber[] = ["1", "3", "4"];

  return (
    <div className="flex gap-2 mb-6">
      <button
        onClick={() => onSelect(null)}
        className={`px-4 py-2 rounded-md flex items-center gap-2 ${
          selectedDormitory === null
            ? 'bg-blue-600 text-white'
            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
        }`}
      >
        <Building2 size={18} />
        Все общежития
      </button>
      {dormitories.map((dorm) => (
        <button
          key={dorm}
          onClick={() => onSelect(dorm)}
          className={`px-4 py-2 rounded-md ${
            selectedDormitory === dorm
              ? 'bg-blue-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
        >
          №{dorm}
        </button>
      ))}
    </div>
  );
}