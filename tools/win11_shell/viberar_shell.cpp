#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <shobjidl.h>
#include <shlwapi.h>
#include <strsafe.h>
#include <string>
#include <vector>

#pragma comment(lib, "shlwapi.lib")
#pragma comment(lib, "ole32.lib")
#pragma comment(lib, "shell32.lib")

// {B09A765D-05A3-4C4F-835E-5C87336EE000}
static const GUID CLSID_VibeRARCommand = 
{ 0xb09a765d, 0x05a3, 0x4c4f, { 0x83, 0x5e, 0x5c, 0x87, 0x33, 0x6e, 0xe0, 0x00 } };

HINSTANCE g_hInstance = NULL;
LONG g_cRefModule = 0;

enum VibeActionType {
    ACTION_ROOT = 0,
    ACTION_OPEN,
    ACTION_EXTRACT_DIALOG,
    ACTION_EXTRACT_HERE,
    ACTION_EXTRACT_TO,
    ACTION_ADD_DIALOG,
    ACTION_ADD_ZIP,
    ACTION_ADD_RAR
};

static void GetVibeRARExePath(LPWSTR outPath, DWORD maxLen) {
    GetModuleFileNameW(g_hInstance, outPath, maxLen);
    PathRemoveFileSpecW(outPath);
    PathCombineW(outPath, outPath, L"VibeRAR.exe");
    if (!PathFileExistsW(outPath)) {
        // Fallback to Program Files default location
        ExpandEnvironmentStringsW(L"%ProgramFiles%\\VibeRAR\\VibeRAR.exe", outPath, maxLen);
    }
}

static bool IsArchiveExtension(LPCWSTR ext) {
    if (!ext || !*ext) return false;
    LPCWSTR exts[] = {
        L".rar", L".zip", L".7z", L".tar", L".gz", L".tgz",
        L".bz2", L".tbz2", L".xz", L".txz", L".iso", L".cab",
        L".arj", L".lzh", L".zst", L".wim"
    };
    for (int i = 0; i < (int)(sizeof(exts) / sizeof(exts[0])); i++) {
        if (_wcsicmp(ext, exts[i]) == 0) return true;
    }
    return false;
}

static void InspectShellItems(IShellItemArray* psiItemArray, bool& isArchive, std::wstring& baseName, std::wstring& firstPath) {
    isArchive = false;
    baseName = L"archive";
    firstPath = L"";

    if (!psiItemArray) return;

    DWORD count = 0;
    if (FAILED(psiItemArray->GetCount(&count)) || count == 0) return;

    IShellItem* pItem = nullptr;
    if (SUCCEEDED(psiItemArray->GetItemAt(0, &pItem)) && pItem) {
        LPWSTR pszPath = nullptr;
        if (SUCCEEDED(pItem->GetDisplayName(SIGDN_FILESYSPATH, &pszPath)) && pszPath) {
            firstPath = pszPath;
            LPCWSTR ext = PathFindExtensionW(pszPath);
            if (IsArchiveExtension(ext)) {
                isArchive = true;
            }

            WCHAR nameBuf[MAX_PATH];
            StringCchCopyW(nameBuf, MAX_PATH, PathFindFileNameW(pszPath));
            PathRemoveExtensionW(nameBuf);
            if (nameBuf[0]) {
                baseName = nameBuf;
            }
            CoTaskMemFree(pszPath);
        }
        pItem->Release();
    }
}

// =========================================================================
// Subcommand Implementation
// =========================================================================
class CVibeSubCommand : public IExplorerCommand {
private:
    LONG m_cRef;
    VibeActionType m_action;

public:
    CVibeSubCommand(VibeActionType action) : m_cRef(1), m_action(action) {
        InterlockedIncrement(&g_cRefModule);
    }

    virtual ~CVibeSubCommand() {
        InterlockedDecrement(&g_cRefModule);
    }

