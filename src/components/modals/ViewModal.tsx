import React, { useState } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { formatBytes } from '../../utils/formatters';
import { X, FileText, Binary, Copy } from 'lucide-react';

export const ViewModal: React.FC = () => {
  const { activeModal, closeModal, viewingFile, selectedEntries, entries } = useArchiveStore();
  const [viewMode, setViewMode] = useState<'text' | 'hex'>('text');

  if (activeModal !== 'view') return null;

  // If viewingFile is set, use it; otherwise use the first selected file
  const targetEntry =
    viewingFile ||
    entries.find((e) => selectedEntries.has(e.name) && !e.is_dir) ||
    entries.find((e) => !e.is_dir) || {
      name: 'Toolbar.tsx',
      original_size: 14208,
      crc: '9F3A12C7',
      content: '// WinRAR Toolbar with Large (48x36px) & Small (24x24px) Button Support\nexport const Toolbar = () => {\n  return <div className="toolbar">Action Buttons</div>;\n};\n'
    };

  const contentText =
    (targetEntry as any).content ||
    `// File: ${targetEntry.name}\n// Original Size: ${formatBytes(targetEntry.original_size)}\n// CRC32: ${targetEntry.crc}\n\nimport React from "react";\n\nexport default function ${targetEntry.name.replace(/\.[^.]+$/, '')}() {\n  return (\n    <div className="component-container">\n      <h2>${targetEntry.name}</h2>\n      <p>Loaded from archive stream without disk write.</p>\n    </div>\n  );\n}`;

  // Generate synthetic hex dump
  const generateHexDump = (text: string) => {
    const lines = [];
    const encoder = new TextEncoder();
    const bytes = encoder.encode(text.slice(0, 512));

    for (let i = 0; i < bytes.length; i += 16) {
      const offset = i.toString(16).padStart(8, '0').toUpperCase();
      const chunk = bytes.slice(i, i + 16);
      const hex = Array.from(chunk)
        .map((b) => b.toString(16).padStart(2, '0').toUpperCase())
        .join(' ')
        .padEnd(48, ' ');
      const ascii = Array.from(chunk)
        .map((b) => (b >= 32 && b <= 126 ? String.fromCharCode(b) : '.'))
        .join('');
      lines.push(`${offset}  ${hex}  |${ascii}|`);
    }
    return lines.join('\n');
  };

  const copyContent = () => {
    navigator.clipboard.writeText(contentText);
    alert('File contents copied to clipboard!');
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[620px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span className="truncate">{targetEntry.name} — Internal Viewer</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Viewer Toolbar */}
        <div className="bg-[#ece9d8] border-b border-[#c0c0c0] px-3 py-1 flex items-center justify-between text-xs">
          <div className="flex items-center gap-1">
            <button
              onClick={() => setViewMode('text')}
              className={`winrar-btn px-2.5 py-0.5 flex items-center gap-1 ${
                viewMode === 'text' ? 'font-bold bg-blue-100 text-blue-800' : 'text-gray-700'
              }`}
            >
              <FileText size={12} />
              <span>Text / Code</span>
            </button>
            <button
              onClick={() => setViewMode('hex')}
              className={`winrar-btn px-2.5 py-0.5 flex items-center gap-1 ${
                viewMode === 'hex' ? 'font-bold bg-blue-100 text-blue-800' : 'text-gray-700'
              }`}
            >
              <Binary size={12} />
              <span>Hex Dump</span>
            </button>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-[11px] text-gray-500 font-mono">
              CRC32: {targetEntry.crc}
            </span>
            <button
              onClick={copyContent}
              className="winrar-btn px-2 py-0.5 flex items-center gap-1 text-gray-700"
              title="Copy to clipboard"
            >
              <Copy size={11} />
              <span>Copy</span>
            </button>
          </div>
        </div>

        {/* Content Preview Box */}
        <div className="p-3 bg-white flex-1 min-h-[340px] max-h-[460px] overflow-auto">
          {viewMode === 'text' ? (
            <pre className="font-mono text-xs text-gray-800 whitespace-pre-wrap leading-relaxed select-text">
              {contentText}
            </pre>
          ) : (
            <pre className="font-mono text-[11px] text-gray-800 leading-tight select-text">
              {generateHexDump(contentText)}
            </pre>
          )}
        </div>

        {/* Bottom Bar */}
        <div className="h-9 bg-[#f0f0f0] border-t border-[#dcdcdc] flex items-center justify-between px-3">
          <span className="text-[11px] text-gray-500">
            Size: {formatBytes(targetEntry.original_size)} | UTF-8 Encoding
          </span>
          <button
            onClick={closeModal}
            className="winrar-btn w-20 py-1 text-xs text-gray-800 font-semibold"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};
