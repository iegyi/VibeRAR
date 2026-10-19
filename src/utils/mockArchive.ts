import { ArchiveEntry } from '../types';

export const INITIAL_ARCHIVE_NAME = 'project_backup.rar';

export const MOCK_ARCHIVE_DATA: Record<string, ArchiveEntry[]> = {
  'project_backup.rar': [
    {
      name: 'docs',
      path: 'project_backup.rar/docs',
      is_dir: true,
      original_size: 0,
      packed_size: 0,
      modified: '2026-09-14 09:15',
      crc: '00000000',
      is_encrypted: false,
      attributes: 'D',
      type: 'Folder'
    },
    {
      name: 'assets',
      path: 'project_backup.rar/assets',
      is_dir: true,
      original_size: 0,
      packed_size: 0,
      modified: '2026-09-15 11:20',
      crc: '00000000',
      is_encrypted: false,
      attributes: 'D',
      type: 'Folder'
    },
    {
      name: 'src',
      path: 'project_backup.rar/src',
      is_dir: true,
      original_size: 0,
      packed_size: 0,
      modified: '2026-09-16 14:00',
      crc: '00000000',
      is_encrypted: false,
      attributes: 'D',
      type: 'Folder'
    },
    {
      name: 'README.md',
      path: 'project_backup.rar/README.md',
      is_dir: false,
      original_size: 4520,
      packed_size: 1320,
      modified: '2026-09-16 16:30',
      crc: '78A4E90B',
      is_encrypted: false,
      attributes: 'A',
      type: 'Markdown File',
      content: '# VibeRAR Archiver\n\nHigh performance archive decompression and compression system built with native systems engine and React 19.\n\n### Features\n- Dual-density VibeRAR interface (Large & Small toolbars)\n- RAR4, RAR5, RAR7, ZIP, 7z decompression\n- ZIP64 and 7z compression with AES-256\n- Built-in multi-threaded hardware benchmark\n- Self-extracting archive (SFX) creation\n'
    },
    {
      name: 'license.txt',
      path: 'project_backup.rar/license.txt',
      is_dir: false,
      original_size: 2150,
      packed_size: 780,
      modified: '2026-09-10 12:00',
      crc: '4A12B7C8',
      is_encrypted: false,
      attributes: 'A',
      type: 'Text Document',
      content: 'VIBERAR SOFTWARE LICENSE\n\nThis archive manager is provided for high performance archive operations.\nZero file processing occurs on the UI thread.\nAll read/write operations stream through native channels.\n'
    }
  ],

  'project_backup.rar/src': [
    {
      name: 'components',
      path: 'project_backup.rar/src/components',
      is_dir: true,
      original_size: 0,
      packed_size: 0,
      modified: '2026-09-16 14:02',
      crc: '00000000',
      is_encrypted: false,
      attributes: 'D',
      type: 'Folder'
    },
    {
      name: 'utils',
      path: 'project_backup.rar/src/utils',
      is_dir: true,
      original_size: 0,
      packed_size: 0,
      modified: '2026-09-15 10:14',
      crc: '00000000',
      is_encrypted: false,
      attributes: 'D',
      type: 'Folder'
    },
    {
      name: 'tests',
      path: 'project_backup.rar/src/tests',
      is_dir: true,
      original_size: 0,
      packed_size: 0,
      modified: '2026-09-15 16:22',
      crc: '00000000',
      is_encrypted: false,
      attributes: 'D',
      type: 'Folder'
    },
    {
      name: 'App.tsx',
      path: 'project_backup.rar/src/App.tsx',
      is_dir: false,
      original_size: 15320,
      packed_size: 3840,
      modified: '2026-09-16 18:20',
      crc: '6BC42A90',
      is_encrypted: false,
      attributes: 'A',
      type: 'TSX File',
      content: 'import React from "react";\n\nexport function App() {\n  return <div className="winrar-desktop">VibeRAR</div>;\n}\n'
    },
    {
      name: 'main.tsx',
      path: 'project_backup.rar/src/main.tsx',
      is_dir: false,
      original_size: 1820,
      packed_size: 540,
      modified: '2026-09-14 08:30',
      crc: '8E44C102',
      is_encrypted: false,
      attributes: 'A',
      type: 'TSX File',
      content: 'import ReactDOM from "react-dom/client";\nimport { App } from "./App";\n\nReactDOM.createRoot(document.getElementById("root")!).render(<App />);\n'
    }
  ],

  // Figure 1 items matching the Architecture Specification Report exactly!
  'project_backup.rar/src/components': [
    {
      name: 'Toolbar.tsx',
      path: 'project_backup.rar/src/components/Toolbar.tsx',
      is_dir: false,
      original_size: 14208,
      packed_size: 3911,
      modified: '2026-09-14 10:22',
      crc: '9F3A12C7',
      is_encrypted: false,
      attributes: 'A',
      type: 'TSX File',
      content: '// WinRAR Toolbar with Large (48x36px) & Small (24x24px) Button Support\nexport const Toolbar = () => {\n  return <div className="toolbar">Action Buttons</div>;\n};\n'
    },
    {
      name: 'MenuBar.tsx',
      path: 'project_backup.rar/src/components/MenuBar.tsx',
      is_dir: false,
      original_size: 9874,
      packed_size: 2640,
      modified: '2026-09-14 10:25',
      crc: '1B07D4E2',
      is_encrypted: false,
      attributes: 'A',
      type: 'TSX File',
      content: '// Classic Menu Bar: File, Commands, Tools, Favorites, Options, Help\nexport const MenuBar = () => {\n  return <nav>WinRAR Menu</nav>;\n};\n'
    },
    {
      name: 'FileGrid.tsx',
      path: 'project_backup.rar/src/components/FileGrid.tsx',
      is_dir: false,
      original_size: 22516,
      packed_size: 5102,
      modified: '2026-09-15 08:03',
      crc: 'C45E9A10',
      is_encrypted: false,
      attributes: 'A',
      type: 'TSX File',
      content: '// Virtualized File Grid (TanStack Table) with Name, Size, Packed, Type, Date, Attr, CRC32\nexport const FileGrid = () => {\n  return <table>File Table</table>;\n};\n'
    },
    {
      name: 'StatusBar.tsx',
      path: 'project_backup.rar/src/components/StatusBar.tsx',
      is_dir: false,
      original_size: 3120,
      packed_size: 1018,
      modified: '2026-09-15 08:40',
      crc: '7D21FB66',
      is_encrypted: false,
      attributes: 'A',
      type: 'TSX File',
      content: '// Bottom Status Bar: dynamic items counter, volume count, selection\nexport const StatusBar = () => {\n  return <footer>Status</footer>;\n};\n'
    },
    {
      name: 'ExtractDialog.tsx',
      path: 'project_backup.rar/src/components/ExtractDialog.tsx',
      is_dir: false,
      original_size: 18430,
      packed_size: 4477,
      modified: '2026-09-16 17:11',
      crc: 'E00C3B95',
      is_encrypted: false,
      attributes: 'A',
      type: 'TSX File',
      content: '// Extraction Path and Options modal with General and Advanced tabs\nexport const ExtractDialog = () => {\n  return <div>Extract Options</div>;\n};\n'
    },
    {
      name: 'ProgressModal.tsx',
      path: 'project_backup.rar/src/components/ProgressModal.tsx',
      is_dir: false,
      original_size: 11902,
      packed_size: 3260,
      modified: '2026-09-16 17:48',
      crc: '34AA7FB1',
      is_encrypted: false,
      attributes: 'A',
      type: 'TSX File',
      content: '// Dual-progress bar with Current file progress, Total progress, Transfer rate\nexport const ProgressModal = () => {\n  return <div>Dual Progress</div>;\n};\n'
    },
    {
      name: 'styles.css',
      path: 'project_backup.rar/src/components/styles.css',
      is_dir: false,
      original_size: 6044,
      packed_size: 1390,
      modified: '2026-09-12 09:00',
      crc: 'B2C9D073',
      is_encrypted: false,
      attributes: 'A',
      type: 'CSS File',
      content: '/* Classic WinRAR UI bevels and color tokens */\n.winrar-panel {\n  border: 1px solid #7f9db9;\n}\n'
    }
  ],

  'project_backup.rar/docs': [
    {
      name: 'architecture_spec.pdf',
      path: 'project_backup.rar/docs/architecture_spec.pdf',
      is_dir: false,
      original_size: 769566,
      packed_size: 642100,
      modified: '2026-10-18 11:33',
      crc: 'E824A101',
      is_encrypted: false,
      attributes: 'A',
      type: 'PDF Document'
    },
    {
      name: 'changelog.txt',
      path: 'project_backup.rar/docs/changelog.txt',
      is_dir: false,
      original_size: 8940,
      packed_size: 2100,
      modified: '2026-09-16 15:40',
      crc: 'D452EF99',
      is_encrypted: false,
      attributes: 'A',
      type: 'Text Document',
      content: 'VibeRAR v1.0.0 Release Notes\n- Implemented dual-density toolbars (48x36 and 24x24)\n- Implemented full Extraction Path & Options dialog\n- Added multi-threaded hardware benchmark\n- Added Advanced SFX configuration engine\n- Added dual-progress streaming dialog\n'
    }
  ]
};
