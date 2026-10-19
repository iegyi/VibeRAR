# Product Requirements Document (PRD) - System Archiver Desktop Application (VibeRAR)

## 1. Scope & Objective
High-performance cross-platform archive manager mimicking WinRAR's functional workflows and UI density, engineered using Tauri v2, Rust, React 19, TypeScript, and Tailwind CSS.

## 2. Technical Stack Boundaries
- **Backend Systems Engine**: Rust (stable edition), Tokio async runtime, unrar / zip / sevenz-rust, rayon multi-threaded worker pools.
- **IPC Bridge**: Tauri v2 Core IPC with typed commands (`invoke`) and streaming events (`op-progress`).
- **Presentation Layer**: React 19, TypeScript, Tailwind CSS, TanStack Virtual / Table, Zustand reactive state management.
- **Architectural Rules**: Zero file processing on the UI thread. All disk read/writes must stream through Rust worker threads with real-time byte chunk progress emission.

## 3. Supported Formats Policy
- **Decompression / Inspection**: Full support for RAR (RAR4, RAR5, RAR7), ZIP, TAR, GZ, 7z, BZ2, XZ, Zstandard, CAB, ISO.
- **Creation / Compression**: ZIP (with Zip64 support) and 7z with AES-256 encryption.
- **Licensing Constraint**: Creation of proprietary RAR is out of scope due to license restrictions; decompression uses unrar.

## 4. UI/UX Specifications
- **Retro-Modern Desktop Aesthetic**: Dual-density toolbars (large 48x36px and small 24x24px), tabular data view, split-pane folder hierarchy.
- **Color Palette & Theme**: Classic Windows desktop styling with light neutral slate tones, crisp 1px borders, dense data layout, high-contrast WinRAR-style visual iconography.
- **Essential Modals**:
  1. **Extraction Path and Options** (General & Advanced tabs, folder tree picker, update modes, overwrite modes, broken file preservation).
  2. **Advanced SFX Configuration** (6 tabs: General, Setup, Modes, Text & Icon, License, Advanced).
  3. **Benchmark and Hardware Test** (Real-time speed KB/s, average speed, memory integrity counter, multi-threading control).
  4. **Operation Progress Monitor** (Dual-progress bars for current file and total volume, streaming transfer rate, pause/cancel/background).
  5. **Archive Information / Comments** (Detailed compression ratios, encryption status, dictionary size, headers).
  6. **Add to Archive Dialog** (Compression method, archive format selection, split to volumes, password protection).
