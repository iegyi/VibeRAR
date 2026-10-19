import React from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { X, ExternalLink, Globe, User, Heart } from 'lucide-react';

const GithubIcon: React.FC<{ size?: number; className?: string }> = ({ size = 15, className = '' }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
    <path d="M9 18c-4.51 2-5-2-7-2" />
  </svg>
);

export const AboutModal: React.FC = () => {
  const { activeModal, closeModal } = useArchiveStore();

  if (activeModal !== 'about') return null;

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[460px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>About VibeRAR</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs" title="Close">
            <X size={14} />
          </button>
        </div>

        {/* Content Body */}
        <div className="p-5 bg-white flex flex-col gap-4">
          <div className="flex items-center gap-4">
            <img
              src="/viberar_icon.png"
              alt="VibeRAR 3D Icon"
              className="w-16 h-16 rounded-xl shadow-md shrink-0 object-contain"
            />
            <div className="flex-1">
              <h2 className="text-lg font-bold text-gray-900 tracking-tight">
                VibeRAR v1.1.0
              </h2>
              <p className="text-xs text-gray-600 font-medium mt-0.5">
                A fast, modern desktop archiver for Windows 11 & 10.
              </p>
            </div>
          </div>

          <div className="border-t border-gray-200 pt-3 flex flex-col gap-2.5 text-xs">
            {/* Created by */}
            <div className="flex items-center gap-2 text-gray-800">
              <User size={15} className="text-blue-600 shrink-0" />
              <span>Created by <strong className="text-gray-900 font-semibold">Ahmed Samir</strong></span>
            </div>

            {/* Official Site */}
            <div className="flex items-center gap-2">
              <Globe size={15} className="text-blue-600 shrink-0" />
              <span className="text-gray-700">Official Site:</span>
              <a
                href="https://ahmedsamir.com"
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:text-blue-800 hover:underline font-semibold flex items-center gap-1 cursor-pointer"
              >
                ahmedsamir.com
                <ExternalLink size={11} />
              </a>
            </div>

            {/* GitHub */}
            <div className="flex items-center gap-2">
              <GithubIcon size={15} className="text-gray-800 shrink-0" />
              <span className="text-gray-700">GitHub:</span>
              <a
                href="https://github.com/iegyi/VibeRAR"
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:text-blue-800 hover:underline font-semibold flex items-center gap-1 cursor-pointer"
              >
                github.com/iegyi/VibeRAR
                <ExternalLink size={11} />
              </a>
            </div>

            {/* Sponsored by */}
            <div className="flex items-center gap-2 pt-1 border-t border-gray-100">
              <Heart size={14} className="text-pink-600 shrink-0 fill-pink-500/20" />
              <span className="text-gray-700">sponsored by</span>
              <a
                href="https://www.0webhost.com/"
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-700 hover:text-blue-900 hover:underline font-bold flex items-center gap-1 cursor-pointer"
              >
                0webhost.com
                <ExternalLink size={11} />
              </a>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="h-11 bg-[#f0f0f0] border-t border-[#dcdcdc] flex items-center justify-between px-3">
          <button
            onClick={() => {
              const { openModal } = useArchiveStore.getState();
              openModal('update');
            }}
            className="winrar-btn px-3 py-1 text-xs text-blue-700 hover:text-blue-900 font-semibold cursor-pointer"
          >
            Check for updates...
          </button>
          <button
            onClick={closeModal}
            className="winrar-btn w-20 py-1 text-xs font-semibold text-gray-800 cursor-pointer"
          >
            OK
          </button>
        </div>
      </div>
    </div>
  );
};
