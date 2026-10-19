import React, { useState } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { X, Wrench, CheckCircle, AlertTriangle } from 'lucide-react';

export const RepairModal: React.FC = () => {
  const { activeModal, closeModal, archiveName } = useArchiveStore();
  const [repairType, setRepairType] = useState<'rar' | 'zip'>('rar');
  const [stage, setStage] = useState<'prompt' | 'repairing' | 'done'>('prompt');
  const [logs, setLogs] = useState<string[]>([]);

  if (activeModal !== 'repair') return null;

  const startRepair = () => {
    setStage('repairing');
    setLogs([
      `Scanning ${archiveName}...`,
      `Analyzing binary blocks and Reed-Solomon parity records...`
    ]);

    setTimeout(() => {
      setLogs((prev) => [...prev, `Found 3% recovery record blocks.`]);
    }, 600);

    setTimeout(() => {
      setLogs((prev) => [
        ...prev,
        `Reconstructing file table headers...`,
        `Rebuilding intact archive: fixed.${archiveName}`
      ]);
      setStage('done');
    }, 1400);
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[480px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>Repairing {archiveName}</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Content Body */}
        <div className="p-4 bg-white flex flex-col gap-3 min-h-[260px]">
          {stage === 'prompt' ? (
            <>
              <div className="flex items-center gap-3 p-2 bg-amber-50 border border-amber-200 rounded-xs">
                <AlertTriangle size={24} className="text-amber-600 shrink-0" />
                <div className="text-xs text-gray-700">
                  Select archive type to reconstruct corrupted data blocks and recovery volumes.
                </div>
              </div>

              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Archive type</legend>
                <div className="flex flex-col gap-2 pt-1">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="repairType"
                      checked={repairType === 'rar'}
                      onChange={() => setRepairType('rar')}
                      className="text-blue-600"
                    />
                    <span>Treat the corrupt archive as RAR (Reed-Solomon recovery)</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="repairType"
                      checked={repairType === 'zip'}
                      onChange={() => setRepairType('zip')}
                      className="text-blue-600"
                    />
                    <span>Treat the corrupt archive as ZIP</span>
                  </label>
                </div>
              </fieldset>

              <div>
                <label className="text-[11px] text-gray-600 block mb-1">
                  Folder to put the repaired archive:
                </label>
                <input
                  type="text"
                  defaultValue="D:\\Projects\\extracted"
                  className="w-full border border-[#7f9db9] rounded-xs px-2 py-1 font-mono text-xs text-gray-900"
                />
              </div>
            </>
          ) : (
            <div className="flex flex-col gap-2">
              <div className="flex items-center gap-2 text-gray-700 font-semibold">
                {stage === 'repairing' ? (
                  <>
                    <Wrench size={16} className="text-blue-600 animate-spin" />
                    <span>Processing parity blocks...</span>
                  </>
                ) : (
                  <>
                    <CheckCircle size={16} className="text-emerald-600" />
                    <span>Repair completed successfully!</span>
                  </>
                )}
              </div>

              <div className="border border-gray-300 rounded-xs p-2 bg-gray-900 text-emerald-400 font-mono text-xs h-40 overflow-y-auto">
                {logs.map((log, i) => (
                  <div key={i}>{log}</div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Action Buttons */}
        <div className="h-11 bg-[#f0f0f0] border-t border-[#dcdcdc] flex items-center justify-end px-3 gap-2">
          {stage === 'prompt' ? (
            <>
              <button
                onClick={startRepair}
                className="winrar-btn w-20 py-1 text-xs font-semibold text-gray-800"
              >
                OK
              </button>
              <button
                onClick={closeModal}
                className="winrar-btn w-20 py-1 text-xs text-gray-800"
              >
                Cancel
              </button>
            </>
          ) : (
            <button
              onClick={closeModal}
              className="winrar-btn w-20 py-1 text-xs font-semibold text-gray-800"
            >
              Close
            </button>
          )}
        </div>
      </div>
    </div>
  );
};
