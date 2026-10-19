import React from 'react';
import { Minus, Square, X } from 'lucide-react';
import { useArchiveStore } from '../store/useArchiveStore';

export const TitleBar: React.FC = () => {
  const currentPath = useArchiveStore(state => state.currentPath);
  const archiveName = useArchiveStore(state => state.archiveName);

  return (
    <div className="h-8 bg-[#1e293b] text-slate-100 flex items-center justify-between px-3 select-none text-xs border-b border-slate-700 shadow-sm">
      <div className="flex items-center gap-2 truncate">
        {/* VibeRAR 3D Icon */}
        <img src="/viberar_icon.png" alt="VibeRAR" className="w-4 h-4 object-contain rounded-xs shrink-0" />
        <span className="font-semibold text-slate-200 truncate">
          {currentPath ? `${currentPath.replace(/\//g, '\\')}` : archiveName} — VibeRAR
        </span>
      </div>

      <div className="flex items-center gap-1 -mr-2">
        <button
          className="w-10 h-8 hover:bg-slate-700/60 active:bg-slate-700 flex items-center justify-center text-slate-300 transition-colors"
          title="Minimize"
        >
          <Minus size={14} />
        </button>
        <button
          className="w-10 h-8 hover:bg-slate-700/60 active:bg-slate-700 flex items-center justify-center text-slate-300 transition-colors"
          title="Maximize"
        >
          <Square size={12} />
        </button>
        <button
          className="w-11 h-8 hover:bg-red-600 active:bg-red-700 hover:text-white flex items-center justify-center text-slate-300 transition-colors"
          title="Close"
          onClick={() => {
            if (window.confirm('Do you want to close VibeRAR?')) {
              window.close();
            }
          }}
        >
          <X size={15} />
        </button>
      </div>
    </div>
  );
};
