@echo off
setlocal
echo ========================================================
echo   VibeRAR - Windows 11 Modern Context Menu Removal
echo ========================================================
echo.

echo [*] Removing Windows 11 Package...
powershell -NoProfile -ExecutionPolicy Bypass -Command "Get-AppxPackage *VibeRAR* | Remove-AppxPackage -ErrorAction SilentlyContinue"

echo [*] Refreshing Windows File Explorer...
powershell -NoProfile -ExecutionPolicy Bypass -Command "Stop-Process -Name explorer -Force -ErrorAction SilentlyContinue; Start-Sleep -Milliseconds 500; Start-Process explorer"

echo.
echo [SUCCESS] VibeRAR context menu removed.
pause
