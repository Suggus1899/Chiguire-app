'use client';
import { useToastStore, type ToastVariant } from '@/lib/toast';

const variantStyles: Record<ToastVariant, { bg: string; icon: string }> = {
  success: { bg: 'bg-green-600', icon: '✓' },
  error: { bg: 'bg-red-600', icon: '✕' },
  info: { bg: 'bg-blue-600', icon: 'ℹ' },
};

export function ToastContainer() {
  const toasts = useToastStore((s) => s.toasts);
  const remove = useToastStore((s) => s.remove);

  if (toasts.length === 0) return null;

  return (
    <div className="fixed bottom-4 right-4 z-50 space-y-2 max-w-sm">
      {toasts.map((t) => {
        const style = variantStyles[t.variant];
        return (
          <div
            key={t.id}
            className={`${style.bg} text-white rounded-lg shadow-lg px-4 py-3 flex items-start gap-3 animate-in`}
            role="alert"
          >
            <span className="font-bold flex-shrink-0">{style.icon}</span>
            <span className="text-sm flex-1">{t.message}</span>
            <button
              onClick={() => remove(t.id)}
              className="text-white/80 hover:text-white flex-shrink-0"
              aria-label="Cerrar"
            >
              ×
            </button>
          </div>
        );
      })}
    </div>
  );
}
