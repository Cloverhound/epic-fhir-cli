# FHIR CLI installer for Windows
# Usage: irm https://raw.githubusercontent.com/jbogarin/fhir-cli/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "jbogarin/fhir-cli"
$Binary = "fhir-cli.exe"
$InstallDir = "$env:LOCALAPPDATA\fhir-cli"

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "Unsupported: 32-bit Windows is not supported"
    return
}

# Get latest version
Write-Host "Fetching latest release..."
$Release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
$Version = $Release.tag_name -replace '^v', ''
if (-not $Version) {
    Write-Error "Could not determine latest version"
    return
}
Write-Host "Latest version: v$Version"

# Download
$ZipName = "fhir-cli_${Version}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/v$Version/$ZipName"

$TmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "fhir-cli-install-$(Get-Random)")

try {
    Write-Host "Downloading $Url..."
    Invoke-WebRequest -Uri $Url -OutFile (Join-Path $TmpDir $ZipName)

    # Extract
    Expand-Archive -Path (Join-Path $TmpDir $ZipName) -DestinationPath $TmpDir -Force

    # Install
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Move-Item -Path (Join-Path $TmpDir $Binary) -Destination (Join-Path $InstallDir $Binary) -Force

    # Add to user PATH if not already there
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$InstallDir;$UserPath", "User")
        $env:Path = "$InstallDir;$env:Path"
        Write-Host "Added $InstallDir to your PATH."
    }

    Write-Host ""
    Write-Host "Installed fhir-cli v$Version to $InstallDir\$Binary"
    Write-Host ""
    Write-Host "NOTE: Restart your terminal for PATH changes to take effect."
    Write-Host ""
    Write-Host "Get started:"
    Write-Host "  fhir-cli config init    # Interactive setup"
    Write-Host "  fhir-cli auth token     # Get an access token"
    Write-Host "  fhir-cli patient search --family Smith"
} finally {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
}