    // IUnknown
    IFACEMETHODIMP QueryInterface(REFIID riid, void** ppv) {
        static const QITAB qit[] = {
            QITABENT(CVibeSubCommand, IExplorerCommand),
            { 0 },
        };
        return QISearch(this, qit, riid, ppv);
    }

    IFACEMETHODIMP_(ULONG) AddRef() {
        return InterlockedIncrement(&m_cRef);
    }

    IFACEMETHODIMP_(ULONG) Release() {
        ULONG cRef = InterlockedDecrement(&m_cRef);
        if (cRef == 0) delete this;
        return cRef;
    }

    // IExplorerCommand
    IFACEMETHODIMP GetTitle(IShellItemArray* psiItemArray, LPWSTR* ppszName) {
        *ppszName = nullptr;
        bool isArchive = false;
        std::wstring baseName = L"archive";
        std::wstring firstPath;
        InspectShellItems(psiItemArray, isArchive, baseName, firstPath);

        std::wstring title;
        switch (m_action) {
            case ACTION_OPEN:
                title = L"Open with VibeRAR";
                break;
            case ACTION_EXTRACT_DIALOG:
                title = L"Extract files...";
                break;
            case ACTION_EXTRACT_HERE:
                title = L"Extract Here";
                break;
            case ACTION_EXTRACT_TO:
                title = L"Extract to \"" + baseName + L"\\\"";
                break;
            case ACTION_ADD_DIALOG:
                title = L"Add to archive...";
                break;
            case ACTION_ADD_ZIP:
                title = L"Add to \"" + baseName + L".zip\"";
                break;
            case ACTION_ADD_RAR:
                title = L"Add to \"" + baseName + L".rar\"";
                break;
            default:
                title = L"VibeRAR";
                break;
        }

        return SHStrDupW(title.c_str(), ppszName);
    }

    IFACEMETHODIMP GetIcon(IShellItemArray* /*psiItemArray*/, LPWSTR* ppszIcon) {
        *ppszIcon = nullptr;
        WCHAR exePath[MAX_PATH];
        GetVibeRARExePath(exePath, MAX_PATH);
        WCHAR iconSpec[MAX_PATH + 8];
        StringCchPrintfW(iconSpec, ARRAYSIZE(iconSpec), L"%s,0", exePath);
        return SHStrDupW(iconSpec, ppszIcon);
    }

    IFACEMETHODIMP GetToolTip(IShellItemArray* /*psiItemArray*/, LPWSTR* ppszInfoTip) {
        *ppszInfoTip = nullptr;
        return E_NOTIMPL;
    }

    IFACEMETHODIMP GetCanonicalName(GUID* pguidCommandName) {
        *pguidCommandName = GUID_NULL;
        return S_OK;
    }

    IFACEMETHODIMP GetState(IShellItemArray* psiItemArray, BOOL /*fOkToBeSlow*/, EXPCMDSTATE* pCmdState) {
        *pCmdState = ECS_ENABLED;
        bool isArchive = false;
        std::wstring baseName, firstPath;
        InspectShellItems(psiItemArray, isArchive, baseName, firstPath);

        if (m_action == ACTION_OPEN || m_action == ACTION_EXTRACT_DIALOG || 
            m_action == ACTION_EXTRACT_HERE || m_action == ACTION_EXTRACT_TO) {
            *pCmdState = isArchive ? ECS_ENABLED : ECS_HIDDEN;
        } else if (m_action == ACTION_ADD_DIALOG || m_action == ACTION_ADD_ZIP || m_action == ACTION_ADD_RAR) {
            *pCmdState = ECS_ENABLED;
        }

        return S_OK;
    }

