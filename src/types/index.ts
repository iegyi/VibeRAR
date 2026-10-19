export interface ArchiveEntry {
  name: string;
  path: string;
  is_dir: boolean;
  original_size: number;
  packed_size: number;
  modified: string;
  crc: string;
  is_encrypted: boolean;
  attributes: string;
  type?: string;
  content?: string; // Text preview if available
}

export type ToolbarSize = 'large' | 'small';

export type UpdateMode = 'extract_replace' | 'extract_update' | 'fresh_only';
export type OverwriteMode = 'ask' | 'overwrite' | 'skip' | 'rename_auto';

export interface ExtractionConfig {
  archivePath: string;
  destinationPath: string;
  updateMode: UpdateMode;
  overwriteMode: OverwriteMode;
  keepBrokenFiles: boolean;
  extractToSeparateFolder: boolean;
  openFolderAfterExtraction: boolean;
  password?: string;
}

export type CompressionFormat = 'rar' | 'zip' | '7z' | 'tar';
export type CompressionLevel = 'store' | 'fastest' | 'fast' | 'normal' | 'good' | 'best';

export interface CompressionConfig {
  archiveName: string;
  sourcePath?: string;
  destinationDir: string;
  format: CompressionFormat;
  compressionLevel: CompressionLevel;
  dictionarySize: string;
  createSfx: boolean;
  recoveryRecordPercent: number;
  password?: string;
  encryptFileNames: boolean;
  splitVolumes?: string;
}

export interface SfxConfig {
  destinationPath: string;
  unpackToTemp: boolean;
  setupExecutable: string;
  silentMode: 'display_all' | 'hide_start' | 'hide_all';
  overwriteMode: 'ask' | 'overwrite' | 'skip';
  updateMode: 'extract_replace' | 'extract_update';
  title: string;
  textPrompt: string;
  licenseTitle: string;
  licenseText: string;
  requestAdmin: boolean;
  createDesktopShortcut: boolean;
}

export interface ProgressState {
  isOpen: boolean;
  title: string;
  currentFile: string;
  fileProgress: number;
  totalProgress: number;
  processedBytes: number;
  totalBytes: number;
  speedBytesPerSec: number;
  elapsedSeconds: number;
  remainingSeconds: number;
  compressionRatio: number;
  isPaused: boolean;
  shutdownOnComplete: boolean;
}

export interface BenchmarkStats {
  running: boolean;
  threads: number;
  multiThread: boolean;
  currentSpeedKbps: number;
  averageSpeedKbps: number;
  processedMb: number;
  elapsedSeconds: number;
  memoryErrors: number;
  status: string;
  speedHistory: number[];
}

export interface TreeNode {
  id: string;
  name: string;
  path: string;
  isFolder: boolean;
  children?: TreeNode[];
  expanded?: boolean;
}
