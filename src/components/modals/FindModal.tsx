import React, { useState } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { X, Search } from 'lucide-react';

export const FindModal: React.FC = () => {
  const { activeModal, closeModal, entries, selectEntry } = useArchiveStore();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<any[]>([]);
  const [searched, setSearched] = useState(false);

  if (activeModal !== 'find') return null;

  const handleSearch = () => {
    setSearched(true);
    const q = query.toLowerCase().trim();
    if (!q) {
      setResults(entries);
      return;
    }
    const matched = entries.filter(
      (e) =>
        e.name.toLowerCase().includes(q) ||
        (e.content && e.content.toLowerCase().includes(q))
    );
    setResults(matched);
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[500px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>Find files</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Content Body */}
        <div className="p-3 bg-white flex flex-col gap-3 min-h-[300px]">
          <div>
            <label className="text-[11px] text-gray-600 block mb-1">File name to find:</label>
            <div className="flex gap-2">
              <input
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="*.* or file name pattern"
                className="flex-1 border border-[#7f9db9] rounded-xs px-2 py-1 font-mono text-xs text-gray-900"
                autoFocus
              />
              <button
                onClick={handleSearch}
                className="winrar-btn px-4 py-1 text-xs font-semibold text-gray-800 flex items-center gap-1"
              >
                <Search size={13} />
                <span>Find</span>
              </button>
            </div>
          </div>

          <fieldset className="border border-[#d0d0d0] p-2 rounded-xs">
            <legend className="px-1 text-[11px] font-semibold text-gray-700">Search scope</legend>
            <div className="flex items-center gap-4 text-xs">
              <label className="flex items-center gap-1.5 cursor-pointer">
                <input type="radio" name="scope" defaultChecked className="text-blue-600" />
                <span>Current archive</span>
              </label>
              <label className="flex items-center gap-1.5 cursor-pointer">
                <input type="radio" name="scope" className="text-blue-600" />
                <span>All archives on disk</span>
              </label>
            </div>
          </fieldset>

          {/* Results Table */}
          <div className="border border-gray-300 rounded-xs flex-1 h-44 overflow-y-auto">
            {searched && (
              <table className="w-full text-xs text-left">
                <thead className="bg-gray-100 border-b border-gray-200 sticky top-0 font-semibold text-gray-600">
                  <tr>
                    <th className="px-2 py-1">Name</th>
                    <th className="px-2 py-1">Path</th>
                    <th className="px-2 py-1 text-right">Size</th>
                  </tr>
                </thead>
                <tbody>
                  {results.length === 0 ? (
                    <tr>
                      <td colSpan={3} className="text-center py-6 text-gray-400">
                        No files matching "{query}"
                      </td>
                    </tr>
                  ) : (
                    results.map((r) => (
                      <tr
                        key={r.name}
                        onClick={() => {
                          selectEntry(r.name, false);
                          closeModal();
                        }}
                        className="hover:bg-blue-100 cursor-pointer border-b border-gray-100"
                      >
                        <td className="px-2 py-1 font-medium">{r.name}</td>
                        <td className="px-2 py-1 text-gray-500 font-mono text-[11px]">{r.path}</td>
                        <td className="px-2 py-1 text-right font-mono text-[11px]">
                          {r.original_size.toLocaleString()} B
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            )}
          </div>
        </div>

        {/* Action Buttons */}
        <div className="h-11 bg-[#f0f0f0] border-t border-[#dcdcdc] flex items-center justify-end px-3 gap-2">
          <button
            onClick={closeModal}
            className="winrar-btn w-20 py-1 text-xs text-gray-800"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};
