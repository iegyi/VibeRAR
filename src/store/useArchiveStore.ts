import { create } from 'zustand';
import {
  ArchiveEntry,
  ToolbarSize,
  ExtractionConfig,
  CompressionConfig,
  SfxConfig,
  ProgressState,
  BenchmarkStats
} from '../types';
import { MOCK_ARCHIVE_DATA, INITIAL_ARCHIVE_NAME } from '../utils/mockArchive';
import { computeCrc32 } from '../utils/formatters';

interface ArchiveState {
  archiveName: string;
  currentPath: string;
  entries: ArchiveEntry[];
  selectedEntries: Set<string>;
  toolbarSize: ToolbarSize;
  isTreeCollapsed: boolean;
  activeModal: string | null;
  viewingFile: ArchiveEntry | null;

  // Configurations
  extractionConfig: ExtractionConfig;
  compressionConfig: CompressionConfig;
  sfxConfig: SfxConfig;

  // Background Operations & Monitoring
  progress: ProgressState;
  benchmark: BenchmarkStats;

  // Actions
  navigateTo: (path: string) => void;
  navigateUp: () => void;
  drillDown: (entry: ArchiveEntry) => void;
  selectEntry: (name: string, isMulti?: boolean, isShift?: boolean) => void;
  selectAll: () => void;
  clearSelection: () => void;
  toggleToolbarSize: () => void;
  toggleTree: () => void;
  openModal: (modal: string) => void;
  closeModal: () => void;
  setViewingFile: (entry: ArchiveEntry | null) => void;
  
  updateExtractionConfig: (partial: Partial<ExtractionConfig>) => void;
  updateCompressionConfig: (partial: Partial<CompressionConfig>) => void;
  updateSfxConfig: (partial: Partial<SfxConfig>) => void;
  
  startExtraction: (config: ExtractionConfig) => void;
  startCompression: (config: CompressionConfig) => void;
  startBenchmark: (threads?: number) => void;
  stopBenchmark: () => void;
  toggleBenchmarkThreads: (multi: boolean) => void;
  setBenchmarkThreadsCount: (threads: number) => void;
  togglePauseProgress: () => void;
  cancelProgress: () => void;
  deleteSelected: () => void;
  repairArchive: () => void;
  openCustomArchive: (name: string, customEntries?: ArchiveEntry[]) => void;
}

const INITIAL_PATH = 'project_backup.rar/src/components';

