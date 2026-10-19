import React, { useState } from 'react';
import { useArchiveStore } from '../store/useArchiveStore';
import { WinRarIcons } from '../utils/icons';
import { ChevronDown, ChevronRight, Folder, FolderOpen } from 'lucide-react';

interface TreeItemProps {
  name: string;
  fullPath: string;
  isFolder: boolean;
  children?: TreeItemProps[];
}

export const FolderTree: React.FC = () => {
  const { currentPath, navigateTo, isTreeCollapsed } = useArchiveStore();
  const [expandedPaths, setExpandedPaths] = useState<Set<string>>(
    new Set(['project_backup.rar', 'project_backup.rar/src'])
  );

  if (isTreeCollapsed) return null;

  const toggleExpand = (path: string, e: React.MouseEvent) => {
    e.stopPropagation();
    const next = new Set(expandedPaths);
    if (next.has(path)) {
      next.delete(path);
    } else {
      next.add(path);
    }
    setExpandedPaths(next);
  };

  const treeData: TreeItemProps = {
    name: 'project_backup.rar',
    fullPath: 'project_backup.rar',
    isFolder: true,
    children: [
      { name: 'docs', fullPath: 'project_backup.rar/docs', isFolder: true },
      { name: 'assets', fullPath: 'project_backup.rar/assets', isFolder: true },
      {
        name: 'src',
        fullPath: 'project_backup.rar/src',
        isFolder: true,
        children: [
          { name: 'components', fullPath: 'project_backup.rar/src/components', isFolder: true },
          { name: 'utils', fullPath: 'project_backup.rar/src/utils', isFolder: true },
          { name: 'tests', fullPath: 'project_backup.rar/src/tests', isFolder: true }
        ]
      }
    ]
  };

  const renderNode = (node: TreeItemProps, level: number = 0) => {
    const isExpanded = expandedPaths.has(node.fullPath);
    const isSelected = currentPath === node.fullPath;
    const hasChildren = node.children && node.children.length > 0;

    return (
      <div key={node.fullPath} className="select-none">
        <div
          onClick={() => navigateTo(node.fullPath)}
          className={`flex items-center gap-1.5 py-1 px-1 rounded-xs cursor-pointer text-xs transition-colors ${
            isSelected
              ? 'bg-[#3399ff] text-white font-medium'
              : 'hover:bg-[#e5f1fb] text-[#222]'
          }`}
          style={{ paddingLeft: `${level * 14 + 4}px` }}
        >
          {hasChildren ? (
            <span
              onClick={(e) => toggleExpand(node.fullPath, e)}
              className="p-0.5 hover:bg-black/10 rounded-xs cursor-pointer"
            >
              {isExpanded ? (
                <ChevronDown size={12} className={isSelected ? 'text-white' : 'text-gray-500'} />
              ) : (
                <ChevronRight size={12} className={isSelected ? 'text-white' : 'text-gray-500'} />
              )}
            </span>
          ) : (
            <span className="w-3"></span>
          )}

          {level === 0 ? (
            <WinRarIcons.ArchiveFile size={16} />
          ) : isExpanded ? (
            <FolderOpen size={15} className={isSelected ? 'text-white' : 'text-[#ffb300] fill-[#ffe082]'} />
          ) : (
            <Folder size={15} className={isSelected ? 'text-white' : 'text-[#f57f17] fill-[#ffd54f]'} />
          )}

          <span className="truncate text-[12px]">{node.name}</span>
        </div>

        {hasChildren && isExpanded && (
          <div>
            {node.children!.map((child) => renderNode(child, level + 1))}
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="w-52 min-w-44 max-w-72 bg-white border-r border-[#c8c8c8] flex flex-col h-full overflow-y-auto p-1.5 shadow-inner">
      <div className="text-[10px] uppercase font-bold text-gray-500 px-2 py-1 tracking-wider border-b border-gray-200 mb-1">
        Folders Tree
      </div>
      {renderNode(treeData)}
    </div>
  );
};
