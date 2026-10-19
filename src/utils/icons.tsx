import React from 'react';

// Custom authentic WinRAR vector icons matching Figure 1 and classical WinRAR themes
export const WinRarIcons = {
  // Classic 3 books tied with a belt (WinRAR's signature brand)
  AddArchive: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      {/* Blue Book */}
      <rect x="8" y="10" width="32" height="7" rx="1.5" fill="#1b6ca8" stroke="#104a75" strokeWidth="1.5" />
      <rect x="36" y="10" width="3" height="7" fill="#f4eedb" />
      <line x1="12" y1="13.5" x2="33" y2="13.5" stroke="#4ba3e3" strokeWidth="1.5" strokeLinecap="round" />
      
      {/* Green Book */}
      <rect x="7" y="18" width="34" height="7.5" rx="1.5" fill="#2e7d32" stroke="#1b5e20" strokeWidth="1.5" />
      <rect x="37" y="18" width="3" height="7.5" fill="#f4eedb" />
      <line x1="11" y1="21.75" x2="34" y2="21.75" stroke="#66bb6a" strokeWidth="1.5" strokeLinecap="round" />

      {/* Red Book */}
      <rect x="6" y="26.5" width="36" height="8.5" rx="1.5" fill="#c62828" stroke="#8e0000" strokeWidth="1.5" />
      <rect x="38" y="26.5" width="3" height="8.5" fill="#f4eedb" />
      <line x1="10" y1="30.75" x2="35" y2="30.75" stroke="#ef5350" strokeWidth="1.5" strokeLinecap="round" />

      {/* Leather Belt & Buckle */}
      <rect x="22" y="8" width="5.5" height="29" rx="1" fill="#4e342e" stroke="#2d1d1a" strokeWidth="1" />
      <rect x="20.5" y="20.5" width="8.5" height="7" rx="1.5" fill="#ffd54f" stroke="#ff8f00" strokeWidth="1.2" />
      <rect x="22.5" y="22.5" width="4.5" height="3" fill="#3e2723" />
      <line x1="25" y1="20.5" x2="25" y2="27.5" stroke="#ffe082" strokeWidth="1" />

      {/* Plus badge */}
      <circle cx="38" cy="38" r="8" fill="#2e7d32" stroke="#ffffff" strokeWidth="2" />
      <line x1="38" y1="34" x2="38" y2="42" stroke="#ffffff" strokeWidth="2.5" strokeLinecap="round" />
      <line x1="34" y1="38" x2="42" y2="38" stroke="#ffffff" strokeWidth="2.5" strokeLinecap="round" />
    </svg>
  ),

  // Extract To (Cabinet / Folder with arrow pointing out)
  ExtractTo: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      {/* Folder Back */}
      <path d="M6 12C6 10.8954 6.89543 10 8 10H18L22 14H38C39.1046 14 40 14.8954 40 16V36C40 37.1046 39.1046 38 38 38H8C6.89543 38 6 37.1046 6 36V12Z" fill="#e0a93b" stroke="#b27f18" strokeWidth="1.5" />
      {/* Folder Front Open */}
      <path d="M4 22L8.5 38H41.5L44 22H4Z" fill="#ffd15c" stroke="#cf9518" strokeWidth="1.5" />
      {/* Dynamic Blue Arrow */}
      <path d="M24 6L33 15H27V25H21V15H15L24 6Z" fill="#1976d2" stroke="#0d47a1" strokeWidth="1.5" />
      <path d="M24 9L29 14H25V23H23V14H19L24 9Z" fill="#64b5f6" />
    </svg>
  ),

  // Test Archive (Archive integrity test checkmark)
  Test: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      {/* Document / Clipboard */}
      <rect x="10" y="8" width="28" height="34" rx="2" fill="#ffffff" stroke="#546e7a" strokeWidth="2" />
      <rect x="16" y="5" width="16" height="6" rx="2" fill="#78909c" stroke="#455a64" strokeWidth="1.5" />
      {/* Lines */}
      <line x1="16" y1="18" x2="32" y2="18" stroke="#b0bec5" strokeWidth="2" strokeLinecap="round" />
      <line x1="16" y1="24" x2="28" y2="24" stroke="#b0bec5" strokeWidth="2" strokeLinecap="round" />
      <line x1="16" y1="30" x2="24" y2="30" stroke="#b0bec5" strokeWidth="2" strokeLinecap="round" />
      {/* Green Checkmark */}
      <circle cx="34" cy="34" r="10" fill="#2e7d32" stroke="#ffffff" strokeWidth="2" />
      <path d="M29 34L33 38L39 30" stroke="#ffffff" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  ),

  // View File
  View: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="8" y="8" width="24" height="32" rx="2" fill="#ffffff" stroke="#78909c" strokeWidth="2" />
      <line x1="12" y1="14" x2="24" y2="14" stroke="#90a4ae" strokeWidth="2" strokeLinecap="round" />
      <line x1="12" y1="20" x2="22" y2="20" stroke="#90a4ae" strokeWidth="2" strokeLinecap="round" />
      <line x1="12" y1="26" x2="20" y2="26" stroke="#90a4ae" strokeWidth="2" strokeLinecap="round" />
      {/* Magnifier */}
      <circle cx="28" cy="26" r="9" fill="rgba(144, 202, 249, 0.4)" stroke="#1565c0" strokeWidth="2.5" />
      <line x1="35" y1="33" x2="43" y2="41" stroke="#37474f" strokeWidth="4.5" strokeLinecap="round" />
      <line x1="35" y1="33" x2="43" y2="41" stroke="#90a4ae" strokeWidth="2" strokeLinecap="round" />
    </svg>
  ),

  // Delete
  Delete: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="9" y="9" width="26" height="32" rx="2" fill="#ffffff" stroke="#78909c" strokeWidth="2" />
      <line x1="14" y1="16" x2="26" y2="16" stroke="#b0bec5" strokeWidth="2" />
      <line x1="14" y1="22" x2="24" y2="22" stroke="#b0bec5" strokeWidth="2" />
      {/* Red Cross */}
      <circle cx="34" cy="33" r="10" fill="#d32f2f" stroke="#ffffff" strokeWidth="2" />
      <line x1="29" y1="28" x2="39" y2="38" stroke="#ffffff" strokeWidth="3" strokeLinecap="round" />
      <line x1="39" y1="28" x2="29" y2="38" stroke="#ffffff" strokeWidth="3" strokeLinecap="round" />
    </svg>
  ),

  // Find / Search
  Find: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      {/* Classic Binoculars / Lens */}
      <ellipse cx="17" cy="21" rx="8" ry="8" fill="#e3f2fd" stroke="#0d47a1" strokeWidth="2.5" />
      <ellipse cx="31" cy="21" rx="8" ry="8" fill="#e3f2fd" stroke="#0d47a1" strokeWidth="2.5" />
      <rect x="22" y="19" width="4" height="4" fill="#37474f" />
      <path d="M11 25L8 38H16L18 28" fill="#1565c0" stroke="#0d47a1" strokeWidth="1.5" />
      <path d="M37 25L40 38H32L30 28" fill="#1565c0" stroke="#0d47a1" strokeWidth="1.5" />
    </svg>
  ),

  // Wizard
  Wizard: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M8 40L30 18L34 22L12 44L8 40Z" fill="#3e2723" stroke="#212121" strokeWidth="1.5" />
      <path d="M26 22L30 18L34 22L30 26L26 22Z" fill="#ffd54f" />
      <polygon points="37,8 39,13 44,14 40,18 41,23 37,20 33,23 34,18 30,14 35,13" fill="#ffca28" stroke="#ff8f00" strokeWidth="1" />
      <polygon points="18,10 19,12 22,12 20,15 21,17 18,16 16,17 17,15 15,12 17,12" fill="#4fc3f7" />
      <polygon points="12,24 13,26 15,26 14,28 15,30 12,29 10,30 11,28 9,26 11,26" fill="#81c784" />
    </svg>
  ),

  // Information
  Info: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="24" cy="24" r="18" fill="#0288d1" stroke="#01579b" strokeWidth="2.5" />
      <circle cx="24" cy="15" r="3" fill="#ffffff" />
      <rect x="21" y="21" width="6" height="15" rx="2" fill="#ffffff" />
    </svg>
  ),

  // Repair
  Repair: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      {/* First Aid Box */}
      <rect x="8" y="14" width="32" height="24" rx="3" fill="#ffffff" stroke="#c62828" strokeWidth="2.5" />
      <path d="M19 14V10C19 8.89543 19.8954 8 21 8H27C28.1046 8 29 8.89543 29 10V14" stroke="#455a64" strokeWidth="2.5" />
      {/* Red Cross */}
      <rect x="21" y="19" width="6" height="14" fill="#d32f2f" />
      <rect x="17" y="23" width="14" height="6" fill="#d32f2f" />
    </svg>
  ),

  // Benchmark
  Benchmark: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      {/* Gauge */}
      <circle cx="24" cy="24" r="17" fill="#263238" stroke="#37474f" strokeWidth="2" />
      <circle cx="24" cy="24" r="14" fill="#37474f" />
      {/* Tick Marks */}
      <circle cx="24" cy="24" r="11" stroke="#eceff1" strokeWidth="1.5" strokeDasharray="3 5" />
      {/* Gauge dial needle */}
      <circle cx="24" cy="24" r="4" fill="#ff5252" />
      <line x1="24" y1="24" x2="33" y2="15" stroke="#ff5252" strokeWidth="2.5" strokeLinecap="round" />
      {/* Speed sparks */}
      <polygon points="12,14 17,20 13,20 15,26 9,19 13,19" fill="#ffeb3b" />
    </svg>
  ),

  // SFX Generator
  SFX: ({ size = 32 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="8" y="10" width="32" height="30" rx="3" fill="#37474f" stroke="#263238" strokeWidth="2" />
      <path d="M12 16H36" stroke="#546e7a" strokeWidth="2" />
      <text x="14" y="28" fill="#ffca28" fontSize="10" fontWeight="bold" fontFamily="sans-serif">SFX</text>
      <circle cx="33" cy="31" r="5" fill="#00c853" stroke="#ffffff" strokeWidth="1.5" />
      <text x="30" y="34" fill="#ffffff" fontSize="8" fontWeight="bold" fontFamily="sans-serif">.EXE</text>
    </svg>
  ),

  // Up One Level Folder Icon
  UpFolder: ({ size = 20 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M3 5C3 4.44772 3.44772 4 4 4H9L11 6H19C19.5523 6 20 6.44772 20 7V17C20 17.5523 19.5523 18 19 18H4C3.44772 18 3 17.5523 3 17V5Z" fill="#ffc107" stroke="#ffa000" strokeWidth="1" />
      <path d="M2 10H20L19 18H3L2 10Z" fill="#ffe082" stroke="#ffb300" strokeWidth="1" />
      <path d="M11 15V9M11 9L8 12M11 9L14 12" stroke="#1565c0" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  ),

  // Folder icon
  Folder: ({ size = 18 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M3 5C3 4.44772 3.44772 4 4 4H9L11 6H19C19.5523 6 20 6.44772 20 7V17C20 17.5523 19.5523 18 19 18H4C3.44772 18 3 17.5523 3 17V5Z" fill="#ffb300" stroke="#f57f17" strokeWidth="1" />
      <path d="M2 10H21L19 18H3L2 10Z" fill="#ffe082" stroke="#ffb300" strokeWidth="1" />
    </svg>
  ),

  // Generic File
  File: ({ size = 18 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M5 3C5 2.44772 5.44772 2 6 2H14L19 7V21C19 21.5523 18.5523 22 18 22H6C5.44772 22 5 21.5523 5 21V3Z" fill="#ffffff" stroke="#90a4ae" strokeWidth="1.2" />
      <path d="M14 2V7H19" stroke="#90a4ae" strokeWidth="1.2" fill="#cfd8dc" />
    </svg>
  ),

  // Archive file icon (RAR/ZIP)
  ArchiveFile: ({ size = 18 }: { size?: number }) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="4" y="3" width="16" height="4" rx="1" fill="#1b6ca8" stroke="#104a75" strokeWidth="0.8" />
      <rect x="4" y="8" width="16" height="4.5" rx="1" fill="#2e7d32" stroke="#1b5e20" strokeWidth="0.8" />
      <rect x="4" y="13.5" width="16" height="5.5" rx="1" fill="#c62828" stroke="#8e0000" strokeWidth="0.8" />
      <rect x="11" y="2" width="3" height="18" fill="#5d4037" />
      <rect x="10" y="9" width="5" height="4" rx="0.5" fill="#ffd54f" stroke="#ff8f00" strokeWidth="0.8" />
    </svg>
  )
};
