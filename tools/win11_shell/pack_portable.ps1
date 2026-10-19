$ErrorActionPreference = "Stop"

$root = "c:\Users\Administrator\Desktop\winrarproject"
$tempDir = Join-Path $root "temp_portable_pack"
$zipOut = Join-Path $root "VibeRAR_Portable_v1.1.0.zip"

if (Test-Path $tempDir) { Remove-Item $tempDir -Recurse -Force }
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

Copy-Item (Join-Path $root "VibeRAR_Portable.exe") $tempDir
Copy-Item (Join-Path $root "tools\win11_shell\viberar_shell.dll") $tempDir
Copy-Item (Join-Path $root "tools\win11_shell\VibeRAR.identity.msix") $tempDir
Copy-Item (Join-Path $root "tools\win11_shell\VibeRAR.cer") $tempDir
Copy-Item (Join-Path $root "tools\win11_shell\package\Assets") $tempDir -Recurse
Copy-Item (Join-Path $root "tools\win11_shell\register_win11_menu.bat") $tempDir
Copy-Item (Join-Path $root "tools\win11_shell\unregister_win11_menu.bat") $tempDir
"VibeRAR Portable v1.1.0 - WinRAR Clone with Windows 11 Modern Context Menu Support`r`nRun register_win11_menu.bat as Administrator to enable top-level right-click menu." | Set-Content (Join-Path $tempDir "README.txt")

if (Test-Path $zipOut) { Remove-Item $zipOut -Force }
Compress-Archive -Path "$tempDir\*" -DestinationPath $zipOut -Force
Remove-Item $tempDir -Recurse -Force

Write-Host "[SUCCESS] Created: $zipOut"
