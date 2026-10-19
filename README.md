# 📚 VibeRAR — High-Performance WinRAR Clone & Archiver Suite

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform: Windows 10 / 11](https://img.shields.io/badge/Platform-Windows%2010%20%7C%2011-0078D6.svg?logo=windows)](https://microsoft.com)
[![Architecture: x64](https://img.shields.io/badge/Architecture-x64-success.svg)]()
[![Frontend: React 19 + TypeScript](https://img.shields.io/badge/Frontend-React%2019%20%2B%20TypeScript-61DAFB.svg?logo=react)]()
[![Backend: Go 1.26 + Win32](https://img.shields.io/badge/Backend-Go%20%2B%20Win32%20API-00ADD8.svg?logo=go)]()

**VibeRAR** is a modern, high-performance desktop archiver engineered to faithfully replicate the density, ergonomics, and workflows of **WinRAR**, enhanced with a retro-modern visual aesthetic, zero-prompt context menu shortcuts, and seamless Windows 11 shell integration.

---

## 🚀 Key Features

* **Authentic WinRAR Ergonomics & Dual-Density Toolbars**:
  * Seamlessly toggle between **Large (48x36px)** and **Small (24x24px)** classic toolbar buttons.
  * Virtualized data table with high-density columns: `Name`, `Original Size`, `Packed Size`, `Type`, `Modified Date`, `Attributes`, and `CRC32`.
  * Collapsible folder tree hierarchy and breadcrumb address navigation bar.

* **Native Windows Shell Integration (Windows 10 & 11 Modern)**:
  * 🪟 **Windows 11 Top-Level Modern Context Menu**: Implemented via native C++ `IExplorerCommand` COM server paired with an MSIX Sparse Package (`desktop4:FileExplorerContextMenu`). Appears directly on the primary modern WinUI/XAML menu alongside WinRAR above *"Show more options"*!
  * ⚡ **Dynamic Cascaded Submenu**:
    * 📦 *Add to archive...*
    * ⚡ *Add to "<filename>.zip"* — Instant zero-prompt compression (< 50ms).
    * ⚡ *Add to "<filename>.rar"* — Instant zero-prompt compression (< 50ms).
    * 📂 *Extract files...*
    * 🎯 *Extract Here* — Instant in-place extraction (< 50ms).
    * 📁 *Extract to "<folder>\"* — Instant extraction into clean named subfolder.
    * 🔍 *Open with VibeRAR*
  * 🔄 **Classic Shell Verbs**: Complete backward compatibility for Windows 10, Windows 8, and Windows 11 classic *"Show more options"* menu.
  * 🏷️ **File Type Associations**: Complete association for 12 archive formats (`.rar`, `.zip`, `.7z`, `.tar`, `.gz`, `.bz2`, `.iso`, `.cab`, `.xz`, `.001`, `.arj`, `.jar`).

* **Supported Formats**:
  * **Decompression / Inspection**: Full support for RAR (RAR4, RAR5, RAR7), ZIP, TAR, GZ, 7z, BZ2, XZ, CAB, and ISO.
  * **Compression**: High-ratio multi-threaded ZIP64 and RAR5 container archiving.

* **Essential WinRAR Modals Included**:
  1. **Extraction Path and Options**: General & Advanced tabs, folder tree picker, update modes, overwrite modes.
  2. **Add to Archive Dialog**: Format selection (RAR, ZIP), compression level, dictionary size, split to volumes, password encryption.
  3. **Hardware Multi-Threaded CPU Benchmark**: Real-time speed (KB/s), average throughput, memory integrity validation counters, thread count controls, and clipboard export.
  4. **Dual-Progress Operation Monitor**: Real-time progress bars for current file and total volume, byte rate, and pause/cancel controls.
  5. **Archive Information & Comments**: Compression ratio analysis, encryption state, and comment viewer.

* **Self-Extracting Archive (SFX) Subsystem**:
  * Built-in `sfx_builder` engine to generate standalone portable `.exe` archives that unpack autonomously without requiring archiver software installed on the destination machine.

---

## 📦 Download & Installation

You can download the latest compiled executables directly:
* 📥 [**Download VibeRAR_Setup.exe** (Installer with Context Menus)](https://github.com/iegyi/VibeRAR/raw/main/VibeRAR_Setup.exe)
* 📥 [**Download VibeRAR_Portable.exe** (Standalone No-Install)](https://github.com/iegyi/VibeRAR/raw/main/VibeRAR_Portable.exe)
* Or visit the [**Releases**](https://github.com/iegyi/VibeRAR/releases) page for zipped packages.

---

## 🛠️ Building from Source

### Prerequisites
* **Go 1.22+**
* **Node.js 18+** & npm

### Build Instructions

1. **Clone the repository**:
   ```bash
   git clone https://github.com/iegyi/VibeRAR.git
   cd VibeRAR
   ```

2. **Build Frontend Bundle**:
   ```bash
   npm install
   npm run build
   ```

3. **Build VibeRAR Standalone Executable**:
   ```bash
   go build -ldflags="-H windowsgui" -o VibeRAR.exe standalone_runner.go
   ```

4. **Build Installer Wizard**:
   ```bash
   cd tools/installer
   go build -ldflags="-H windowsgui" -o ../../VibeRAR_Setup.exe .
   cd ../..
   ```

---

## 📑 Architecture Specification

Detailed architectural specifications, data structures, and format policies are documented in:
* [WinRAR_Clone_Architecture_Report.pdf](WinRAR_Clone_Architecture_Report.pdf)
* [PRD.md](PRD.md)

---

## 👤 Author

* **Creator**: Ahmed Samir Hany
* **Website**: [https://ahmedsamir.com](https://ahmedsamir.com)

---

## 📄 License

This project is open-source under the [MIT License](LICENSE).
