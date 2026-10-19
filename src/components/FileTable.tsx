import React, { useMemo, useState } from 'react';
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  SortingState,
  flexRender,
  ColumnDef
} from '@tanstack/react-table';
import { ArchiveEntry } from '../types';
import { useArchiveStore } from '../store/useArchiveStore';
import { WinRarIcons } from '../utils/icons';
import { formatNumberWithCommas } from '../utils/formatters';
import { Folder, ArrowUpDown, Lock } from 'lucide-react';

export const FileTable: React.FC = () => {
  const {
    entries,
    selectedEntries,
    selectEntry,
    drillDown,
    navigateUp,
    currentPath,
    archiveName,
    openModal
  } = useArchiveStore();

  const [sorting, setSorting] = useState<SortingState>([]);
  const [contextMenu, setContextMenu] = useState<{ x: number; y: number; entry: ArchiveEntry } | null>(null);

  const isSubfolder = currentPath !== archiveName;

  const columns = useMemo<ColumnDef<ArchiveEntry>[]>(
    () => [
      {
        accessorKey: 'name',
        header: 'Name',
        cell: ({ row }) => {
          const entry = row.original;
          const isSelected = selectedEntries.has(entry.name);
          return (
            <div className="flex items-center gap-2 truncate">
              <input
                type="checkbox"
                checked={isSelected}
                onChange={() => selectEntry(entry.name, true)}
                onClick={(e) => e.stopPropagation()}
                className="w-3.5 h-3.5 text-blue-600 rounded-xs border-gray-300"
              />
              {entry.is_dir ? (
                <Folder size={16} className="text-[#f57f17] fill-[#ffd54f] shrink-0" />
              ) : entry.name.endsWith('.rar') || entry.name.endsWith('.zip') ? (
                <WinRarIcons.ArchiveFile size={16} />
              ) : (
                <WinRarIcons.File size={16} />
              )}
              <span className="truncate text-xs font-medium text-gray-900">
                {entry.name}
              </span>
              {entry.is_encrypted && (
                <span title="Password Encrypted">
                  <Lock size={12} className="text-amber-600 shrink-0" />
                </span>
              )}
            </div>
          );
        }
      },
      {
        accessorKey: 'original_size',
        header: 'Original Size',
        cell: ({ getValue, row }) => {
          if (row.original.is_dir) return null;
          const val = getValue<number>();
          return (
            <span className="font-mono text-[11px] text-gray-700 block text-right pr-2">
              {formatNumberWithCommas(val)}
            </span>
          );
        }
      },
      {
        accessorKey: 'packed_size',
        header: 'Packed Size',
        cell: ({ getValue, row }) => {
          if (row.original.is_dir) return null;
          const val = getValue<number>();
          return (
            <span className="font-mono text-[11px] text-gray-600 block text-right pr-2">
              {formatNumberWithCommas(val)}
            </span>
          );
        }
      },
      {
        accessorKey: 'type',
        header: 'Type',
        cell: ({ getValue, row }) => {
          const val = getValue<string>() || (row.original.is_dir ? 'Folder' : 'File');
          return <span className="text-xs text-gray-600 truncate">{val}</span>;
        }
      },
      {
        accessorKey: 'modified',
        header: 'Modified',
        cell: ({ getValue }) => (
          <span className="font-mono text-[11px] text-gray-600 whitespace-nowrap">
            {getValue<string>()}
          </span>
        )
      },
      {
        accessorKey: 'attributes',
        header: 'Attr',
        cell: ({ getValue }) => (
          <span className="font-mono text-[11px] text-gray-500 uppercase">
            {getValue<string>() || 'A'}
          </span>
        )
      },
      {
        accessorKey: 'crc',
        header: 'CRC32',
        cell: ({ getValue, row }) => {
          if (row.original.is_dir) return null;
          return (
            <span className="font-mono text-[11px] text-gray-700 tracking-wider">
              {getValue<string>()}
            </span>
          );
        }
      }
    ],
    [selectedEntries, selectEntry]
  );

  const table = useReactTable({
    data: entries,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel()
  });

  const handleContextMenu = (e: React.MouseEvent, entry: ArchiveEntry) => {
    e.preventDefault();
    selectEntry(entry.name, false);
    setContextMenu({ x: e.clientX, y: e.clientY, entry });
  };

  return (
    <div
      className="flex-1 bg-white flex flex-col h-full overflow-hidden select-none relative"
      onClick={() => setContextMenu(null)}
    >
      <div className="flex-1 overflow-auto">
        <table className="w-full border-collapse text-left text-xs table-fixed">
          <thead className="bg-[#f0f0f0] sticky top-0 border-b border-[#c0c0c0] z-10 shadow-xs">
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  const isSorted = header.column.getIsSorted();
                  return (
                    <th
                      key={header.id}
                      onClick={header.column.getToggleSortingHandler()}
                      className="px-2.5 py-1.5 font-semibold text-gray-700 border-r border-[#d4d4d4] cursor-pointer hover:bg-[#e4e4e4] transition-colors overflow-hidden"
                      style={{
                        width:
                          header.id === 'name'
                            ? '35%'
                            : header.id === 'original_size' || header.id === 'packed_size'
                            ? '13%'
                            : header.id === 'type'
                            ? '12%'
                            : header.id === 'modified'
                            ? '16%'
                            : header.id === 'crc'
                            ? '10%'
                            : '6%'
                      }}
                    >
                      <div className="flex items-center justify-between">
                        <span className="truncate">
                          {flexRender(header.column.columnDef.header, header.getContext())}
                        </span>
                        {isSorted && (
                          <span className="text-blue-600 text-[10px] ml-1">
                            {isSorted === 'asc' ? '▲' : '▼'}
                          </span>
                        )}
                      </div>
                    </th>
                  );
                })}
              </tr>
            ))}
          </thead>
          <tbody>
            {/* Up one level row if inside a subfolder */}
            {isSubfolder && (
              <tr
                onDoubleClick={navigateUp}
                className="hover:bg-[#eaf4fe] cursor-pointer border-b border-gray-100 h-6"
              >
                <td className="px-2.5 py-1 flex items-center gap-2" colSpan={7}>
                  <WinRarIcons.UpFolder size={14} />
                  <span className="font-mono font-bold text-gray-700 text-xs">..</span>
                  <span className="text-[11px] text-gray-400 italic ml-2">[Up One Level]</span>
                </td>
              </tr>
            )}

            {table.getRowModel().rows.map((row) => {
              const entry = row.original;
              const isSelected = selectedEntries.has(entry.name);
              return (
                <tr
                  key={row.id}
                  onClick={(e) => selectEntry(entry.name, e.ctrlKey || e.metaKey, e.shiftKey)}
                  onDoubleClick={() => drillDown(entry)}
                  onContextMenu={(e) => handleContextMenu(e, entry)}
                  className={`border-b border-gray-100 cursor-pointer transition-colors h-6.5 ${
                    isSelected
                      ? 'bg-[#cce8ff] text-black border-blue-200'
                      : 'hover:bg-[#f2f8fe] text-[#222]'
                  }`}
                >
                  {row.getVisibleCells().map((cell) => (
                    <td
                      key={cell.id}
                      className="px-2.5 py-1 truncate border-r border-transparent"
                    >
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Context Menu */}
      {contextMenu && (
        <div
          className="fixed bg-white border border-gray-400 shadow-xl rounded-xs py-1 z-50 text-xs min-w-44 text-gray-800"
          style={{ top: contextMenu.y, left: contextMenu.x }}
          onClick={(e) => e.stopPropagation()}
        >
          <div className="px-3 py-1 font-semibold border-b border-gray-200 bg-gray-50 text-gray-600 truncate">
            {contextMenu.entry.name}
          </div>
          <button
            onClick={() => {
              setContextMenu(null);
              drillDown(contextMenu.entry);
            }}
            className="w-full text-left px-3 py-1.5 hover:bg-blue-600 hover:text-white"
          >
            {contextMenu.entry.is_dir ? 'Open folder' : 'View file (Alt+V)'}
          </button>
          <button
            onClick={() => {
              setContextMenu(null);
              openModal('extract');
            }}
            className="w-full text-left px-3 py-1.5 hover:bg-blue-600 hover:text-white"
          >
            Extract to... (Alt+E)
          </button>
          <button
            onClick={() => {
              setContextMenu(null);
              openModal('test');
            }}
            className="w-full text-left px-3 py-1.5 hover:bg-blue-600 hover:text-white"
          >
            Test archive (Alt+T)
          </button>
          <div className="my-1 border-t border-gray-200" />
          <button
            onClick={() => {
              setContextMenu(null);
              useArchiveStore.getState().deleteSelected();
            }}
            className="w-full text-left px-3 py-1.5 hover:bg-red-600 hover:text-white text-red-600"
          >
            Delete (Del)
          </button>
        </div>
      )}
    </div>
  );
};
