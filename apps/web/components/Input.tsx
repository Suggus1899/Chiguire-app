import { InputHTMLAttributes } from 'react';

export function Input({
  label,
  className = '',
  ...props
}: InputHTMLAttributes<HTMLInputElement> & { label?: string }) {
  return (
    <label className="block space-y-1">
      {label && <span className="block text-sm font-medium text-gray-700">{label}</span>}
      <input
        className={`w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${className}`}
        {...props}
      />
    </label>
  );
}
