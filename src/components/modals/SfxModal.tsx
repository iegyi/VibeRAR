import React, { useState } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { X, HelpCircle } from 'lucide-react';

export const SfxModal: React.FC = () => {
  const { activeModal, closeModal, sfxConfig, updateSfxConfig } = useArchiveStore();
  const [activeTab, setActiveTab] = useState<'general' | 'setup' | 'modes' | 'text' | 'license' | 'advanced'>('modes');

  if (activeModal !== 'sfx') return null;

  const tabs = [
    { id: 'general', label: 'General' },
    { id: 'setup', label: 'Setup' },
    { id: 'modes', label: 'Modes' },
    { id: 'text', label: 'Text and icon' },
    { id: 'license', label: 'License' },
    { id: 'advanced', label: 'Advanced' }
  ] as const;

  const handleOk = () => {
    alert(`SFX Executable built successfully!\nPayload stub concatenated with marker and SFX config.`);
    closeModal();
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[500px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Modal Title */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>Advanced SFX options</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-[#a0a0a0] px-2 pt-2 bg-[#ece9d8] overflow-x-auto">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-3 py-1 text-xs border-t border-l border-r rounded-t-xs -mb-[1px] mr-1 whitespace-nowrap ${
                activeTab === tab.id
                  ? 'bg-white border-[#707070] font-semibold text-gray-900 border-b-transparent z-10'
                  : 'bg-[#ece9d8] border-[#a0a0a0] text-gray-600 hover:text-black'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Tab Body */}
        <div className="p-3 bg-white flex flex-col gap-3 min-h-[300px]">
          {/* General Tab */}
          {activeTab === 'general' && (
            <div className="flex flex-col gap-3">
              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Path to extract</legend>
                <div className="flex flex-col gap-2">
                  <input
                    type="text"
                    value={sfxConfig.destinationPath}
                    onChange={(e) => updateSfxConfig({ destinationPath: e.target.value })}
                    className="border border-[#7f9db9] rounded-xs px-2 py-1 font-mono text-xs w-full"
                  />
                  <label className="flex items-center gap-2 cursor-pointer pt-1">
                    <input
                      type="checkbox"
                      checked={sfxConfig.unpackToTemp}
                      onChange={(e) => updateSfxConfig({ unpackToTemp: e.target.checked })}
                      className="text-blue-600"
                    />
                    <span>Unpack to temporary folder (%TEMP%)</span>
                  </label>
                </div>
              </fieldset>
            </div>
          )}

          {/* Setup Tab */}
          {activeTab === 'setup' && (
            <div className="flex flex-col gap-3">
              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Run after extraction</legend>
                <div className="flex flex-col gap-2">
                  <label className="text-gray-600 text-[11px]">Run program:</label>
                  <input
                    type="text"
                    value={sfxConfig.setupExecutable}
                    onChange={(e) => updateSfxConfig({ setupExecutable: e.target.value })}
                    className="border border-[#7f9db9] rounded-xs px-2 py-1 font-mono text-xs w-full"
                    placeholder="e.g. setup.exe or installer.cmd"
                  />
                </div>
              </fieldset>
            </div>
          )}

          {/* Modes Tab matching Figure 3 */}
          {activeTab === 'modes' && (
            <div className="flex flex-col gap-3">
              {/* Silent mode */}
              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Silent mode</legend>
                <div className="flex flex-col gap-1.5 pt-0.5">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="silentMode"
                      checked={sfxConfig.silentMode === 'display_all'}
                      onChange={() => updateSfxConfig({ silentMode: 'display_all' })}
                      className="text-blue-600"
                    />
                    <span>Display all</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="silentMode"
                      checked={sfxConfig.silentMode === 'hide_start'}
                      onChange={() => updateSfxConfig({ silentMode: 'hide_start' })}
                      className="text-blue-600"
                    />
                    <span>Hide start dialog</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="silentMode"
                      checked={sfxConfig.silentMode === 'hide_all'}
                      onChange={() => updateSfxConfig({ silentMode: 'hide_all' })}
                      className="text-blue-600"
                    />
                    <span>Hide all (Silent extraction)</span>
                  </label>
                </div>
              </fieldset>

              {/* Overwrite mode */}
              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Overwrite mode</legend>
                <div className="flex flex-col gap-1.5 pt-0.5">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="sfxOverwrite"
                      checked={sfxConfig.overwriteMode === 'ask'}
                      onChange={() => updateSfxConfig({ overwriteMode: 'ask' })}
                      className="text-blue-600"
                    />
                    <span>Ask before overwrite</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="sfxOverwrite"
                      checked={sfxConfig.overwriteMode === 'overwrite'}
                      onChange={() => updateSfxConfig({ overwriteMode: 'overwrite' })}
                      className="text-blue-600"
                    />
                    <span>Overwrite all files</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="sfxOverwrite"
                      checked={sfxConfig.overwriteMode === 'skip'}
                      onChange={() => updateSfxConfig({ overwriteMode: 'skip' })}
                      className="text-blue-600"
                    />
                    <span>Skip existing files</span>
                  </label>
                </div>
              </fieldset>

              {/* Update mode checkboxes */}
              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Update mode</legend>
                <div className="flex flex-col gap-1.5 pt-0.5">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      defaultChecked
                      className="text-blue-600"
                    />
                    <span>Extract and replace files</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={sfxConfig.requestAdmin}
                      onChange={(e) => updateSfxConfig({ requestAdmin: e.target.checked })}
                      className="text-blue-600"
                    />
                    <span>Request administrative access (-iadm UAC prompt)</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={sfxConfig.unpackToTemp}
                      onChange={(e) => updateSfxConfig({ unpackToTemp: e.target.checked })}
                      className="text-blue-600"
                    />
                    <span>Unpack to temporary folder</span>
                  </label>
                </div>
              </fieldset>
            </div>
          )}

          {/* Text and Icon */}
          {activeTab === 'text' && (
            <div className="flex flex-col gap-3">
              <div>
                <label className="text-gray-600 text-[11px] block mb-1">Title of SFX window:</label>
                <input
                  type="text"
                  value={sfxConfig.title}
                  onChange={(e) => updateSfxConfig({ title: e.target.value })}
                  className="border border-[#7f9db9] rounded-xs px-2 py-1 text-xs w-full"
                />
              </div>
              <div>
                <label className="text-gray-600 text-[11px] block mb-1">Text to display in SFX window:</label>
                <textarea
                  rows={4}
                  value={sfxConfig.textPrompt}
                  onChange={(e) => updateSfxConfig({ textPrompt: e.target.value })}
                  className="border border-[#7f9db9] rounded-xs px-2 py-1 text-xs w-full"
                />
              </div>
            </div>
          )}

          {/* License */}
          {activeTab === 'license' && (
            <div className="flex flex-col gap-3">
              <div>
                <label className="text-gray-600 text-[11px] block mb-1">License dialog title:</label>
                <input
                  type="text"
                  value={sfxConfig.licenseTitle}
                  onChange={(e) => updateSfxConfig({ licenseTitle: e.target.value })}
                  className="border border-[#7f9db9] rounded-xs px-2 py-1 text-xs w-full"
                />
              </div>
              <div>
                <label className="text-gray-600 text-[11px] block mb-1">License text to display:</label>
                <textarea
                  rows={5}
                  value={sfxConfig.licenseText}
                  onChange={(e) => updateSfxConfig({ licenseText: e.target.value })}
                  className="border border-[#7f9db9] rounded-xs px-2 py-1 text-xs w-full font-mono"
                />
              </div>
            </div>
          )}

          {/* Advanced */}
          {activeTab === 'advanced' && (
            <div className="flex flex-col gap-3">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={sfxConfig.createDesktopShortcut}
                  onChange={(e) => updateSfxConfig({ createDesktopShortcut: e.target.checked })}
                  className="text-blue-600"
                />
                <span>Create shortcut on Desktop</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={sfxConfig.requestAdmin}
                  onChange={(e) => updateSfxConfig({ requestAdmin: e.target.checked })}
                  className="text-blue-600"
                />
                <span>Enforce high administrator privileges (UAC elevation)</span>
              </label>
            </div>
          )}
        </div>

        {/* Modal Action Buttons */}
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
        </div>
      </div>
    </div>
  );
};