    IFACEMETHODIMP Invoke(IShellItemArray* psiItemArray, IBindCtx* /*pbc*/) {
        WCHAR exePath[MAX_PATH];
        GetVibeRARExePath(exePath, MAX_PATH);

        std::wstring cmdLine = L"\"";
        cmdLine += exePath;
        cmdLine += L"\" ";

        switch (m_action) {
            case ACTION_OPEN:
                break;
            case ACTION_EXTRACT_DIALOG:
                cmdLine += L"--extract ";
                break;
            case ACTION_EXTRACT_HERE:
                cmdLine += L"--extract-here ";
                break;
            case ACTION_EXTRACT_TO:
                cmdLine += L"--extract-to ";
                break;
            case ACTION_ADD_DIALOG:
                cmdLine += L"--add ";
                break;
            case ACTION_ADD_ZIP:
                cmdLine += L"--quick-zip ";
                break;
            case ACTION_ADD_RAR:
                cmdLine += L"--quick-rar ";
                break;
            default:
                break;
        }

        if (psiItemArray) {
            DWORD count = 0;
            if (SUCCEEDED(psiItemArray->GetCount(&count))) {
                for (DWORD i = 0; i < count; i++) {
                    IShellItem* pItem = nullptr;
                    if (SUCCEEDED(psiItemArray->GetItemAt(i, &pItem)) && pItem) {
                        LPWSTR pszPath = nullptr;
                        if (SUCCEEDED(pItem->GetDisplayName(SIGDN_FILESYSPATH, &pszPath)) && pszPath) {
                            cmdLine += L"\"";
                            cmdLine += pszPath;
                            cmdLine += L"\" ";
                            CoTaskMemFree(pszPath);
                        }
                        pItem->Release();
                    }
                }
            }
        }

        WCHAR dirPath[MAX_PATH];
        StringCchCopyW(dirPath, MAX_PATH, exePath);
        PathRemoveFileSpecW(dirPath);

        STARTUPINFOW si = { sizeof(si) };
        PROCESS_INFORMATION pi = {};
        std::vector<WCHAR> cmdVec(cmdLine.begin(), cmdLine.end());
        cmdVec.push_back(L'\0');

        if (CreateProcessW(NULL, cmdVec.data(), NULL, NULL, FALSE, 0, NULL, dirPath, &si, &pi)) {
            CloseHandle(pi.hThread);
            CloseHandle(pi.hProcess);
        }

        return S_OK;
    }

    IFACEMETHODIMP GetFlags(EXPCMDFLAGS* pFlags) {
        *pFlags = ECF_DEFAULT;
        return S_OK;
    }

    IFACEMETHODIMP EnumSubCommands(IEnumExplorerCommand** ppEnum) {
        *ppEnum = nullptr;
        return E_NOTIMPL;
    }
};

// =========================================================================
// Enumerator Implementation
// =========================================================================
class CVibeEnumExplorerCommand : public IEnumExplorerCommand {
private:
    LONG m_cRef;
    ULONG m_index;
    std::vector<IExplorerCommand*> m_commands;

public:
    CVibeEnumExplorerCommand() : m_cRef(1), m_index(0) {
        InterlockedIncrement(&g_cRefModule);

        m_commands.push_back(new CVibeSubCommand(ACTION_OPEN));
        m_commands.push_back(new CVibeSubCommand(ACTION_EXTRACT_DIALOG));
        m_commands.push_back(new CVibeSubCommand(ACTION_EXTRACT_HERE));
        m_commands.push_back(new CVibeSubCommand(ACTION_EXTRACT_TO));
        m_commands.push_back(new CVibeSubCommand(ACTION_ADD_DIALOG));
        m_commands.push_back(new CVibeSubCommand(ACTION_ADD_ZIP));
        m_commands.push_back(new CVibeSubCommand(ACTION_ADD_RAR));
    }

    virtual ~CVibeEnumExplorerCommand() {
        for (auto pCmd : m_commands) {
            pCmd->Release();
        }
        m_commands.clear();
        InterlockedDecrement(&g_cRefModule);
    }

