'use client';
import { SelectHTMLAttributes } from 'react';

export function Select({
  label,
  className = '',
  children,
  ...props
}: SelectHTMLAttributes<HTMLSelectElement> & { label?: string }) {
  return (
    <label className="block space-y-1">
      {label && <span className="block text-sm font-medium text-gray-700">{label}</span>}
      <select
        className={`w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white ${className}`}
        {...props}
      >
        {children}
      </select>
    </label>
  );
}
