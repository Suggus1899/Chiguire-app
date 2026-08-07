import Link from 'next/link';

export default function NotFound() {
  return (
    <div className="flex items-center justify-center min-h-screen px-4">
      <div className="text-center space-y-4 max-w-md">
        <h1 className="text-4xl font-bold text-gray-900">404</h1>
        <h2 className="text-xl font-semibold text-gray-700">Página no encontrada</h2>
        <p className="text-sm text-gray-500">
          La página que buscas no existe o fue movida.
        </p>
        <Link
          href="/"
          className="inline-block bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700"
        >
          Volver al inicio
        </Link>
      </div>
    </div>
  );
}
