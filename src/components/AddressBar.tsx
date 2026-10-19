import React, { useState, useEffect } from 'react';
import { WinRarIcons } from '../utils/icons';
import { useArchiveStore } from '../store/useArchiveStore';
import { ChevronRight, HardDrive, RefreshCw } from 'lucide-react';

export const AddressBar: React.FC = () => {
  const { currentPath, navigateTo, navigateUp } = useArchiveStore();
  const [isEditing, setIsEditing] = useState(false);
  const [inputVal, setInputVal] = useState(currentPath);

  useEffect(() => {
    setInputVal(currentPath);
  }, [currentPath]);

  // Clickable path segments
  const segments = currentPath.split('/');

  const handleSegmentClick = (index: number) => {
    const targetPath = segments.slice(0, index + 1).join('/');
    navigateTo(targetPath);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      navigateTo(inputVal.trim());
      setIsEditing(false);
    } else if (e.key === 'Escape') {
      setInputVal(currentPath);
      setIsEditing(false);
    }
  };

  return (
    <div className="bg-[#f0f0f0] border-b border-[#d8d8d8] px-2 py-1 flex items-center gap-1.5 text-xs">
      {/* Up One Level button */}
      <button
        onClick={navigateUp}
        className="winrar-btn flex items-center gap-1 px-2 py-0.5 text-[#222] hover:text-blue-700 shadow-xs"
        title="Up One Level (Backspace)"
      >
        <WinRarIcons.UpFolder size={16} />
        <span className="font-semibold text-[11px]">Up</span>
      </button>

      {/* Address path box */}
      <div className="flex-1 flex items-center bg-white border border-[#7f9db9] rounded-xs px-2 py-0.5 shadow-inner">
        <HardDrive size={13} className="text-blue-600 mr-1.5 shrink-0" />
        
        {isEditing ? (
          <input
            type="text"
            className="w-full text-xs font-mono outline-hidden bg-transparent text-gray-800"
            value={inputVal}
            onChange={(e) => setInputVal(e.target.value)}
            onKeyDown={handleKeyDown}
            onBlur={() => setIsEditing(false)}
            autoFocus
          />
        ) : (
          <div
            className="flex-1 flex items-center gap-0.5 overflow-x-auto cursor-text py-0.5"
            onClick={() => setIsEditing(true)}
          >
            {segments.map((seg, idx) => (
              <React.Fragment key={idx}>
                <span
                  onClick={(e) => {
                    e.stopPropagation();
                    handleSegmentClick(idx);
                  }}
                  className="hover:underline hover:text-blue-700 font-mono text-[11px] text-gray-800 cursor-pointer px-0.5 rounded-xs hover:bg-blue-50"
                >
                  {seg}
                </span>
                {idx < segments.length - 1 && (
                  <ChevronRight size={11} className="text-gray-400 shrink-0" />
                )}
              </React.Fragment>
            ))}
          </div>
        )}

        <button
          onClick={() => navigateTo(currentPath)}
          className="text-gray-400 hover:text-blue-600 ml-1 p-0.5"
          title="Refresh"
        >
          <RefreshCw size={11} />
        </button>
      </div>

      {/* Quick Drive Selector */}
      <select
        value={segments[0]}
        onChange={(e) => navigateTo(e.target.value)}
        className="bg-white border border-[#7f9db9] rounded-xs text-[11px] px-1 py-0.5 text-gray-700 outline-hidden"
      >
        <option value="project_backup.rar">Archive: project_backup.rar</option>
        <option value="C:">Local Disk (C:)</option>
        <option value="D:">Data Disk (D:)</option>
      </select>
    </div>
  );
};
