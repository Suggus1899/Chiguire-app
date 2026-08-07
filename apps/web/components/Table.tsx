import { ReactNode } from 'react';

type Column<T> = {
  key: string;
  label: string;
  render?: (row: T) => ReactNode;
};

export function Table<T extends { id: string }>({
  columns,
  data,
  empty = 'Sin datos',
}: {
  columns: Column<T>[];
  data: T[];
  empty?: string;
}) {
  if (data.length === 0) {
    return <p className="text-sm text-gray-400 py-4 text-center">{empty}</p>;
  }
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-200 text-left text-gray-500">
            {columns.map((col) => (
              <th key={col.key} className="py-2 px-3 font-medium">{col.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row) => (
            <tr key={row.id} className="border-b border-gray-100 hover:bg-gray-50">
              {columns.map((col) => (
                <td key={col.key} className="py-2 px-3 text-gray-800">
                  {col.render ? col.render(row) : (row as Record<string, ReactNode>)[col.key]}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
