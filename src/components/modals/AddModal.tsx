import React, { useState } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { CompressionFormat, CompressionLevel } from '../../types';
import { X, Lock, KeyRound } from 'lucide-react';

export const AddModal: React.FC = () => {
  const { activeModal, closeModal, compressionConfig, updateCompressionConfig, startCompression } = useArchiveStore();
  const [activeTab, setActiveTab] = useState<'general' | 'advanced' | 'options' | 'files'>('general');
  const [showPasswordInput, setShowPasswordInput] = useState(false);

  if (activeModal !== 'add') return null;

  const handleOk = () => {
    startCompression(compressionConfig);
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[520px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>Archive name and parameters</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-[#a0a0a0] px-2 pt-2 bg-[#ece9d8]">
          {['General', 'Advanced', 'Options', 'Files'].map((t) => (
            <button
              key={t}
              onClick={() => setActiveTab(t.toLowerCase() as any)}
              className={`px-3 py-1 text-xs border-t border-l border-r rounded-t-xs -mb-[1px] mr-1 ${
                activeTab === t.toLowerCase()
                  ? 'bg-white border-[#707070] font-semibold text-gray-900 border-b-transparent z-10'
                  : 'bg-[#ece9d8] border-[#a0a0a0] text-gray-600 hover:text-black'
              }`}
            >
              {t}
            </button>
          ))}
        </div>

        {/* Tab Body */}
        <div className="p-3 bg-white flex flex-col gap-3 min-h-[340px]">
          {activeTab === 'general' ? (
            <>
              {/* Archive Name */}
              <div>
                <label className="text-[11px] text-gray-600 block mb-1">Archive name</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={compressionConfig.archiveName}
                    onChange={(e) => updateCompressionConfig({ archiveName: e.target.value })}
                    className="flex-1 border border-[#7f9db9] rounded-xs px-2 py-1 font-mono text-xs text-gray-900"
                  />
                  <button className="winrar-btn px-3 text-xs">Browse...</button>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                {/* Archive Format */}
                <fieldset className="border border-[#d0d0d0] p-2 rounded-xs">
                  <legend className="px-1 text-[11px] font-semibold text-gray-700">Archive format</legend>
                  <div className="flex flex-col gap-1.5 pt-0.5">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="archFormat"
                        checked={compressionConfig.format === 'rar'}
                        onChange={() => updateCompressionConfig({ format: 'rar', archiveName: compressionConfig.archiveName.replace(/\.[^.]+$/, '.rar') })}
                        className="text-blue-600"
                      />
                      <span className="font-semibold text-blue-900">RAR (RAR5 / WinRAR format)</span>
                    </label>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="archFormat"
                        checked={compressionConfig.format === 'zip'}
                        onChange={() => updateCompressionConfig({ format: 'zip', archiveName: compressionConfig.archiveName.replace(/\.[^.]+$/, '.zip') })}
                        className="text-blue-600"
                      />
                      <span className="font-semibold">ZIP (Zip64)</span>
                    </label>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="archFormat"
                        checked={compressionConfig.format === '7z'}
                        onChange={() => updateCompressionConfig({ format: '7z', archiveName: compressionConfig.archiveName.replace(/\.[^.]+$/, '.7z') })}
                        className="text-blue-600"
                      />
                      <span>7Z (LZMA2 / High Ratio)</span>
                    </label>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="archFormat"
                        checked={compressionConfig.format === 'tar'}
                        onChange={() => updateCompressionConfig({ format: 'tar', archiveName: compressionConfig.archiveName.replace(/\.[^.]+$/, '.tar') })}
                        className="text-blue-600"
                      />
                      <span>TAR (POSIX archive)</span>
                    </label>
                  </div>
                </fieldset>

                {/* Compression Method */}
                <fieldset className="border border-[#d0d0d0] p-2 rounded-xs">
                  <legend className="px-1 text-[11px] font-semibold text-gray-700">Compression method</legend>
                  <select
                    value={compressionConfig.compressionLevel}
                    onChange={(e) => updateCompressionConfig({ compressionLevel: e.target.value as CompressionLevel })}
                    className="w-full border border-[#7f9db9] rounded-xs px-2 py-1 text-xs bg-white text-gray-900 mt-1"
                  >
                    <option value="store">Store (No compression)</option>
                    <option value="fastest">Fastest</option>
                    <option value="fast">Fast</option>
                    <option value="normal">Normal (Recommended)</option>
                    <option value="good">Good</option>
                    <option value="best">Best (Maximum ratio)</option>
                  </select>

                  <div className="mt-2 text-[11px] text-gray-600">Dictionary size:</div>
                  <select
                    value={compressionConfig.dictionarySize}
                    onChange={(e) => updateCompressionConfig({ dictionarySize: e.target.value })}
                    className="w-full border border-[#7f9db9] rounded-xs px-2 py-1 text-xs bg-white text-gray-900 mt-0.5"
                  >
                    <option value="4 MB">4 MB</option>
                    <option value="32 MB">32 MB (Standard)</option>
                    <option value="64 MB">64 MB</option>
                    <option value="128 MB">128 MB</option>
                    <option value="1 GB">1 GB (Ultra Long Range)</option>
                  </select>
                </fieldset>
              </div>

              {/* Archiving Options */}
              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Archiving options</legend>
                <div className="grid grid-cols-2 gap-2">
                  <label className="flex items-center gap-1.5 cursor-pointer">
                    <input type="checkbox" defaultChecked className="text-blue-600 rounded-xs" />
                    <span>Create solid archive</span>
                  </label>
                  <label className="flex items-center gap-1.5 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={compressionConfig.createSfx}
                      onChange={(e) => updateCompressionConfig({ createSfx: e.target.checked })}
                      className="text-blue-600 rounded-xs"
                    />
                    <span className="font-semibold text-blue-800">Create SFX archive (.exe)</span>
                  </label>
                  <label className="flex items-center gap-1.5 cursor-pointer">
                    <input type="checkbox" defaultChecked className="text-blue-600 rounded-xs" />
                    <span>Put recovery record ({compressionConfig.recoveryRecordPercent}%)</span>
                  </label>
                  <label className="flex items-center gap-1.5 cursor-pointer">
                    <input type="checkbox" defaultChecked className="text-blue-600 rounded-xs" />
                    <span>Test archived files</span>
                  </label>
                </div>
              </fieldset>

              {/* Set Password Button / Section */}
              <div className="flex items-center justify-between p-2 bg-[#f8f9fa] border border-[#e0e0e0] rounded-xs">
                <div className="flex items-center gap-2">
                  <Lock size={15} className="text-amber-600" />
                  <span className="font-semibold text-gray-700">
                    {compressionConfig.password ? 'Password protection enabled' : 'Password protection'}
                  </span>
                </div>
                <button
                  onClick={() => setShowPasswordInput(!showPasswordInput)}
                  className="winrar-btn px-3 py-1 text-xs flex items-center gap-1"
                >
                  <KeyRound size={12} />
                  <span>{showPasswordInput ? 'Hide' : 'Set password...'}</span>
                </button>
              </div>

              {showPasswordInput && (
                <div className="p-2 border border-amber-200 bg-amber-50 rounded-xs flex flex-col gap-2">
                  <div className="flex gap-2 items-center">
                    <span className="text-[11px] text-gray-700 w-24">Enter password:</span>
                    <input
                      type="password"
                      value={compressionConfig.password || ''}
                      onChange={(e) => updateCompressionConfig({ password: e.target.value })}
                      className="flex-1 border border-gray-300 rounded-xs px-2 py-0.5 text-xs bg-white"
                      placeholder="AES-256 encryption password"
                    />
                  </div>
                  <label className="flex items-center gap-2 cursor-pointer ml-24">
                    <input
                      type="checkbox"
                      checked={compressionConfig.encryptFileNames}
                      onChange={(e) => updateCompressionConfig({ encryptFileNames: e.target.checked })}
                      className="text-blue-600"
                    />
                    <span className="text-[11px] text-gray-700">Encrypt file names (-hp mode)</span>
                  </label>
                </div>
              )}
            </>
          ) : (
            <div className="p-2 text-gray-600">
              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Volumes & Splitting</legend>
                <div className="flex flex-col gap-2">
                  <label className="text-[11px]">Split to volumes, size:</label>
                  <select className="border border-[#7f9db9] rounded-xs px-2 py-1 text-xs bg-white w-48">
                    <option>Autodetect</option>
                    <option>5 MB</option>
                    <option>100 MB</option>
                    <option>650 MB (CD)</option>
                    <option>4480 MB (DVD)</option>
                  </select>
                </div>
              </fieldset>
            </div>
          )}
        </div>

        {/* Action Buttons */}
        <div className="h-11 bg-[#f0f0f0] border-t border-[#dcdcdc] flex items-center justify-end px-3 gap-2">
          <button
            onClick={handleOk}
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
          <button
            onClick={() => alert('Compression parameters help.')}
            className="winrar-btn w-20 py-1 text-xs text-gray-800"
          >
            Help
          </button>
        </div>
      </div>
    </div>
  );
};
