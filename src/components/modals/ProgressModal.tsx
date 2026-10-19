import React from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { formatBytes, formatSpeed, formatTime } from '../../utils/formatters';
import { X, Play, Pause, ArrowDown } from 'lucide-react';

export const ProgressModal: React.FC = () => {
  const { progress, togglePauseProgress, cancelProgress, activeModal } = useArchiveStore();

  if (activeModal !== 'progress' || !progress.isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[500px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar matching Figure 5 */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span className="truncate">{progress.title}</span>
          <button onClick={cancelProgress} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Content Body */}
        <div className="p-4 bg-white flex flex-col gap-3">
          {/* Current File Section */}
          <div>
            <div className="flex justify-between items-center mb-1 text-gray-700">
              <span className="font-semibold text-gray-800">Extracting</span>
              <span className="font-mono text-xs text-blue-700 truncate max-w-[320px]">
                {progress.currentFile}
              </span>
            </div>
            {/* Current File Progress Bar */}
            <div className="h-4 bg-[#e8e8e8] border border-[#7f9db9] rounded-xs overflow-hidden p-[1px]">
              <div
                className="h-full bg-linear-to-r from-[#4caf50] to-[#81c784] transition-all duration-150 rounded-xs"
                style={{ width: `${progress.fileProgress}%` }}
              />
            </div>
          </div>

          {/* Total Progress Section */}
          <div>
            <div className="flex justify-between items-center mb-1 text-gray-700">
              <span className="font-semibold text-gray-800">Total progress</span>
              <span className="font-mono font-bold text-gray-900 text-xs">
                {progress.totalProgress}%
              </span>
            </div>
            {/* Total Progress Bar */}
            <div className="h-4 bg-[#e8e8e8] border border-[#7f9db9] rounded-xs overflow-hidden p-[1px]">
              <div
                className="h-full bg-linear-to-r from-[#1976d2] to-[#64b5f6] transition-all duration-150 rounded-xs"
                style={{ width: `${progress.totalProgress}%` }}
              />
            </div>
          </div>

          {/* Metrics Grid matching Figure 5 */}
          <div className="mt-2 pt-2 border-t border-gray-100 flex flex-col gap-1.5 text-xs">
            <div className="flex justify-between items-center">
              <span className="text-gray-600">Elapsed time</span>
              <span className="font-mono text-gray-900">{formatTime(progress.elapsedSeconds)}</span>
            </div>

            <div className="flex justify-between items-center">
              <span className="text-gray-600">Time left</span>
              <span className="font-mono text-gray-900">{formatTime(progress.remainingSeconds)}</span>
            </div>

            <div className="flex justify-between items-center">
              <span className="text-gray-600">Processed</span>
              <span className="font-mono text-gray-900">
                {formatBytes(progress.processedBytes, 1)} of {formatBytes(progress.totalBytes, 1)}
              </span>
            </div>

            <div className="flex justify-between items-center">
              <span className="text-gray-600">Speed</span>
              <span className="font-mono text-gray-900 font-semibold">
                {formatSpeed(progress.speedBytesPerSec)}
              </span>
            </div>

            <div className="flex justify-between items-center">
              <span className="text-gray-600">Compression ratio</span>
              <span className="font-mono text-gray-900">{progress.compressionRatio}%</span>
            </div>
          </div>

          {/* Operational Buttons matching Figure 5 */}
          <div className="flex items-center justify-center gap-2 mt-2">
            {progress.totalProgress >= 100 ? (
              <button
                onClick={cancelProgress}
                className="winrar-btn w-28 py-1.5 text-xs font-bold text-gray-900 bg-blue-100 border-blue-400"
              >
                Close
              </button>
            ) : (
              <>
                <button
                  onClick={togglePauseProgress}
                  className="winrar-btn w-24 py-1 text-xs font-semibold text-gray-800 flex items-center justify-center gap-1"
                >
                  {progress.isPaused ? (
                    <>
                      <Play size={12} className="text-green-600" />
                      <span>Resume</span>
                    </>
                  ) : (
                    <>
                      <Pause size={12} className="text-amber-600" />
                      <span>Pause</span>
                    </>
                  )}
                </button>

                <button
                  onClick={cancelProgress}
                  className="winrar-btn w-24 py-1 text-xs text-gray-800 flex items-center justify-center gap-1"
                  title="Run task in background"
                >
                  <ArrowDown size={12} className="text-blue-600" />
                  <span>Background</span>
                </button>

                <button
                  onClick={cancelProgress}
                  className="winrar-btn w-24 py-1 text-xs text-gray-800"
                >
                  Cancel
                </button>
              </>
            )}
          </div>

          {/* Auto Shutdown Checkbox matching Figure 5 */}
          <div className="mt-1 pt-2 border-t border-gray-100 flex items-center gap-2">
            <input
              type="checkbox"
              id="shutdown-check"
              checked={progress.shutdownOnComplete}
              onChange={() =>
                useArchiveStore.setState((s) => ({
                  progress: { ...s.progress, shutdownOnComplete: !s.progress.shutdownOnComplete }
                }))
              }
              className="text-blue-600 rounded-xs"
            />
            <label htmlFor="shutdown-check" className="text-gray-700 cursor-pointer text-[11px]">
              Shut down computer after the operation completes
            </label>
          </div>
        </div>
      </div>
    </div>
  );
};
