import React, { useEffect, useState } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import {
  X,
  RefreshCw,
  ExternalLink,
  CheckCircle,
  AlertTriangle,
  Download,
  ShieldCheck,
  Sparkles,
  ArrowRight,
  Wrench
} from 'lucide-react';

interface UpdateInfo {
  current_version: string;
  latest_version: string;
  has_update: boolean;
  release_name: string;
  release_notes: string;
  download_url: string;
  html_url: string;
  asset_name: string;
  asset_size: number;
}

export const UpdateModal: React.FC = () => {
  const { activeModal, closeModal } = useArchiveStore();
  const [checking, setChecking] = useState(false);
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Download state
  const [downloading, setDownloading] = useState(false);
  const [downloadProgress, setDownloadProgress] = useState(0);
  const [downloadBytes, setDownloadBytes] = useState(0);
  const [downloadTotal, setDownloadTotal] = useState(0);
  const [downloadError, setDownloadError] = useState<string | null>(null);

  // Shell repair state
  const [repairing, setRepairing] = useState(false);
  const [repairMsg, setRepairMsg] = useState<string | null>(null);

  const checkForUpdates = async () => {
    setChecking(true);
    setError(null);
    try {
      const res = await fetch('/api/check-update');
      if (!res.ok) {
        throw new Error(`Server returned HTTP ${res.status}`);
      }
      const data: UpdateInfo = await res.json();
      setUpdateInfo(data);
    } catch (e: any) {
      setError(e.message || 'Failed to check for updates');
    } finally {
      setChecking(false);
    }
  };

  useEffect(() => {
    if (activeModal === 'update') {
      checkForUpdates();
    } else {
      // Reset state on close
      setDownloading(false);
      setRepairMsg(null);
    }
  }, [activeModal]);

  const handleStartUpdate = async () => {
    if (!updateInfo || !updateInfo.download_url) return;
    setDownloading(true);
    setDownloadError(null);
    setDownloadProgress(0);

    try {
      const res = await fetch('/api/download-update', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ download_url: updateInfo.download_url })
      });

      if (!res.ok) {
        throw new Error('Failed to initiate update download');
      }

      // Poll progress
      const pollTimer = setInterval(async () => {
        try {
          const statusRes = await fetch('/api/update-status');
          if (statusRes.ok) {
            const statusData = await statusRes.json();
            setDownloadProgress(statusData.percent || 0);
            setDownloadBytes(statusData.bytes || 0);
            setDownloadTotal(statusData.total || 0);

            if (statusData.status === 'done') {
              clearInterval(pollTimer);
            } else if (statusData.status === 'error') {
              clearInterval(pollTimer);
              setDownloadError(statusData.error || 'Download failed');
              setDownloading(false);
            }
          }
        } catch {
          // Ignore polling errors
        }
      }, 300);
    } catch (e: any) {
      setDownloadError(e.message || 'Error downloading update');
      setDownloading(false);
    }
  };

  const handleRepairWin11Menu = async () => {
    setRepairing(true);
    setRepairMsg(null);
    try {
      const res = await fetch('/api/repair-win11-menu', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshExplorer: true })
      });
      const data = await res.json();
      setRepairMsg(data.message || 'Windows 11 modern context menu repaired!');
    } catch (e: any) {
      setRepairMsg('Failed to repair menu: ' + e.message);
    } finally {
      setRepairing(false);
    }
  };

  if (activeModal !== 'update') return null;

  const formatSize = (bytes: number) => {
    if (!bytes || bytes === 0) return '';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-xs flex items-center justify-center z-50 select-none animate-in fade-in duration-150">
      <div className="w-[520px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <div className="flex items-center gap-1.5">
            <Sparkles size={13} className="text-yellow-300" />
            <span>VibeRAR - Software Update & Shell Integrator</span>
          </div>
          <button
            onClick={closeModal}
            className="hover:bg-red-600 p-0.5 rounded-xs transition-colors cursor-pointer"
            title="Close"
          >
            <X size={14} />
          </button>
        </div>

        {/* Content Body */}
        <div className="p-5 bg-white flex flex-col gap-4">
          {/* Header Banner */}
          <div className="flex items-center gap-3.5 pb-3 border-b border-gray-200">
            <img
              src="/viberar_icon.png"
              alt="VibeRAR Logo"
              className="w-14 h-14 rounded-lg shadow-sm shrink-0 object-contain"
            />
            <div className="flex-1">
              <h2 className="text-base font-bold text-gray-900 tracking-tight flex items-center gap-2">
                VibeRAR Update Selection
                <span className="text-xs px-2 py-0.5 bg-blue-100 text-blue-700 rounded-full font-mono font-medium">
                  {updateInfo?.current_version || 'v1.1.0'}
                </span>
              </h2>
              <p className="text-xs text-gray-600 mt-0.5">
                Keep your VibeRAR archiver and Windows 11 shell integrations always up to date.
              </p>
            </div>
          </div>

          {/* Status Section */}
          {checking && (
            <div className="py-8 flex flex-col items-center justify-center gap-3 text-gray-600">
              <RefreshCw size={24} className="animate-spin text-blue-600" />
              <div className="text-sm font-medium">Checking GitHub for latest VibeRAR release...</div>
            </div>
          )}

          {!checking && error && (
            <div className="p-3 bg-red-50 border border-red-200 rounded text-red-700 flex items-start gap-2.5">
              <AlertTriangle size={16} className="text-red-500 shrink-0 mt-0.5" />
              <div className="flex-1">
                <div className="font-semibold text-xs">Failed to check updates</div>
                <div className="text-xs text-red-600 mt-0.5">{error}</div>
                <button
                  onClick={checkForUpdates}
                  className="mt-2 text-xs font-semibold text-blue-600 hover:underline flex items-center gap-1 cursor-pointer"
                >
                  <RefreshCw size={11} /> Retry Check
                </button>
              </div>
            </div>
          )}

          {!checking && !error && updateInfo && (
            <>
              {/* Case 1: Update Available */}
              {updateInfo.has_update ? (
                <div className="flex flex-col gap-3">
                  <div className="p-3 bg-green-50 border border-green-200 rounded flex items-center justify-between">
                    <div className="flex items-center gap-2.5">
                      <div className="w-8 h-8 rounded-full bg-green-100 text-green-700 flex items-center justify-center font-bold">
                        <Download size={16} />
                      </div>
                      <div>
                        <div className="font-bold text-green-900 text-sm">
                          New Version Available: {updateInfo.latest_version}
                        </div>
                        <div className="text-xs text-green-700 flex items-center gap-1.5 mt-0.5">
                          <span>Current: {updateInfo.current_version}</span>
                          <ArrowRight size={11} />
                          <span className="font-semibold">New: {updateInfo.latest_version}</span>
                          {updateInfo.asset_size > 0 && (
                            <span>• {formatSize(updateInfo.asset_size)}</span>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Release Notes */}
                  <div>
                    <div className="text-xs font-semibold text-gray-700 mb-1 flex items-center justify-between">
                      <span>What's New in {updateInfo.latest_version}:</span>
                      {updateInfo.html_url && (
                        <a
                          href={updateInfo.html_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-600 hover:underline flex items-center gap-1 font-normal"
                        >
                          View release on GitHub <ExternalLink size={10} />
                        </a>
                      )}
                    </div>
                    <div className="p-2.5 bg-gray-50 border border-gray-200 rounded font-mono text-[11px] text-gray-700 max-h-32 overflow-y-auto whitespace-pre-wrap leading-relaxed">
                      {updateInfo.release_notes || 'No release notes provided.'}
                    </div>
                  </div>

                  {/* Downloading Progress */}
                  {downloading && (
                    <div className="p-3 bg-blue-50 border border-blue-200 rounded flex flex-col gap-2">
                      <div className="flex justify-between text-xs font-semibold text-blue-900">
                        <span>Downloading VibeRAR Update Setup...</span>
                        <span>{downloadProgress}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-3 overflow-hidden border border-gray-300">
                        <div
                          className="bg-blue-600 h-full transition-all duration-200"
                          style={{ width: `${downloadProgress}%` }}
                        />
                      </div>
                      <div className="text-[11px] text-blue-700 flex justify-between">
                        <span>
                          {formatSize(downloadBytes)} of {formatSize(downloadTotal)}
                        </span>
                        <span>Installer will launch automatically</span>
                      </div>
                    </div>
                  )}

                  {downloadError && (
                    <div className="text-xs text-red-600 font-medium">
                      Download error: {downloadError}
                    </div>
                  )}
                </div>
              ) : (
                /* Case 2: Already Up To Date */
                <div className="p-4 bg-emerald-50/70 border border-emerald-200 rounded flex items-center gap-3">
                  <div className="w-9 h-9 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center shrink-0">
                    <CheckCircle size={20} />
                  </div>
                  <div className="flex-1">
                    <div className="font-bold text-gray-900 text-sm">
                      You are using the latest version of VibeRAR!
                    </div>
                    <div className="text-xs text-gray-600 mt-0.5">
                      Installed version <strong className="text-gray-800">{updateInfo.current_version}</strong> is the most recent release.
                    </div>
                  </div>
                </div>
              )}
            </>
          )}

          {/* Quick Windows 11 Context Menu Self-Repair Section */}
          <div className="p-3 bg-gray-50 border border-gray-200 rounded flex flex-col gap-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 text-gray-800 font-semibold text-xs">
                <Wrench size={13} className="text-blue-600" />
                <span>Windows 11 Modern Context Menu Repair</span>
              </div>
              <button
                onClick={handleRepairWin11Menu}
                disabled={repairing}
                className="px-2.5 py-1 bg-white hover:bg-gray-100 border border-gray-300 rounded text-xs font-medium text-gray-700 shadow-2xs transition-colors cursor-pointer disabled:opacity-50 flex items-center gap-1.5"
              >
                {repairing ? (
                  <>
                    <RefreshCw size={11} className="animate-spin text-blue-600" />
                    <span>Repairing...</span>
                  </>
                ) : (
                  <>
                    <ShieldCheck size={12} className="text-green-600" />
                    <span>Repair / Register Menu</span>
                  </>
                )}
              </button>
            </div>
            <p className="text-[11px] text-gray-500 leading-tight">
              Instantly ensures COM CLSID registration and Windows 11 Sparse Package integration for this PC.
            </p>
            {repairMsg && (
              <div className="text-[11px] font-medium text-green-700 bg-green-50 p-1.5 rounded border border-green-200">
                {repairMsg}
              </div>
            )}
          </div>
        </div>

        {/* Modal Footer Controls */}
        <div className="px-5 py-3 bg-[#e8e8e8] border-t border-[#d8d8d8] flex items-center justify-between">
          <button
            onClick={checkForUpdates}
            disabled={checking || downloading}
            className="px-3 py-1 bg-white hover:bg-gray-100 border border-gray-300 rounded text-xs text-gray-700 shadow-2xs transition-colors cursor-pointer flex items-center gap-1.5 disabled:opacity-50"
          >
            <RefreshCw size={12} className={checking ? 'animate-spin' : ''} />
            <span>Check Again</span>
          </button>

          <div className="flex items-center gap-2">
            {updateInfo?.has_update && (
              <button
                onClick={handleStartUpdate}
                disabled={downloading}
                className="px-4 py-1.5 bg-[#0055ea] hover:bg-blue-700 text-white font-semibold rounded text-xs shadow-sm transition-colors cursor-pointer flex items-center gap-1.5 disabled:opacity-50"
              >
                <Download size={13} />
                <span>{downloading ? 'Downloading...' : 'Install Update Now'}</span>
              </button>
            )}

            <button
              onClick={closeModal}
              className="px-4 py-1.5 bg-white hover:bg-gray-100 border border-gray-300 rounded text-xs font-medium text-gray-700 shadow-2xs transition-colors cursor-pointer"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
