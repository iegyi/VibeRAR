@echo off
setlocal
echo ========================================================
echo   VibeRAR - Windows 11 Modern Context Menu Registration
echo ========================================================
echo.

cd /d "%~dp0"

echo [*] Registering COM InProc Server...
reg add "HKCU\Software\Classes\CLSID\{b09a765d-05a3-4c4f-835e-5c87336ee000}" /ve /d "VibeRAR Modern Shell Command" /f >nul 2>&1
reg add "HKCU\Software\Classes\CLSID\{b09a765d-05a3-4c4f-835e-5c87336ee000}\InProcServer32" /ve /d "%~dp0viberar_shell.dll" /f >nul 2>&1
reg add "HKCU\Software\Classes\CLSID\{b09a765d-05a3-4c4f-835e-5c87336ee000}\InProcServer32" /v "ThreadingModel" /d "Apartment" /f >nul 2>&1

echo [*] Importing certificate to TrustedPeople...
certutil -user -addstore -f "TrustedPeople" "%~dp0VibeRAR.cer" >nul 2>&1
powershell -NoProfile -ExecutionPolicy Bypass -Command "Import-Certificate -FilePath '%~dp0VibeRAR.cer' -CertStoreLocation Cert:\CurrentUser\TrustedPeople -ErrorAction SilentlyContinue" >nul 2>&1

echo [*] Registering Windows 11 Sparse Package...
powershell -NoProfile -ExecutionPolicy Bypass -Command "Add-AppxPackage -Path '%~dp0VibeRAR.identity.msix' -ExternalLocation '%~dp0' -ErrorAction SilentlyContinue"

echo [*] Refreshing Windows File Explorer...
powershell -NoProfile -ExecutionPolicy Bypass -Command "Stop-Process -Name explorer -Force -ErrorAction SilentlyContinue; Start-Sleep -Milliseconds 500; Start-Process explorer"

echo.
echo [SUCCESS] VibeRAR is now integrated into Windows 11 top-level context menu!
pause

