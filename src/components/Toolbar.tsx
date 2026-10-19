import React from 'react';
import { WinRarIcons } from '../utils/icons';
import { useArchiveStore } from '../store/useArchiveStore';
import { LayoutGrid, Maximize2 } from 'lucide-react';

export const Toolbar: React.FC = () => {
  const {
    toolbarSize,
    toggleToolbarSize,
    openModal,
    deleteSelected,
    selectedEntries
  } = useArchiveStore();

  const isLarge = toolbarSize === 'large';
  const iconSize = isLarge ? 36 : 22;

  const tools = [
    {
      id: 'add',
      name: 'Add',
      icon: <WinRarIcons.AddArchive size={iconSize} />,
      onClick: () => openModal('add'),
      title: 'Add files to archive (Alt+A)'
    },
    {
      id: 'extract',
      name: 'Extract To',
      icon: <WinRarIcons.ExtractTo size={iconSize} />,
      onClick: () => openModal('extract'),
      title: 'Extract files to specified folder (Alt+E)'
    },
    {
      id: 'test',
      name: 'Test',
      icon: <WinRarIcons.Test size={iconSize} />,
      onClick: () => openModal('test'),
      title: 'Test archive integrity and checksums (Alt+T)'
    },
    {
      id: 'view',
      name: 'View',
      icon: <WinRarIcons.View size={iconSize} />,
      onClick: () => openModal('view'),
      title: 'View internal file content (Alt+V)'
    },
    {
      id: 'delete',
      name: 'Delete',
      icon: <WinRarIcons.Delete size={iconSize} />,
      onClick: deleteSelected,
      title: 'Delete selected files from archive (Del)',
      disabled: selectedEntries.size === 0
    },
    {
      id: 'find',
      name: 'Find',
      icon: <WinRarIcons.Find size={iconSize} />,
      onClick: () => openModal('find'),
      title: 'Find files inside archive or folders'
    },
    {
      id: 'wizard',
      name: 'Wizard',
      icon: <WinRarIcons.Wizard size={iconSize} />,
      onClick: () => openModal('wizard'),
      title: 'Easy step-by-step extraction and compression wizard'
    },
    {
      id: 'info',
      name: 'Info',
      icon: <WinRarIcons.Info size={iconSize} />,
      onClick: () => openModal('info'),
      title: 'Show detailed archive information and compression ratios (Alt+I)'
    },
    {
      id: 'repair',
      name: 'Repair',
      icon: <WinRarIcons.Repair size={iconSize} />,
      onClick: () => openModal('repair'),
      title: 'Repair damaged archive using Reed-Solomon parity records (Alt+R)'
    },
    {
      id: 'benchmark',
      name: 'Benchmark',
      icon: <WinRarIcons.Benchmark size={iconSize} />,
      onClick: () => openModal('benchmark'),
      title: 'Hardware compression throughput and CPU benchmark (Alt+B)'
    },
    {
      id: 'sfx',
      name: 'SFX',
      icon: <WinRarIcons.SFX size={iconSize} />,
      onClick: () => openModal('sfx'),
      title: 'Create self-extracting executable archive (.exe)'
    }
  ];

  return (
    <div className="bg-[#f2f2f2] border-b border-[#c8c8c8] px-2 py-1 flex items-center justify-between select-none shadow-[inset_0_1px_0_rgba(255,255,255,0.8)]">
      <div className="flex items-center gap-1 overflow-x-auto py-0.5">
        {tools.map((tool, idx) => (
          <React.Fragment key={tool.id}>
            {/* Visual separator before Info and Benchmark */}
            {(tool.id === 'find' || tool.id === 'benchmark') && (
              <div className="w-[1px] h-8 bg-[#c0c0c0] mx-1 self-center" />
            )}
            <button
              onClick={tool.onClick}
              disabled={tool.disabled}
              title={tool.title}
              className={`winrar-tool-btn flex flex-col items-center justify-center p-1 text-[#222] ${
                tool.disabled ? 'opacity-40 cursor-not-allowed' : 'cursor-pointer'
              } ${
                isLarge
                  ? 'w-[58px] min-w-[58px] h-[54px]'
                  : 'flex-row gap-1.5 px-2.5 h-[32px]'
              }`}
            >
              <div className="flex items-center justify-center">{tool.icon}</div>
              <span
                className={`font-sans tracking-tight text-[#222] ${
                  isLarge ? 'text-[11px] mt-0.5' : 'text-xs'
                }`}
              >
                {tool.name}
              </span>
            </button>
          </React.Fragment>
        ))}
      </div>

      {/* Density toggle button */}
      <div className="flex items-center gap-1 pl-2 border-l border-[#c8c8c8]">
        <button
          onClick={toggleToolbarSize}
          className="winrar-btn px-2 py-1 text-xs text-gray-700 flex items-center gap-1"
          title={`Switch to ${isLarge ? 'Small' : 'Large'} buttons`}
        >
          {isLarge ? (
            <>
              <Maximize2 size={13} className="text-gray-600" />
              <span className="text-[11px]">24×24</span>
            </>
          ) : (
            <>
              <LayoutGrid size={13} className="text-gray-600" />
              <span className="text-[11px]">48×36</span>
            </>
          )}
        </button>
      </div>
    </div>
  );
};