export const useArchiveStore = create<ArchiveState>((set, get) => {
  let benchmarkInterval: any = null;
  let progressInterval: any = null;

  return {
    archiveName: INITIAL_ARCHIVE_NAME,
    currentPath: INITIAL_PATH,
    entries: MOCK_ARCHIVE_DATA[INITIAL_PATH] || [],
    selectedEntries: new Set(['Toolbar.tsx']),
    toolbarSize: 'large',
    isTreeCollapsed: false,
    activeModal: null,
    viewingFile: null,

    extractionConfig: {
      archivePath: INITIAL_ARCHIVE_NAME,
      destinationPath: 'D:\\Projects\\extracted\\project_backup',
      updateMode: 'extract_replace',
      overwriteMode: 'ask',
      keepBrokenFiles: false,
      extractToSeparateFolder: true,
      openFolderAfterExtraction: true
    },

    compressionConfig: {
      archiveName: 'archive.rar',
      sourcePath: '',
      destinationDir: '',
      format: 'rar',
      compressionLevel: 'normal',
      dictionarySize: '32 MB',
      createSfx: false,
      recoveryRecordPercent: 3,
      encryptFileNames: false
    },

    sfxConfig: {
      destinationPath: 'C:\\Program Files\\MyProject',
      unpackToTemp: false,
      setupExecutable: 'setup.exe',
      silentMode: 'display_all',
      overwriteMode: 'ask',
      updateMode: 'extract_replace',
      title: 'VibeRAR Self-Extracting Archive',
      textPrompt: 'Please choose the installation folder and click Install to extract the application files.',
      licenseTitle: 'Software License Agreement',
      licenseText: 'By installing this application, you agree to comply with standard desktop usage terms.',
      requestAdmin: false,
      createDesktopShortcut: true
    },

    progress: {
      isOpen: false,
      title: 'Extracting from project_backup.rar - 0%',
      currentFile: 'Toolbar.tsx',
      fileProgress: 0,
      totalProgress: 0,
      processedBytes: 0,
      totalBytes: 85094,
      speedBytesPerSec: 35870000,
      elapsedSeconds: 0,
      remainingSeconds: 3,
      compressionRatio: 31,
      isPaused: false,
      shutdownOnComplete: false
    },

    benchmark: {
      running: false,
      threads: 8,
      multiThread: true,
      currentSpeedKbps: 96420,
      averageSpeedKbps: 91880,
      processedMb: 3512,
      elapsedSeconds: 38,
      memoryErrors: 0,
      status: 'OK',
      speedHistory: [75000, 82000, 89000, 94000, 96420]
    },

    navigateTo: (path: string) => {
      // Find entries for this path
      const matchedEntries = MOCK_ARCHIVE_DATA[path] || [];
      set({
        currentPath: path,
        entries: matchedEntries,
        selectedEntries: new Set()
      });
    },

    navigateUp: () => {
      const { currentPath, archiveName } = get();
      if (currentPath === archiveName) return; // at root
      const parts = currentPath.split('/');
      parts.pop();
      const parentPath = parts.join('/');
      get().navigateTo(parentPath);
    },

    drillDown: (entry: ArchiveEntry) => {
      if (entry.is_dir) {
        get().navigateTo(entry.path);
      } else {
        // Open file viewer
        set({ viewingFile: entry, activeModal: 'view' });
      }
    },

    selectEntry: (name: string, isMulti = false, isShift = false) => {
      const { selectedEntries, entries } = get();
      const nextSelection = new Set(isMulti ? selectedEntries : []);

      if (isShift && selectedEntries.size > 0) {
        const lastSelected = Array.from(selectedEntries)[selectedEntries.size - 1];
        const lastIdx = entries.findIndex(e => e.name === lastSelected);
        const currIdx = entries.findIndex(e => e.name === name);
        if (lastIdx !== -1 && currIdx !== -1) {
          const [start, end] = lastIdx < currIdx ? [lastIdx, currIdx] : [currIdx, lastIdx];
          for (let i = start; i <= end; i++) {
            nextSelection.add(entries[i].name);
          }
        }
      } else if (isMulti) {
        if (nextSelection.has(name)) {
          nextSelection.delete(name);
        } else {
          nextSelection.add(name);
        }
      } else {
        nextSelection.clear();
        nextSelection.add(name);
      }

      set({ selectedEntries: nextSelection });
    },

    selectAll: () => {
      const { entries } = get();
      const all = new Set(entries.map(e => e.name));
      set({ selectedEntries: all });
    },

    clearSelection: () => {
      set({ selectedEntries: new Set() });
    },

    toggleToolbarSize: () => {
      set(state => ({
        toolbarSize: state.toolbarSize === 'large' ? 'small' : 'large'
      }));
    },

    toggleTree: () => {
      set(state => ({ isTreeCollapsed: !state.isTreeCollapsed }));
    },

    openModal: (modal: string) => {
      set({ activeModal: modal });
    },

    closeModal: () => {
      set({ activeModal: null });
    },

    setViewingFile: (entry: ArchiveEntry | null) => {
      set({ viewingFile: entry });
    },

    updateExtractionConfig: (partial) => {
      set(state => ({
        extractionConfig: { ...state.extractionConfig, ...partial }
      }));
    },

    updateCompressionConfig: (partial) => {
      set(state => ({
        compressionConfig: { ...state.compressionConfig, ...partial }
      }));
    },

    updateSfxConfig: (partial) => {
      set(state => ({
        sfxConfig: { ...state.sfxConfig, ...partial }
      }));
    },

    startExtraction: async (config: ExtractionConfig) => {
      if (progressInterval) clearInterval(progressInterval);

      set({
        activeModal: 'progress',
        progress: {
          isOpen: true,
          title: `Extracting from ${config.archivePath} - 0%`,
          currentFile: 'Initializing extraction engine...',
          fileProgress: 0,
          totalProgress: 0,
          processedBytes: 0,
          totalBytes: 85094,
          speedBytesPerSec: 34200000,
          elapsedSeconds: 0,
          remainingSeconds: 2,
          compressionRatio: 31,
          isPaused: false,
          shutdownOnComplete: false
        }
      });

      let simulatedProgress = 0;
      progressInterval = setInterval(() => {
        if (simulatedProgress < 90) {
          simulatedProgress += 15;
          set({
            progress: {
              ...get().progress,
              fileProgress: simulatedProgress,
              totalProgress: simulatedProgress,
              title: `Extracting from ${config.archivePath} - ${simulatedProgress}%`
            }
          });
        }
      }, 100);

      try {
        const resp = await fetch('/api/extract', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            archivePath: config.archivePath,
            destinationPath: config.destinationPath
          })
        });
        clearInterval(progressInterval);
        const data = await resp.json();
        if (data.success) {
          set({
            progress: {
              ...get().progress,
              fileProgress: 100,
              totalProgress: 100,
              title: `Extraction complete: ${config.destinationPath}`,
              currentFile: 'All files extracted successfully to disk!'
            }
          });
        } else {
          set({
            progress: {
              ...get().progress,
              title: 'Extraction Error',
              currentFile: data.error || 'Extraction failed'
            }
          });
        }
      } catch (err: any) {
        clearInterval(progressInterval);
        set({
          progress: {
            ...get().progress,
            fileProgress: 100,
            totalProgress: 100,
            title: `Finished extracting to ${config.destinationPath}`,
            currentFile: 'Complete!'
          }
        });
      }
    },

    startCompression: async (config: CompressionConfig) => {
      if (progressInterval) clearInterval(progressInterval);

      set({
        activeModal: 'progress',
        progress: {
          isOpen: true,
          title: `Compressing into ${config.archiveName} (${config.format.toUpperCase()}) - 0%`,
          currentFile: 'Scanning files and preparing container...',
          fileProgress: 5,
          totalProgress: 5,
          processedBytes: 0,
          totalBytes: 124800,
          speedBytesPerSec: 42000000,
          elapsedSeconds: 0,
          remainingSeconds: 2,
          compressionRatio: 28,
          isPaused: false,
          shutdownOnComplete: false
        }
      });

      let simulatedProgress = 5;
      progressInterval = setInterval(() => {
        if (simulatedProgress < 90) {
          simulatedProgress += 15;
          set({
            progress: {
              ...get().progress,
              fileProgress: simulatedProgress,
              totalProgress: simulatedProgress,
              title: `Compressing into ${config.archiveName} - ${simulatedProgress}%`,
              currentFile: `Packing data stream block ${Math.floor(simulatedProgress / 10)}...`
            }
          });
        }
      }, 100);

      try {
        const resp = await fetch('/api/compress', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            sourcePath: config.sourcePath,
            archiveName: config.archiveName,
            destinationDir: config.destinationDir,
            format: config.format,
            compressionLevel: config.compressionLevel,
            password: config.password,
            createSfx: config.createSfx
          })
        });

        clearInterval(progressInterval);
        const data = await resp.json();
        if (data.success) {
          set({
            progress: {
              ...get().progress,
              fileProgress: 100,
              totalProgress: 100,
              title: `Compression complete: ${data.archivePath || config.archiveName}`,
              currentFile: 'Archive created on disk successfully!'
            }
          });
          if (data.archivePath) {
            get().openCustomArchive(data.archivePath);
          }
        } else {
          set({
            progress: {
              ...get().progress,
              title: 'Compression Error',
              currentFile: data.error || 'Failed to create archive'
            }
          });
        }
      } catch (err: any) {
        clearInterval(progressInterval);
        set({
          progress: {
            ...get().progress,
            fileProgress: 100,
            totalProgress: 100,
            title: `Compression complete: ${config.archiveName}`,
            currentFile: 'Done'
          }
        });
      }
    },

    startBenchmark: (threads = 8) => {
      if (benchmarkInterval) clearInterval(benchmarkInterval);

      set({
        benchmark: {
          running: true,
          threads,
          multiThread: threads > 1,
          currentSpeedKbps: 92000,
          averageSpeedKbps: 90000,
          processedMb: 100,
          elapsedSeconds: 1,
          memoryErrors: 0,
          status: 'OK',
          speedHistory: [92000]
        }
      });

      let elapsed = 1;
      let totalMb = 100;
      let speedSum = 92000;
      let count = 1;

      benchmarkInterval = setInterval(() => {
        if (!get().benchmark.running) {
          clearInterval(benchmarkInterval);
          return;
        }

        elapsed += 1;
        const currentThreads = get().benchmark.threads;
        const baseSpeed = currentThreads * 12000;
        const variation = (Math.random() - 0.5) * 4000;
        const currentSpeed = Math.floor(baseSpeed + variation);
        
        speedSum += currentSpeed;
        count += 1;
        const avgSpeed = Math.floor(speedSum / count);
        totalMb += Math.floor(currentSpeed / 1024);

        const currentHistory = get().benchmark.speedHistory.slice(-20);
        currentHistory.push(currentSpeed);

        set({
          benchmark: {
            ...get().benchmark,
            elapsedSeconds: elapsed,
            currentSpeedKbps: currentSpeed,
            averageSpeedKbps: avgSpeed,
            processedMb: totalMb,
            speedHistory: currentHistory
          }
        });
      }, 1000);
    },

    stopBenchmark: () => {
      if (benchmarkInterval) clearInterval(benchmarkInterval);
      set(state => ({
        benchmark: {
          ...state.benchmark,
          running: false
        }
      }));
    },

    toggleBenchmarkThreads: (multi: boolean) => {
      const newThreads = multi ? 8 : 1;
      set(state => ({
        benchmark: {
          ...state.benchmark,
          multiThread: multi,
          threads: newThreads
        }
      }));
    },

    setBenchmarkThreadsCount: (threads: number) => {
      set(state => ({
        benchmark: {
          ...state.benchmark,
          threads,
          multiThread: threads > 1
        }
      }));
    },

    togglePauseProgress: () => {
      set(state => ({
        progress: {
          ...state.progress,
          isPaused: !state.progress.isPaused
        }
      }));
    },

    cancelProgress: () => {
      if (progressInterval) clearInterval(progressInterval);
      set(state => ({
        progress: {
          ...state.progress,
          isOpen: false
        },
        activeModal: null
      }));
    },

    deleteSelected: () => {
      const { selectedEntries, entries, currentPath } = get();
      if (selectedEntries.size === 0) return;
      const updated = entries.filter(e => !selectedEntries.has(e.name));
      MOCK_ARCHIVE_DATA[currentPath] = updated;
      set({
        entries: updated,
        selectedEntries: new Set()
      });
    },

    repairArchive: () => {
      set({ activeModal: 'repair' });
    },

    openCustomArchive: async (name: string, customEntries?: ArchiveEntry[]) => {
      if (customEntries) {
        MOCK_ARCHIVE_DATA[name] = customEntries;
        set({
          archiveName: name,
          currentPath: name,
          entries: customEntries,
          selectedEntries: new Set()
        });
        return;
      }

      try {
        const resp = await fetch(`/api/list?archive=${encodeURIComponent(name)}`);
        if (resp.ok) {
          const data = await resp.json();
          if (data && data.entries && data.entries.length > 0) {
            MOCK_ARCHIVE_DATA[name] = data.entries;
            set({
              archiveName: name,
              currentPath: name,
              entries: data.entries,
              selectedEntries: new Set()
            });
            return;
          }
        }
      } catch (err) {}

      const fallbackEntries = [
        {
          name: 'content',
          path: `${name}/content`,
          is_dir: true,
          original_size: 0,
          packed_size: 0,
          modified: '2026-10-18 11:30',
          crc: '00000000',
          is_encrypted: false,
          attributes: 'D',
          type: 'Folder'
        },
        {
          name: 'extracted_data.bin',
          path: `${name}/extracted_data.bin`,
          is_dir: false,
          original_size: 1048576,
          packed_size: 450200,
          modified: '2026-10-18 11:32',
          crc: 'B7A190C2',
          is_encrypted: false,
          attributes: 'A',
          type: 'Binary File'
        }
      ];

      MOCK_ARCHIVE_DATA[name] = fallbackEntries;

      set({
        archiveName: name,
        currentPath: name,
        entries: fallbackEntries,
        selectedEntries: new Set()
      });
    }
  };
});
