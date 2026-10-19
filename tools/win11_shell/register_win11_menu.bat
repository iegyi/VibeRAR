@echo off
setlocal
echo ========================================================
echo   VibeRAR - Windows 11 Modern Context Menu Registration
echo ========================================================
echo.

cd /d "%~dp0"

echo [*] Importing certificate to TrustedPeople...
powershell -NoProfile -ExecutionPolicy Bypass -Command "Import-Certificate -FilePath 'VibeRAR.cer' -CertStoreLocation Cert:\CurrentUser\TrustedPeople -ErrorAction SilentlyContinue; Import-Certificate -FilePath 'VibeRAR.cer' -CertStoreLocation Cert:\LocalMachine\TrustedPeople -ErrorAction SilentlyContinue"

echo [*] Registering Windows 11 Sparse Package...
powershell -NoProfile -ExecutionPolicy Bypass -Command "Add-AppxPackage -Path '%~dp0VibeRAR.identity.msix' -ExternalLocation '%~dp0' -ErrorAction SilentlyContinue"

echo [*] Refreshing Windows File Explorer...
powershell -NoProfile -ExecutionPolicy Bypass -Command "Stop-Process -Name explorer -Force -ErrorAction SilentlyContinue; Start-Sleep -Milliseconds 500; Start-Process explorer"

echo.
echo [SUCCESS] VibeRAR is now integrated into Windows 11 top-level context menu!
pause
