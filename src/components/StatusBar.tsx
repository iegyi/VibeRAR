import React from 'react';
import { useArchiveStore } from '../store/useArchiveStore';
import { formatNumberWithCommas } from '../utils/formatters';

export const StatusBar: React.FC = () => {
  const { entries, selectedEntries } = useArchiveStore();

  const totalFiles = entries.length;
  const totalBytes = entries.reduce((acc, curr) => acc + curr.original_size, 0);

  const selectedCount = selectedEntries.size;
  const selectedBytes = entries
    .filter((e) => selectedEntries.has(e.name))
    .reduce((acc, curr) => acc + curr.original_size, 0);

  return (
    <div className="h-6 bg-[#f0f0f0] border-t border-[#c0c0c0] flex items-center text-xs select-none px-1 text-gray-700">
      {/* Selected Items Panel */}
      <div className="flex-1 px-3 py-0.5 border-r border-[#d4d4d4] truncate">
        {selectedCount > 0 ? (
          <span>
            Selected <strong className="font-semibold">{selectedCount}</strong> {selectedCount === 1 ? 'file' : 'files'},{' '}
            <strong className="font-mono text-[11px]">{formatNumberWithCommas(selectedBytes)}</strong> bytes
          </span>
        ) : (
          <span className="text-gray-500">No selection</span>
        )}
      </div>

      {/* Total Items Panel */}
      <div className="px-4 py-0.5 border-r border-[#d4d4d4] shrink-0 text-right">
        <span>
          Total <strong className="font-semibold">{totalFiles}</strong> {totalFiles === 1 ? 'file' : 'files'},{' '}
          <strong className="font-mono text-[11px]">{formatNumberWithCommas(totalBytes)}</strong> bytes
        </span>
      </div>

      {/* Format & Solid indicator */}
      <div className="px-3 py-0.5 text-[11px] font-mono text-gray-600 shrink-0 bg-[#e8e8e8] border border-[#d0d0d0] rounded-xs ml-1">
        Solid RAR5 (Dict: 32MB)
      </div>
    </div>
  );
};
