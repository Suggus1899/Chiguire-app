'use client';
import { ReactNode } from 'react';

export function Card({
  title,
  children,
  className = '',
  actions,
}: {
  title?: string;
  children: ReactNode;
  className?: string;
  actions?: ReactNode;
}) {
  return (
    <section className={`bg-white rounded-xl border border-gray-200 p-6 ${className}`}>
      {(title || actions) && (
        <div className="flex justify-between items-center mb-4">
          {title && <h2 className="font-medium text-gray-900">{title}</h2>}
          {actions}
        </div>
      )}
      {children}
    </section>
  );
}
