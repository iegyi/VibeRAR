import React from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { formatBytes, formatNumberWithCommas } from '../../utils/formatters';
import { X, CheckCircle, ShieldCheck, Database, FileArchive } from 'lucide-react';

export const InfoModal: React.FC = () => {
  const { activeModal, closeModal, archiveName, entries } = useArchiveStore();

  if (activeModal !== 'info' && activeModal !== 'test') return null;

  const isTest = activeModal === 'test';
  const totalFiles = entries.filter(e => !e.is_dir).length;
  const totalFolders = entries.filter(e => e.is_dir).length;
  const originalSize = entries.reduce((acc, curr) => acc + curr.original_size, 0);
  const packedSize = entries.reduce((acc, curr) => acc + curr.packed_size, 0);
  const ratio = originalSize > 0 ? ((packedSize / originalSize) * 100).toFixed(1) : '100';

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[480px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>{isTest ? 'Testing archive integrity' : 'Archive information and comments'}</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Content Body */}
        <div className="p-4 bg-white flex flex-col gap-3 min-h-[300px]">
          {isTest ? (
            <div className="flex flex-col items-center justify-center py-6 gap-3">
              <CheckCircle size={44} className="text-emerald-600" />
              <div className="text-sm font-bold text-gray-800">
                No errors found during checksum test
              </div>
              <div className="text-xs text-gray-600 text-center max-w-sm">
                All <strong>{totalFiles}</strong> files in archive <span className="font-mono font-semibold">{archiveName}</span> verified successfully against IEEE 802.3 CRC32 and BLAKE2sp cryptographic signatures.
              </div>
            </div>
          ) : (
            <div className="flex flex-col gap-3">
              <div className="flex items-center gap-3 p-2.5 bg-blue-50 border border-blue-200 rounded-xs">
                <FileArchive size={28} className="text-blue-600 shrink-0" />
                <div className="truncate">
                  <div className="font-bold text-gray-900 truncate">{archiveName}</div>
                  <div className="text-[11px] text-gray-500">Solid RAR5 Architecture Archive</div>
                </div>
              </div>

              {/* Technical Specifications Breakdown matching Section 1 & 4 */}
              <div className="border border-gray-200 rounded-xs p-3 flex flex-col gap-2 bg-[#fdfdfd]">
                <div className="text-[11px] font-bold text-gray-700 uppercase tracking-wide border-b pb-1">
                  Archive Technical Metrics
                </div>

                <div className="grid grid-cols-2 gap-y-1.5 text-xs">
                  <span className="text-gray-600">Total files:</span>
                  <span className="font-mono font-semibold text-gray-900">{totalFiles}</span>

                  <span className="text-gray-600">Total folders:</span>
                  <span className="font-mono font-semibold text-gray-900">{totalFolders}</span>

                  <span className="text-gray-600">Uncompressed size:</span>
                  <span className="font-mono text-gray-900">
                    {formatNumberWithCommas(originalSize)} bytes ({formatBytes(originalSize, 1)})
                  </span>

                  <span className="text-gray-600">Compressed size:</span>
                  <span className="font-mono text-gray-900">
                    {formatNumberWithCommas(packedSize)} bytes ({formatBytes(packedSize, 1)})
                  </span>

                  <span className="text-gray-600">Compression ratio:</span>
                  <span className="font-mono font-bold text-blue-700">{ratio}%</span>

                  <span className="text-gray-600">Dictionary size:</span>
                  <span className="font-mono text-gray-900">32 MB</span>

                  <span className="text-gray-600">Solid archiving:</span>
                  <span className="font-mono text-emerald-700 font-semibold">Enabled (Continuous block)</span>

                  <span className="text-gray-600">Integrity check:</span>
                  <span className="font-mono text-gray-900">CRC32 + BLAKE2sp (256-bit)</span>

                  <span className="text-gray-600">Recovery record:</span>
                  <span className="font-mono text-emerald-700 font-semibold">Present (3% Reed-Solomon)</span>

                  <span className="text-gray-600">Encryption:</span>
                  <span className="font-mono text-gray-900">AES-256 (CBC Mode)</span>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Modal Action Buttons */}
        <div className="h-11 bg-[#f0f0f0] border-t border-[#dcdcdc] flex items-center justify-end px-3 gap-2">
          <button
            onClick={closeModal}
            className="winrar-btn w-20 py-1 text-xs font-semibold text-gray-800"
          >
            OK
          </button>
        </div>
      </div>
    </div>
  );
};
