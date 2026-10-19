import React, { useEffect, useState } from 'react';
import { TitleBar } from './components/TitleBar';
import { MenuBar } from './components/MenuBar';
import { Toolbar } from './components/Toolbar';
import { AddressBar } from './components/AddressBar';
import { FolderTree } from './components/FolderTree';
import { FileTable } from './components/FileTable';
import { StatusBar } from './components/StatusBar';
import { ExtractModal } from './components/modals/ExtractModal';
import { AddModal } from './components/modals/AddModal';
import { SfxModal } from './components/modals/SfxModal';
import { BenchmarkModal } from './components/modals/BenchmarkModal';
import { ProgressModal } from './components/modals/ProgressModal';
import { InfoModal } from './components/modals/InfoModal';
import { ViewModal } from './components/modals/ViewModal';
import { RepairModal } from './components/modals/RepairModal';
import { FindModal } from './components/modals/FindModal';
import { WizardModal } from './components/modals/WizardModal';
import { AboutModal } from './components/modals/AboutModal';
import { useArchiveStore } from './store/useArchiveStore';

export const App: React.FC = () => {
  const {
    openCustomArchive,
    currentPath,
    archiveName,
    openModal,
    updateCompressionConfig,
    updateExtractionConfig
  } = useArchiveStore();
  const [isDragOver, setIsDragOver] = useState(false);

  // Synchronize native Windows title bar with current archive path
  useEffect(() => {
    const displayPath = currentPath ? currentPath.replace(/\//g, '\\') : archiveName;
    document.title = `${displayPath} — VibeRAR`;
  }, [currentPath, archiveName]);

  // Handle command line invocation from Explorer context menus (?action=add|extract|open&target=...)
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const action = params.get('action');
    const target = params.get('target');

    if (target) {
      const cleanTarget = target.replace(/\\/g, '/');
      const filename = cleanTarget.split('/').pop() || 'archive';
      const lastSlash = target.lastIndexOf('\\');
      const dir = lastSlash !== -1 ? target.substring(0, lastSlash) : '';

      if (action === 'add') {
        const baseName = filename.replace(/\.[^.]+$/, '');
        const proposedArchive = `${baseName}.rar`;
        updateCompressionConfig({
          sourcePath: target,
          archiveName: proposedArchive,
          destinationDir: dir,
          format: 'rar'
        });
        openModal('add');
      } else if (action === 'extract') {
        openCustomArchive(cleanTarget);
        updateExtractionConfig({
          archivePath: target,
          destinationPath: dir
        });
        openModal('extract');
      } else {
        openCustomArchive(cleanTarget);
      }
    } else {
      if (action === 'add') {
        openModal('add');
      } else if (action === 'extract') {
        openModal('extract');
      }
    }
  }, [openCustomArchive, openModal, updateCompressionConfig, updateExtractionConfig]);

  // Keep-alive heartbeat to backend engine
  useEffect(() => {
    fetch('/api/ping').catch(() => {});
    const timer = setInterval(() => {
      fetch('/api/ping').catch(() => {});
    }, 2500);
    return () => clearInterval(timer);
  }, []);

  // Graceful shutdown when user closes the window
  useEffect(() => {
    const handleUnload = () => {
      try {
        if (navigator.sendBeacon) {
          navigator.sendBeacon('/api/shutdown');
        } else {
          fetch('/api/shutdown', { method: 'POST', keepalive: true }).catch(() => {});
        }
      } catch (e) {}
    };
    window.addEventListener('beforeunload', handleUnload);
    window.addEventListener('pagehide', handleUnload);
    return () => {
      window.removeEventListener('beforeunload', handleUnload);
      window.removeEventListener('pagehide', handleUnload);
    };
  }, []);

  // Native OS drag-and-drop support
  useEffect(() => {
    const handleDragOver = (e: DragEvent) => {
      e.preventDefault();
      setIsDragOver(true);
    };

    const handleDragLeave = () => {
      setIsDragOver(false);
    };

    const handleDrop = (e: DragEvent) => {
      e.preventDefault();
      setIsDragOver(false);
      if (e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files.length > 0) {
        const file = e.dataTransfer.files[0];
        openCustomArchive(file.name);
      }
    };

    window.addEventListener('dragover', handleDragOver);
    window.addEventListener('dragleave', handleDragLeave);
    window.addEventListener('drop', handleDrop);

    return () => {
      window.removeEventListener('dragover', handleDragOver);
      window.removeEventListener('dragleave', handleDragLeave);
      window.removeEventListener('drop', handleDrop);
    };
  }, [openCustomArchive]);

  return (
    <div className="h-screen w-screen flex flex-col overflow-hidden bg-[#f0f0f0] text-gray-900 select-none relative">
      {/* Classic Menu Bar (directly beneath native Windows title bar) */}
      <MenuBar />

      {/* Main Action Toolbar */}
      <Toolbar />

      {/* Address & Navigation Bar */}
      <AddressBar />

      {/* Resizable Split Workspace */}
      <div className="flex-1 flex overflow-hidden relative">
        <FolderTree />
        <FileTable />

        {/* Drag and Drop Overlay */}
        {isDragOver && (
          <div className="absolute inset-0 bg-blue-600/20 backdrop-blur-xs border-4 border-dashed border-blue-500 z-50 flex items-center justify-center pointer-events-none">
            <div className="bg-white/90 px-6 py-4 rounded-lg shadow-xl text-center">
              <div className="text-lg font-bold text-blue-700">Drop Archive to Inspect</div>
              <div className="text-xs text-gray-600 mt-1">Supports .rar, .zip, .7z, .tar, .gz</div>
            </div>
          </div>
        )}
      </div>

      {/* Dynamic Status Bar */}
      <StatusBar />

      {/* Interactive Modals */}
      <ExtractModal />
      <AddModal />
      <SfxModal />
      <BenchmarkModal />
      <ProgressModal />
      <InfoModal />
      <ViewModal />
      <RepairModal />
      <FindModal />
      <WizardModal />
      <AboutModal />
    </div>
  );
};

export default App;
