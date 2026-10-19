import React, { useState, useRef, useEffect } from 'react';
import { useArchiveStore } from '../store/useArchiveStore';

interface MenuItem {
  label: string;
  shortcut?: string;
  action?: () => void;
  divider?: boolean;
  checked?: boolean;
}

export const MenuBar: React.FC = () => {
  const [openMenu, setOpenMenu] = useState<string | null>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  const {
    openModal,
    selectAll,
    clearSelection,
    deleteSelected,
    toolbarSize,
    toggleToolbarSize,
    isTreeCollapsed,
    toggleTree,
    navigateUp
  } = useArchiveStore();

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setOpenMenu(null);
      }
    };
    window.addEventListener('mousedown', handleClickOutside);
    return () => window.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Keyboard shortcut listener
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.altKey && e.key.toLowerCase() === 'a') {
        e.preventDefault();
        openModal('add');
      } else if (e.altKey && e.key.toLowerCase() === 'e') {
        e.preventDefault();
        openModal('extract');
      } else if (e.altKey && e.key.toLowerCase() === 't') {
        e.preventDefault();
        openModal('test');
      } else if (e.altKey && e.key.toLowerCase() === 'v') {
        e.preventDefault();
        openModal('view');
      } else if (e.altKey && e.key.toLowerCase() === 'r') {
        e.preventDefault();
        openModal('repair');
      } else if (e.altKey && e.key.toLowerCase() === 'i') {
        e.preventDefault();
        openModal('info');
      } else if (e.altKey && e.key.toLowerCase() === 'b') {
        e.preventDefault();
        openModal('benchmark');
      } else if (e.key === 'Delete') {
        deleteSelected();
      } else if (e.key === 'Backspace') {
        navigateUp();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [openModal, deleteSelected, navigateUp]);

  const menus: Record<string, MenuItem[]> = {
    File: [
      {
        label: 'Open archive...',
        shortcut: 'Ctrl+O',
        action: () => {
          const input = document.createElement('input');
          input.type = 'file';
          input.accept = '.rar,.zip,.7z,.tar,.gz,.bz2';
          input.onchange = (e: any) => {
            const file = e.target.files?.[0];
            if (file) {
              useArchiveStore.getState().openCustomArchive(file.name);
            }
          };
          input.click();
        }
      },
      { label: 'Save archive copy as...', shortcut: 'Ctrl+S' },
      { divider: true, label: '' },
      { label: 'Select all', shortcut: 'Ctrl+A', action: selectAll },
      { label: 'Clear selection', action: clearSelection },
      { divider: true, label: '' },
      {
        label: 'Exit',
        shortcut: 'Alt+F4',
        action: () => alert('Exit archiver application.')
      }
    ],
    Commands: [
      {
        label: 'Add files to archive',
        shortcut: 'Alt+A',
        action: () => openModal('add')
      },
      {
        label: 'Extract to specified folder',
        shortcut: 'Alt+E',
        action: () => openModal('extract')
      },
      {
        label: 'Test archive files',
        shortcut: 'Alt+T',
        action: () => openModal('test')
      },
      {
        label: 'View file',
        shortcut: 'Alt+V',
        action: () => openModal('view')
      },
      {
        label: 'Delete files',
        shortcut: 'Del',
        action: deleteSelected
      },
      { divider: true, label: '' },
      { label: 'Rename', shortcut: 'F2' },
      { label: 'Print file' },
      { divider: true, label: '' },
      { label: 'Protect archive from damage (Recovery Record)' },
      { label: 'Lock archive' }
    ],
    Tools: [
      { label: 'Wizard', action: () => openModal('wizard') },
      {
        label: 'Show information',
        shortcut: 'Alt+I',
        action: () => openModal('info')
      },
      { label: 'Convert archives' },
      {
        label: 'Repair archive',
        shortcut: 'Alt+R',
        action: () => openModal('repair')
      },
      {
        label: 'Benchmark and hardware test',
        shortcut: 'Alt+B',
        action: () => openModal('benchmark')
      },
      {
        label: 'Create SFX archive (.exe)',
        action: () => openModal('sfx')
      },
      { divider: true, label: '' },
      { label: 'Generate technical report' }
    ],
    Favorites: [
      { label: 'Add current folder to favorites' },
      { label: 'Organize favorites...' },
      { divider: true, label: '' },
      { label: 'D:\\Projects\\extracted' },
      { label: 'Desktop\\VibeRAR Archives' }
    ],
    Options: [
      {
        label: 'Large buttons (48x36)',
        checked: toolbarSize === 'large',
        action: () => {
          if (toolbarSize !== 'large') toggleToolbarSize();
        }
      },
      {
        label: 'Small buttons (24x24)',
        checked: toolbarSize === 'small',
        action: () => {
          if (toolbarSize !== 'small') toggleToolbarSize();
        }
      },
      { divider: true, label: '' },
      {
        label: 'Folder tree pane',
        checked: !isTreeCollapsed,
        action: toggleTree
      },
      { divider: true, label: '' },
      { label: 'Settings / Configuration...' },
      { label: 'View extraction log' }
    ],
    Help: [
      {
        label: 'Check for updates...',
        action: () => openModal('update')
      },
      {
        label: 'Repair Windows 11 Context Menu...',
        action: () => openModal('update')
      },
      { divider: true, label: '' },
      {
        label: 'VibeRAR Homepage (ahmedsamir.com)',
        action: () => window.open('https://ahmedsamir.com', '_blank')
      },
      {
        label: 'GitHub Repository (Releases & Issues)',
        action: () => window.open('https://github.com/iegyi/VibeRAR', '_blank')
      },
      {
        label: 'Buy Cheap Hosting (0webhost.com)',
        action: () => window.open('https://www.0webhost.com/', '_blank')
      },
      { divider: true, label: '' },
      {
        label: 'Architecture Report (PRD)',
        action: () => openModal('about')
      },
      { label: 'VibeRAR Help Topics', shortcut: 'F1' },
      { divider: true, label: '' },
      { label: 'About VibeRAR', action: () => openModal('about') }
    ]
  };

  return (
    <div
      ref={menuRef}
      className="bg-[#f0f0f0] border-b border-[#d0d0d0] text-xs flex items-center px-1 select-none z-40"
    >
      {Object.keys(menus).map(menuKey => {
        const isOpen = openMenu === menuKey;
        return (
          <div key={menuKey} className="relative">
            <button
              onClick={() => setOpenMenu(isOpen ? null : menuKey)}
              onMouseEnter={() => {
                if (openMenu) setOpenMenu(menuKey);
              }}
              className={`px-2.5 py-1 rounded-sm text-[#1e1e1e] cursor-default transition-colors ${
                isOpen ? 'bg-[#3399ff] text-white' : 'hover:bg-[#e5f1fb]'
              }`}
            >
              <span className="underline decoration-black/40">
                {menuKey[0]}
              </span>
              {menuKey.slice(1)}
            </button>

            {isOpen && (
              <div className="absolute left-0 top-full mt-0.5 min-w-[220px] bg-white border border-[#979797] shadow-lg rounded-xs py-1 z-50 text-[#1e1e1e]">
                {menus[menuKey].map((item, idx) => {
                  if (item.divider) {
                    return (
                      <div
                        key={idx}
                        className="my-1 border-t border-[#e0e0e0] mx-1"
                      />
                    );
                  }
                  return (
                    <button
                      key={idx}
                      onClick={() => {
                        setOpenMenu(null);
                        item.action?.();
                      }}
                      className="w-full text-left px-4 py-1.5 hover:bg-[#2b88d8] hover:text-white flex items-center justify-between text-xs transition-colors"
                    >
                      <div className="flex items-center gap-2">
                        {item.checked !== undefined ? (
                          <span className="w-3 text-center font-bold">
                            {item.checked ? '✓' : ''}
                          </span>
                        ) : (
                          <span className="w-3"></span>
                        )}
                        <span>{item.label}</span>
                      </div>
                      {item.shortcut && (
                        <span className="text-gray-400 group-hover:text-white ml-4 font-mono text-[11px]">
                          {item.shortcut}
                        </span>
                      )}
                    </button>
                  );
                })}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
};
