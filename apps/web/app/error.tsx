'use client';

import { useEffect } from 'react';

export default function Error({
  error,
  retry,
}: {
  error: Error & { digest?: string };
  retry: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex items-center justify-center min-h-[50vh] px-4">
      <div className="text-center space-y-4 max-w-md">
        <h2 className="text-xl font-semibold text-gray-900">Algo salió mal</h2>
        <p className="text-sm text-gray-500">
          Ocurrió un error inesperado. Puedes intentar de nuevo.
        </p>
        <button
          onClick={retry}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700"
        >
          Reintentar
        </button>
      </div>
    </div>
  );
}