    // IUnknown
    IFACEMETHODIMP QueryInterface(REFIID riid, void** ppv) {
        static const QITAB qit[] = {
            QITABENT(CVibeEnumExplorerCommand, IEnumExplorerCommand),
            { 0 },
        };
        return QISearch(this, qit, riid, ppv);
    }

    IFACEMETHODIMP_(ULONG) AddRef() {
        return InterlockedIncrement(&m_cRef);
    }

    IFACEMETHODIMP_(ULONG) Release() {
        ULONG cRef = InterlockedDecrement(&m_cRef);
        if (cRef == 0) delete this;
        return cRef;
    }

    // IEnumExplorerCommand
    IFACEMETHODIMP Next(ULONG celt, IExplorerCommand** pUICommand, ULONG* pceltFetched) {
        if (pceltFetched) *pceltFetched = 0;
        if (!pUICommand || (celt > 1 && !pceltFetched)) return E_INVALIDARG;

        ULONG fetched = 0;
        while (fetched < celt && m_index < m_commands.size()) {
            pUICommand[fetched] = m_commands[m_index];
            pUICommand[fetched]->AddRef();
            fetched++;
            m_index++;
        }

        if (pceltFetched) *pceltFetched = fetched;
        return (fetched == celt) ? S_OK : S_FALSE;
    }

    IFACEMETHODIMP Skip(ULONG celt) {
        m_index += celt;
        return (m_index <= m_commands.size()) ? S_OK : S_FALSE;
    }

    IFACEMETHODIMP Reset() {
        m_index = 0;
        return S_OK;
    }

    IFACEMETHODIMP Clone(IEnumExplorerCommand** ppEnum) {
        *ppEnum = nullptr;
        return E_NOTIMPL;
    }
};

// =========================================================================
// Parent Root Command Implementation (VibeRAR >)
// =========================================================================
class CVibeParentCommand : public IExplorerCommand {
private:
    LONG m_cRef;

public:
    CVibeParentCommand() : m_cRef(1) {
        InterlockedIncrement(&g_cRefModule);
    }

    virtual ~CVibeParentCommand() {
        InterlockedDecrement(&g_cRefModule);
    }

    // IUnknown
    IFACEMETHODIMP QueryInterface(REFIID riid, void** ppv) {
        static const QITAB qit[] = {
            QITABENT(CVibeParentCommand, IExplorerCommand),
            { 0 },
        };
        return QISearch(this, qit, riid, ppv);
    }

    IFACEMETHODIMP_(ULONG) AddRef() {
        return InterlockedIncrement(&m_cRef);
    }

    IFACEMETHODIMP_(ULONG) Release() {
        ULONG cRef = InterlockedDecrement(&m_cRef);
        if (cRef == 0) delete this;
        return cRef;
    }

    // IExplorerCommand
    IFACEMETHODIMP GetTitle(IShellItemArray* /*psiItemArray*/, LPWSTR* ppszName) {
        *ppszName = nullptr;
        return SHStrDupW(L"VibeRAR", ppszName);
    }

    IFACEMETHODIMP GetIcon(IShellItemArray* /*psiItemArray*/, LPWSTR* ppszIcon) {
        *ppszIcon = nullptr;
        WCHAR exePath[MAX_PATH];
        GetVibeRARExePath(exePath, MAX_PATH);
        WCHAR iconSpec[MAX_PATH + 8];
        StringCchPrintfW(iconSpec, ARRAYSIZE(iconSpec), L"%s,0", exePath);
        return SHStrDupW(iconSpec, ppszIcon);
    }

    IFACEMETHODIMP GetToolTip(IShellItemArray* /*psiItemArray*/, LPWSTR* ppszInfoTip) {
        *ppszInfoTip = nullptr;
        return SHStrDupW(L"VibeRAR Archive Manager", ppszInfoTip);
    }

    IFACEMETHODIMP GetCanonicalName(GUID* pguidCommandName) {
        *pguidCommandName = CLSID_VibeRARCommand;
        return S_OK;
    }

