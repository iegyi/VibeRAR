@echo off
setlocal
call "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat" x64

cl.exe /O2 /W3 /D_USRDLL /D_WINDLL /EHsc /MT tools\win11_shell\viberar_shell.cpp /link /DLL /DEF:tools\win11_shell\viberar_shell.def /OUT:tools\win11_shell\viberar_shell.dll shlwapi.lib ole32.lib shell32.lib advapi32.lib

if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] viberar_shell.dll compiled successfully.
) else (
    echo [ERROR] Compilation failed.
)
