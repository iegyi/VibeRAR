import React, { useEffect, useRef } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { formatNumberWithCommas, formatTime } from '../../utils/formatters';
import { X } from 'lucide-react';

export const BenchmarkModal: React.FC = () => {
  const {
    activeModal,
    closeModal,
    benchmark,
    startBenchmark,
    stopBenchmark,
    toggleBenchmarkThreads,
    setBenchmarkThreadsCount
  } = useArchiveStore();

  const canvasRef = useRef<HTMLCanvasElement>(null);

  // Auto-start benchmark when opened
  useEffect(() => {
    if (activeModal === 'benchmark') {
      startBenchmark(benchmark.threads);
    } else {
      stopBenchmark();
    }
  }, [activeModal]);

  // Render throughput canvas graph matching Figure 4
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const width = canvas.width;
    const height = canvas.height;

    // Dark background matching Figure 4
    ctx.fillStyle = '#0f172a';
    ctx.fillRect(0, 0, width, height);

    // Subtle grid lines
    ctx.strokeStyle = '#1e293b';
    ctx.lineWidth = 1;
    for (let y = 15; y < height; y += 20) {
      ctx.beginPath();
      ctx.moveTo(0, y);
      ctx.lineTo(width, y);
      ctx.stroke();
    }

    const history = benchmark.speedHistory;
    if (history.length > 1) {
      ctx.strokeStyle = '#38bdf8';
      ctx.lineWidth = 2;
      ctx.beginPath();

      const maxSpeed = 150000;
      const step = width / (history.length - 1);

      history.forEach((speed, i) => {
        const x = i * step;
        const normalized = Math.min(1, speed / maxSpeed);
        const y = height - normalized * (height - 10) - 5;
        if (i === 0) {
          ctx.moveTo(x, y);
        } else {
          ctx.lineTo(x, y);
        }
      });
      ctx.stroke();

      // Soft glow fill
      ctx.lineTo(width, height);
      ctx.lineTo(0, height);
      ctx.fillStyle = 'rgba(56, 189, 248, 0.15)';
      ctx.fill();
    }
  }, [benchmark.speedHistory]);

  if (activeModal !== 'benchmark') return null;

  const copyResults = () => {
    const text = `VibeRAR Benchmark and Hardware Test:
Threads: ${benchmark.threads}
Current Speed: ${formatNumberWithCommas(benchmark.currentSpeedKbps)} KB/s
Average Speed: ${formatNumberWithCommas(benchmark.averageSpeedKbps)} KB/s
Processed: ${formatNumberWithCommas(benchmark.processedMb)} MB
Elapsed: ${formatTime(benchmark.elapsedSeconds)}
Memory Errors: ${benchmark.memoryErrors}
Status: ${benchmark.status}`;
    navigator.clipboard.writeText(text);
    alert('Benchmark metrics copied to clipboard!');
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[490px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>Benchmark and hardware test</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Top Controls */}
        <div className="p-3 bg-white flex items-center justify-between border-b border-[#e0e0e0]">
          <div className="flex items-center gap-2">
            <span className="font-semibold text-gray-700">Threads:</span>
            <select
              value={benchmark.threads}
              onChange={(e) => setBenchmarkThreadsCount(Number(e.target.value))}
              disabled={!benchmark.multiThread}
              className="border border-[#7f9db9] rounded-xs px-2 py-0.5 font-mono text-xs bg-white text-gray-900 outline-hidden"
            >
              {[1, 2, 4, 8, 12, 16, 24, 32, 64].map((num) => (
                <option key={num} value={num}>
                  {num}
                </option>
              ))}
            </select>
          </div>

          <label className="flex items-center gap-1.5 cursor-pointer">
            <input
              type="checkbox"
              checked={benchmark.multiThread}
              onChange={(e) => toggleBenchmarkThreads(e.target.checked)}
              className="text-blue-600 rounded-xs"
            />
            <span className="font-semibold text-gray-700">Multi-thread</span>
          </label>
        </div>

        {/* Hardware Throughput Live Visualizer Canvas matching Figure 4 */}
        <div className="p-3 bg-white">
          <div className="rounded-xs overflow-hidden border border-slate-700 shadow-inner">
            <canvas ref={canvasRef} width={460} height={110} className="w-full block" />
          </div>
        </div>

        {/* Metrics Grid matching Figure 4 */}
        <div className="px-5 py-2 bg-white grid grid-cols-2 gap-x-8 gap-y-2 text-xs">
          <div className="flex justify-between items-center py-0.5 border-b border-gray-100">
            <span className="text-gray-600">Current speed</span>
            <span className="font-mono font-bold text-gray-900 text-sm">
              {formatNumberWithCommas(benchmark.currentSpeedKbps)} KB/s
            </span>
          </div>

          <div className="flex justify-between items-center py-0.5 border-b border-gray-100">
            <span className="text-gray-600">Elapsed</span>
            <span className="font-mono font-semibold text-gray-900">
              {formatTime(benchmark.elapsedSeconds)}
            </span>
          </div>

          <div className="flex justify-between items-center py-0.5 border-b border-gray-100">
            <span className="text-gray-600">Average speed</span>
            <span className="font-mono font-bold text-gray-900">
              {formatNumberWithCommas(benchmark.averageSpeedKbps)} KB/s
            </span>
          </div>

          <div className="flex justify-between items-center py-0.5 border-b border-gray-100">
            <span className="text-gray-600">Memory errors</span>
            <span className="font-mono font-semibold text-emerald-600">
              {benchmark.memoryErrors}
            </span>
          </div>

          <div className="flex justify-between items-center py-0.5 border-b border-gray-100">
            <span className="text-gray-600">Processed</span>
            <span className="font-mono font-semibold text-gray-900">
              {formatNumberWithCommas(benchmark.processedMb)} MB
            </span>
          </div>

          <div className="flex justify-between items-center py-0.5 border-b border-gray-100">
            <span className="text-gray-600">Status</span>
            <span className="font-mono font-bold text-emerald-600">
              {benchmark.status}
            </span>
          </div>
        </div>

        {/* Action Buttons matching Figure 4 */}
        <div className="h-12 bg-[#f0f0f0] border-t border-[#dcdcdc] flex items-center justify-end px-3 gap-2 mt-1">
          <button
            onClick={copyResults}
            className="winrar-btn px-3 py-1 text-xs text-gray-800"
          >
            Copy to clipboard
          </button>
          <button
            onClick={() => {
              if (benchmark.running) {
                stopBenchmark();
              } else {
                startBenchmark(benchmark.threads);
              }
            }}
            className="winrar-btn w-20 py-1 text-xs font-semibold text-gray-800"
          >
            {benchmark.running ? 'Stop' : 'Start'}
          </button>
          <button
            onClick={closeModal}
            className="winrar-btn w-20 py-1 text-xs text-gray-800"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};
