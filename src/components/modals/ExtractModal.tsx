import React, { useState } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { UpdateMode, OverwriteMode } from '../../types';
import { Folder, FolderOpen, HelpCircle, X } from 'lucide-react';

export const ExtractModal: React.FC = () => {
  const { activeModal, closeModal, extractionConfig, updateExtractionConfig, startExtraction } = useArchiveStore();
  const [activeTab, setActiveTab] = useState<'general' | 'advanced'>('general');

  if (activeModal !== 'extract') return null;

  const folderBrowserTree = [
    { name: 'Desktop', path: 'C:\\Users\\Administrator\\Desktop' },
    { name: 'Documents', path: 'C:\\Users\\Administrator\\Documents' },
    { name: 'Downloads', path: 'C:\\Users\\Administrator\\Downloads' },
    { name: 'D: Data', path: 'D:\\' },
    { name: 'Projects', path: 'D:\\Projects' },
    { name: 'extracted', path: 'D:\\Projects\\extracted\\project_backup' }
  ];

  const handleOk = () => {
    startExtraction(extractionConfig);
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[520px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Modal Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>Extraction path and options</span>
          <div className="flex items-center gap-1">
            <button className="hover:bg-blue-700 p-0.5 rounded-xs" title="Help">
              <HelpCircle size={14} />
            </button>
            <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs" title="Close">
              <X size={14} />
            </button>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-[#a0a0a0] px-2 pt-2 bg-[#ece9d8]">
          <button
            onClick={() => setActiveTab('general')}
            className={`px-4 py-1 text-xs border-t border-l border-r rounded-t-xs -mb-[1px] ${
              activeTab === 'general'
                ? 'bg-white border-[#707070] font-semibold text-gray-900 border-b-transparent z-10'
                : 'bg-[#ece9d8] border-[#a0a0a0] text-gray-600 hover:text-black'
            }`}
          >
            General
          </button>
          <button
            onClick={() => setActiveTab('advanced')}
            className={`px-4 py-1 text-xs border-t border-l border-r rounded-t-xs -mb-[1px] ml-1 ${
              activeTab === 'advanced'
                ? 'bg-white border-[#707070] font-semibold text-gray-900 border-b-transparent z-10'
                : 'bg-[#ece9d8] border-[#a0a0a0] text-gray-600 hover:text-black'
            }`}
          >
            Advanced
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-3 bg-white flex flex-col gap-3 min-h-[360px]">
          {activeTab === 'general' ? (
            <>
              {/* Destination Path Box & Tree */}
              <div>
                <label className="text-[11px] text-gray-600 block mb-1">
                  Destination path (will be created if it does not exist)
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={extractionConfig.destinationPath}
                    onChange={(e) => updateExtractionConfig({ destinationPath: e.target.value })}
                    className="flex-1 bg-white border border-[#7f9db9] rounded-xs px-2 py-1 font-mono text-xs text-gray-900"
                  />
                </div>
              </div>

              {/* Main Layout: Left Options & Right Directory Tree */}
              <div className="flex gap-3">
                {/* Left Options Column */}
                <div className="w-[58%] flex flex-col gap-2">
                  {/* Update Mode Group */}
                  <fieldset className="border border-[#d0d0d0] p-2 rounded-xs">
                    <legend className="px-1 text-[11px] font-semibold text-gray-700">Update mode</legend>
                    <div className="flex flex-col gap-1.5 pt-0.5">
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="radio"
                          name="updateMode"
                          checked={extractionConfig.updateMode === 'extract_replace'}
                          onChange={() => updateExtractionConfig({ updateMode: 'extract_replace' })}
                          className="text-blue-600"
                        />
                        <span>Extract and replace files</span>
                      </label>
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="radio"
                          name="updateMode"
                          checked={extractionConfig.updateMode === 'extract_update'}
                          onChange={() => updateExtractionConfig({ updateMode: 'extract_update' })}
                          className="text-blue-600"
                        />
                        <span>Extract and update files</span>
                      </label>
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="radio"
                          name="updateMode"
                          checked={extractionConfig.updateMode === 'fresh_only'}
                          onChange={() => updateExtractionConfig({ updateMode: 'fresh_only' })}
                          className="text-blue-600"
                        />
                        <span>Fresh existing files only</span>
                      </label>
                    </div>
                  </fieldset>

                  {/* Overwrite Mode Group */}
                  <fieldset className="border border-[#d0d0d0] p-2 rounded-xs">
                    <legend className="px-1 text-[11px] font-semibold text-gray-700">Overwrite mode</legend>
                    <div className="flex flex-col gap-1.5 pt-0.5">
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="radio"
                          name="overwriteMode"
                          checked={extractionConfig.overwriteMode === 'ask'}
                          onChange={() => updateExtractionConfig({ overwriteMode: 'ask' })}
                          className="text-blue-600"
                        />
                        <span>Ask before overwrite</span>
                      </label>
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="radio"
                          name="overwriteMode"
                          checked={extractionConfig.overwriteMode === 'overwrite'}
                          onChange={() => updateExtractionConfig({ overwriteMode: 'overwrite' })}
                          className="text-blue-600"
                        />
                        <span>Overwrite without prompt</span>
                      </label>
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="radio"
                          name="overwriteMode"
                          checked={extractionConfig.overwriteMode === 'skip'}
                          onChange={() => updateExtractionConfig({ overwriteMode: 'skip' })}
                          className="text-blue-600"
                        />
                        <span>Skip existing files</span>
                      </label>
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="radio"
                          name="overwriteMode"
                          checked={extractionConfig.overwriteMode === 'rename_auto'}
                          onChange={() => updateExtractionConfig({ overwriteMode: 'rename_auto' })}
                          className="text-blue-600"
                        />
                        <span>Rename automatically</span>
                      </label>
                    </div>
                  </fieldset>

                  {/* Miscellaneous Group */}
                  <fieldset className="border border-[#d0d0d0] p-2 rounded-xs">
                    <legend className="px-1 text-[11px] font-semibold text-gray-700">Miscellaneous</legend>
                    <div className="flex flex-col gap-1.5 pt-0.5">
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={extractionConfig.keepBrokenFiles}
                          onChange={(e) => updateExtractionConfig({ keepBrokenFiles: e.target.checked })}
                          className="text-blue-600 rounded-xs"
                        />
                        <span>Keep broken files</span>
                      </label>
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={extractionConfig.extractToSeparateFolder}
                          onChange={(e) => updateExtractionConfig({ extractToSeparateFolder: e.target.checked })}
                          className="text-blue-600 rounded-xs"
                        />
                        <span>Extract each archive to separate folder</span>
                      </label>
                      <label className="flex items-center gap-1.5 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={extractionConfig.openFolderAfterExtraction}
                          onChange={(e) => updateExtractionConfig({ openFolderAfterExtraction: e.target.checked })}
                          className="text-blue-600 rounded-xs"
                        />
                        <span>Display files in explorer</span>
                      </label>
                    </div>
                  </fieldset>
                </div>

                {/* Right Folder Hierarchy Browser */}
                <div className="w-[42%] border border-[#7f9db9] rounded-xs p-1.5 bg-[#fbfbfb] flex flex-col h-[280px] overflow-y-auto">
                  <div className="text-[11px] font-bold text-gray-600 mb-1 border-b pb-0.5">
                    Target Destination
                  </div>
                  <div className="flex flex-col gap-1">
                    {folderBrowserTree.map((item, idx) => (
                      <div
                        key={idx}
                        onClick={() => updateExtractionConfig({ destinationPath: item.path })}
                        className={`flex items-center gap-1.5 px-1 py-1 rounded-xs cursor-pointer text-xs ${
                          extractionConfig.destinationPath === item.path
                            ? 'bg-blue-600 text-white font-medium'
                            : 'hover:bg-blue-50 text-gray-700'
                        }`}
                        style={{ paddingLeft: `${idx * 6 + 4}px` }}
                      >
                        {extractionConfig.destinationPath === item.path ? (
                          <FolderOpen size={14} className="text-yellow-300" />
                        ) : (
                          <Folder size={14} className="text-[#f57f17] fill-[#ffd54f]" />
                        )}
                        <span className="truncate">{item.name}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </>
          ) : (
            <div className="flex flex-col gap-3 p-2">
              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">Security & File System</legend>
                <div className="flex flex-col gap-2">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" defaultChecked className="text-blue-600" />
                    <span>Restore NTFS Security permissions (ACLs)</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" defaultChecked className="text-blue-600" />
                    <span>Preserve high precision UTC timestamps (100ns)</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" className="text-blue-600" />
                    <span>Delete archive file automatically after extraction</span>
                  </label>
                </div>
              </fieldset>

              <fieldset className="border border-[#d0d0d0] p-2.5 rounded-xs">
                <legend className="px-1 text-[11px] font-semibold text-gray-700">File Links & Streams</legend>
                <div className="flex flex-col gap-2">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" defaultChecked className="text-blue-600" />
                    <span>Extract hard links and symbolic links</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" defaultChecked className="text-blue-600" />
                    <span>Propagate Mark of the Web (MOTW) zone info</span>
                  </label>
                </div>
              </fieldset>
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
          <button
            onClick={() => alert('Extraction options help documentation.')}
            className="winrar-btn w-20 py-1 text-xs text-gray-800"
          >
            Help
          </button>
        </div>
      </div>
    </div>
  );
};
