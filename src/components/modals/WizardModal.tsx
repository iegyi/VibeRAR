import React, { useState } from 'react';
import { useArchiveStore } from '../../store/useArchiveStore';
import { WinRarIcons } from '../../utils/icons';
import { X } from 'lucide-react';

export const WizardModal: React.FC = () => {
  const { activeModal, closeModal, openModal } = useArchiveStore();
  const [step, setStep] = useState<number>(1);
  const [action, setAction] = useState<'unpack' | 'pack'>('unpack');

  if (activeModal !== 'wizard') return null;

  const handleNext = () => {
    if (step === 1) {
      setStep(2);
    } else {
      closeModal();
      if (action === 'unpack') {
        openModal('extract');
      } else {
        openModal('add');
      }
    }
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 select-none">
      <div className="w-[490px] bg-[#f0f0f0] border border-[#707070] shadow-2xl rounded-xs flex flex-col text-xs text-gray-800 winrar-window">
        {/* Title Bar */}
        <div className="h-7 bg-[#0055ea] text-white flex items-center justify-between px-2 font-semibold">
          <span>Wizard</span>
          <button onClick={closeModal} className="hover:bg-red-600 p-0.5 rounded-xs">
            <X size={14} />
          </button>
        </div>

        {/* Content Body */}
        <div className="p-4 bg-white flex gap-4 min-h-[250px]">
          <div className="w-20 flex flex-col items-center pt-2">
            <WinRarIcons.Wizard size={56} />
          </div>

          <div className="flex-1 flex flex-col gap-3">
            <div className="text-sm font-bold text-gray-800">
              {step === 1 ? 'Welcome to VibeRAR Wizard' : 'Choose Operation Destination'}
            </div>

            {step === 1 ? (
              <>
                <p className="text-gray-600 text-xs">
                  This wizard simplifies unpacking files from existing archives or creating new archive files.
                </p>
                <div className="flex flex-col gap-2 pt-2">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="wizardAction"
                      checked={action === 'unpack'}
                      onChange={() => setAction('unpack')}
                      className="text-blue-600"
                    />
                    <span className="font-semibold">Unpack an archive</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="wizardAction"
                      checked={action === 'pack'}
                      onChange={() => setAction('pack')}
                      className="text-blue-600"
                    />
                    <span className="font-semibold">Create a new archive</span>
                  </label>
                </div>
              </>
            ) : (
              <p className="text-gray-600 text-xs">
                Click Finish to open the dedicated {action === 'unpack' ? 'Extraction' : 'Compression'} dialog with advanced parameters.
              </p>
            )}
          </div>
        </div>

        {/* Action Buttons */}
        <div className="h-11 bg-[#f0f0f0] border-t border-[#dcdcdc] flex items-center justify-end px-3 gap-2">
          {step === 2 && (
            <button
              onClick={() => setStep(1)}
              className="winrar-btn w-20 py-1 text-xs text-gray-800"
            >
              &lt; Back
            </button>
          )}
          <button
            onClick={handleNext}
            className="winrar-btn w-20 py-1 text-xs font-semibold text-gray-800"
          >
            {step === 1 ? 'Next >' : 'Finish'}
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
