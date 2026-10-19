$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$pkgDir = Join-Path $scriptDir "package"
$outMsix = Join-Path $scriptDir "VibeRAR.identity.msix"
$pfxFile = Join-Path $scriptDir "VibeRAR.pfx"
$cerFile = Join-Path $scriptDir "VibeRAR.cer"
$pfxPass = "VibeRAR2026"

# Remove any certs inside package dir if present
Remove-Item (Join-Path $pkgDir "VibeRAR.pfx") -ErrorAction SilentlyContinue
Remove-Item (Join-Path $pkgDir "VibeRAR.cer") -ErrorAction SilentlyContinue

Write-Host ">>> Generating Code Signing Certificate..."
$cert = New-SelfSignedCertificate -Type CodeSigningCert -Subject "CN=VibeRAR" -KeyUsage DigitalSignature -FriendlyName "VibeRAR Context Menu Certificate" -CertStoreLocation "Cert:\CurrentUser\My"
$securePass = ConvertTo-SecureString -String $pfxPass -Force -AsPlainText
Export-PfxCertificate -Cert $cert -FilePath $pfxFile -Password $securePass | Out-Null
Export-Certificate -Cert $cert -FilePath $cerFile | Out-Null
Write-Host ">>> Certificate generated: $cerFile"

$makeAppx = "C:\Program Files (x86)\Windows Kits\10\bin\10.0.26100.0\x64\makeappx.exe"
$signTool = "C:\Program Files (x86)\Windows Kits\10\bin\10.0.26100.0\x64\signtool.exe"

Write-Host ">>> Packing Sparse MSIX Package..."
& $makeAppx pack /d "$pkgDir" /nv /p "$outMsix" /o
if ($LASTEXITCODE -ne 0) {
    throw "makeappx failed with exit code $LASTEXITCODE"
}

Write-Host ">>> Signing Sparse MSIX Package..."
& $signTool sign /fd SHA256 /f "$pfxFile" /p "$pfxPass" "$outMsix"
if ($LASTEXITCODE -ne 0) {
    throw "signtool failed with exit code $LASTEXITCODE"
}

Write-Host ">>> [SUCCESS] Sparse package built and signed: $outMsix"