    IFACEMETHODIMP GetState(IShellItemArray* /*psiItemArray*/, BOOL /*fOkToBeSlow*/, EXPCMDSTATE* pCmdState) {
        *pCmdState = ECS_ENABLED;
        return S_OK;
    }

    IFACEMETHODIMP Invoke(IShellItemArray* psiItemArray, IBindCtx* /*pbc*/) {
        // Direct click fallback: open with selection
        CVibeSubCommand openCmd(ACTION_OPEN);
        return openCmd.Invoke(psiItemArray, nullptr);
    }

    IFACEMETHODIMP GetFlags(EXPCMDFLAGS* pFlags) {
        *pFlags = ECF_HASSUBCOMMANDS;
        return S_OK;
    }

    IFACEMETHODIMP EnumSubCommands(IEnumExplorerCommand** ppEnum) {
        *ppEnum = new (std::nothrow) CVibeEnumExplorerCommand();
        return (*ppEnum != nullptr) ? S_OK : E_OUTOFMEMORY;
    }
};

// =========================================================================
// Class Factory Implementation
// =========================================================================
class CVibeClassFactory : public IClassFactory {
private:
    LONG m_cRef;

public:
    CVibeClassFactory() : m_cRef(1) {
        InterlockedIncrement(&g_cRefModule);
    }

    virtual ~CVibeClassFactory() {
        InterlockedDecrement(&g_cRefModule);
    }

    IFACEMETHODIMP QueryInterface(REFIID riid, void** ppv) {
        static const QITAB qit[] = {
            QITABENT(CVibeClassFactory, IClassFactory),
            { 0 },
        };
        return QISearch(this, qit, riid, ppv);
    }

    IFACEMETHODIMP_(ULONG) AddRef() {
        return InterlockedIncrement(&m_cRef);
    }

    IFACEMETHODIMP_(ULONG) Release() {
        ULONG cRef = InterlockedDecrement(&m_cRef);
        if (cRef == 0) delete this;
        return cRef;
    }

    IFACEMETHODIMP CreateInstance(IUnknown* pUnkOuter, REFIID riid, void** ppv) {
        *ppv = nullptr;
        if (pUnkOuter) return CLASS_E_NOAGGREGATION;

        CVibeParentCommand* pCmd = new (std::nothrow) CVibeParentCommand();
        if (!pCmd) return E_OUTOFMEMORY;

        HRESULT hr = pCmd->QueryInterface(riid, ppv);
        pCmd->Release();
        return hr;
    }

    IFACEMETHODIMP LockServer(BOOL fLock) {
        if (fLock) {
            InterlockedIncrement(&g_cRefModule);
        } else {
            InterlockedDecrement(&g_cRefModule);
        }
        return S_OK;
    }
};

// =========================================================================
// DLL Entrypoints
// =========================================================================
BOOL APIENTRY DllMain(HMODULE hModule, DWORD ul_reason_for_call, LPVOID /*lpReserved*/) {
    if (ul_reason_for_call == DLL_PROCESS_ATTACH) {
        g_hInstance = (HINSTANCE)hModule;
        DisableThreadLibraryCalls(hModule);
    }
    return TRUE;
}

STDAPI DllGetClassObject(REFCLSID rclsid, REFIID riid, void** ppv) {
    *ppv = nullptr;
    if (IsEqualCLSID(rclsid, CLSID_VibeRARCommand)) {
        CVibeClassFactory* pFactory = new (std::nothrow) CVibeClassFactory();
        if (!pFactory) return E_OUTOFMEMORY;
        HRESULT hr = pFactory->QueryInterface(riid, ppv);
        pFactory->Release();
        return hr;
    }
    return CLASS_E_CLASSNOTAVAILABLE;
}

STDAPI DllCanUnloadNow(void) {
    return (g_cRefModule == 0) ? S_OK : S_FALSE;
}
