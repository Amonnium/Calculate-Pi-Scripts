# Calculate Pi Scripts CLI Installer - Windows
# Installs the precompiled CLI from the calculate-pi-cli-v1.0 GitHub release.

[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ProjectName = 'Calculate Pi Scripts'
$CliVersion = '1.0'
$ReleaseTag = 'calculate-pi-cli-v1.0'
$AssetName = 'calculate-pi-cli.exe'
$AssetSize = 12623360
$AssetSha256 = 'bb28c9e5c98f62a8311fc18c67ae19750e998696b5b7480649163b544173e211'
$AssetUrl = "https://github.com/Amonnium/Calculate-Pi-Scripts/releases/download/$ReleaseTag/$AssetName"
$InstallDirectory = Join-Path $env:LOCALAPPDATA 'Programs\Calculate-Pi-Scripts\bin'
$InstallPath = Join-Path $InstallDirectory $AssetName
$TempDirectory = Join-Path ([IO.Path]::GetTempPath()) ("calculate-pi-cli-{0}" -f ([Guid]::NewGuid().ToString('N')))
$DownloadPath = Join-Path $TempDirectory $AssetName
$StagePath = Join-Path $InstallDirectory '.calculate-pi-cli.new.exe'
$BackupPath = Join-Path $InstallDirectory '.calculate-pi-cli.previous.exe'

$RollbackNeeded = $false
$BackupCreated = $false

function Write-Stage {
    param(
        [Parameter(Mandatory = $true)][int]$Number,
        [Parameter(Mandatory = $true)][string]$Message
    )

    Write-Host ("[{0}/6] {1}" -f $Number, $Message) -ForegroundColor Cyan
}

function Write-Success {
    param([Parameter(Mandatory = $true)][string]$Message)
    Write-Host $Message -ForegroundColor Green
}

function Fail-Installer {
    param([Parameter(Mandatory = $true)][string]$Message)
    throw $Message
}

function Test-WindowsX64 {
    if (-not [Environment]::Is64BitOperatingSystem) {
        Fail-Installer 'This release supports Windows x64 only. A 32-bit Windows installation was detected.'
    }

    $architecture = $env:PROCESSOR_ARCHITECTURE
    $architectureWow = $env:PROCESSOR_ARCHITEW6432

    if ($architecture -eq 'ARM64' -or $architectureWow -eq 'ARM64') {
        Fail-Installer 'This release provides a Windows x64 executable, but Windows ARM64 was detected. No ARM64 CLI asset is available in this release.'
    }
}

function Test-PeX64 {
    param([Parameter(Mandatory = $true)][string]$Path)

    $bytes = [IO.File]::ReadAllBytes($Path)
    try {
        if ($bytes.Length -lt 64) {
            return $false
        }

        if ($bytes[0] -ne 0x4D -or $bytes[1] -ne 0x5A) {
            return $false
        }

        $peOffset = [BitConverter]::ToInt32($bytes, 0x3C)
        if ($peOffset -lt 0 -or ($peOffset + 6) -gt $bytes.Length) {
            return $false
        }

        if ($bytes[$peOffset] -ne 0x50 -or
            $bytes[$peOffset + 1] -ne 0x45 -or
            $bytes[$peOffset + 2] -ne 0x00 -or
            $bytes[$peOffset + 3] -ne 0x00) {
            return $false
        }

        $machine = [BitConverter]::ToUInt16($bytes, $peOffset + 4)
        return ($machine -eq 0x8664)
    }
    finally {
        $bytes = $null
    }
}

function Get-Sha256 {
    param([Parameter(Mandatory = $true)][string]$Path)
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Add-UserPath {
    param([Parameter(Mandatory = $true)][string]$Directory)

    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if ($null -eq $userPath) {
        $userPath = ''
    }

    $entries = @($userPath -split ';' | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
    $alreadyPresent = $entries | Where-Object {
        [string]::Equals($_.TrimEnd('\'), $Directory.TrimEnd('\'), [StringComparison]::OrdinalIgnoreCase)
    }

    if ($null -ne $alreadyPresent) {
        return $false
    }

    $newPath = if ([string]::IsNullOrWhiteSpace($userPath)) {
        $Directory
    }
    else {
        "$userPath;$Directory"
    }

    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    $env:Path = "$Directory;$env:Path"
    return $true
}

function Remove-TemporaryFiles {
    if (Test-Path -LiteralPath $StagePath -PathType Leaf) {
        Remove-Item -LiteralPath $StagePath -Force -ErrorAction SilentlyContinue
    }

    if (Test-Path -LiteralPath $DownloadPath -PathType Leaf) {
        Remove-Item -LiteralPath $DownloadPath -Force -ErrorAction SilentlyContinue
    }

    if (Test-Path -LiteralPath $TempDirectory -PathType Container) {
        Remove-Item -LiteralPath $TempDirectory -Recurse -Force -ErrorAction SilentlyContinue
    }
}

try {
    Clear-Host

    Write-Host '   ____      _            _       ____  _'
    Write-Host '  / ___|__ _| | ___ _   _| | __ _|  _ \(_)' 
    Write-Host ' | |   / _` | |/ __| | | | |/ _` | |_) | |'
    Write-Host ' | |__| (_| | | (__| |_| | | (_| |  __/| |'
    Write-Host '  \\____\\__,_|_|\\___|\\__,_|_|\\__,_|_|   |_|'
    Write-Host ''
    Write-Host 'Calculate Pi Scripts CLI Installer' -ForegroundColor White
    Write-Host ("Version {0} | Windows x64 | User installation" -f $CliVersion)
    Write-Host ''
    Write-Host ("Project: {0}" -f $ProjectName)
    Write-Host 'License: MIT'
    Write-Host 'Author: Amonnium'
    Write-Host ''

    Write-Stage 1 'Checking system...'
    if ($env:OS -ne 'Windows_NT') {
        Fail-Installer 'This installer supports Windows only.'
    }

    if ([string]::IsNullOrWhiteSpace($env:LOCALAPPDATA)) {
        Fail-Installer 'LOCALAPPDATA is not available. A user-level Windows installation cannot continue safely.'
    }

    Test-WindowsX64

    Write-Stage 2 'Preparing installation...'
    New-Item -ItemType Directory -Path $TempDirectory -Force | Out-Null
    New-Item -ItemType Directory -Path $InstallDirectory -Force | Out-Null

    if (Test-Path -LiteralPath $InstallPath -PathType Container) {
        Fail-Installer "The install target '$InstallPath' is a directory. It was left unchanged."
    }

    Write-Host ("Install location: {0}" -f $InstallPath)
    Write-Host ''

    Write-Stage 3 'Downloading Calculate Pi Scripts...'
    try {
        [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    }
    catch {
        # PowerShell 7 / newer .NET versions may manage TLS automatically.
    }

    Invoke-WebRequest -Uri $AssetUrl -OutFile $DownloadPath -UseBasicParsing

    if (-not (Test-Path -LiteralPath $DownloadPath -PathType Leaf)) {
        Fail-Installer 'The download did not produce a file.'
    }

    $downloadInfo = Get-Item -LiteralPath $DownloadPath
    if ($downloadInfo.Length -ne $AssetSize) {
        Fail-Installer ("The downloaded file has an unexpected size ({0} bytes); expected {1}." -f $downloadInfo.Length, $AssetSize)
    }

    $downloadHash = Get-Sha256 -Path $DownloadPath
    if ($downloadHash -ne $AssetSha256) {
        Fail-Installer 'The downloaded file failed SHA-256 verification. It will not be installed.'
    }

    if (-not (Test-PeX64 -Path $DownloadPath)) {
        Fail-Installer 'The downloaded release is not a valid Windows x64 executable. It will not be installed.'
    }

    Write-Stage 4 'Installing CLI...'

    if (Test-Path -LiteralPath $StagePath) {
        Remove-Item -LiteralPath $StagePath -Force
    }

    Copy-Item -LiteralPath $DownloadPath -Destination $StagePath -Force

    $stageHash = Get-Sha256 -Path $StagePath
    if ($stageHash -ne $AssetSha256) {
        Fail-Installer 'The staged executable failed SHA-256 verification. The previous installation was left unchanged.'
    }

    if (Test-Path -LiteralPath $BackupPath) {
        Remove-Item -LiteralPath $BackupPath -Force
    }

    if (Test-Path -LiteralPath $InstallPath -PathType Leaf) {
        Move-Item -LiteralPath $InstallPath -Destination $BackupPath -Force
        $BackupCreated = $true
    }

    try {
        Move-Item -LiteralPath $StagePath -Destination $InstallPath -Force
        $RollbackNeeded = $true
    }
    catch {
        if ($BackupCreated -and (Test-Path -LiteralPath $BackupPath -PathType Leaf)) {
            Move-Item -LiteralPath $BackupPath -Destination $InstallPath -Force
            $BackupCreated = $false
        }
        throw
    }

    Write-Stage 5 'Configuring PATH...'
    $pathChanged = $false
    try {
        $pathChanged = Add-UserPath -Directory $InstallDirectory
    }
    catch {
        Write-Warning ("Could not update the user PATH automatically. Add '{0}' to your user PATH manually." -f $InstallDirectory)
    }

    if ($pathChanged) {
        Write-Host ("Added {0} to the user PATH." -f $InstallDirectory)
    }
    else {
        Write-Host 'User PATH already contains the CLI installation directory.'
    }

    Write-Stage 6 'Verifying installation...'

    if (-not (Test-Path -LiteralPath $InstallPath -PathType Leaf)) {
        Fail-Installer 'The installed CLI executable is missing.'
    }

    $installedHash = Get-Sha256 -Path $InstallPath
    if ($installedHash -ne $AssetSha256) {
        Fail-Installer 'The installed CLI failed SHA-256 verification.'
    }

    if (-not (Test-PeX64 -Path $InstallPath)) {
        Fail-Installer 'The installed CLI failed its Windows x64 executable check.'
    }

    # Do not launch the TUI automatically. The CLI is interactive and the
    # installer must never block waiting for keyboard input.
    $RollbackNeeded = $false

    if ($BackupCreated -and (Test-Path -LiteralPath $BackupPath -PathType Leaf)) {
        Remove-Item -LiteralPath $BackupPath -Force
        $BackupCreated = $false
    }

    Write-Host ''
    Write-Success 'Installation complete!'
    Write-Host ''
    Write-Host 'Run:'
    Write-Host ''
    Write-Host '    calculate-pi-cli.exe'
    Write-Host ''
    Write-Host 'Open a new PowerShell or Command Prompt window if the command is not available yet.'
    Write-Host ''
    Write-Host "$ProjectName is ready."
    Write-Host ("Released CLI version: {0} | Project author: Amonnium | License: MIT" -f $CliVersion)
}
catch {
    if ($RollbackNeeded -and $BackupCreated -and (Test-Path -LiteralPath $BackupPath -PathType Leaf)) {
        try {
            if (Test-Path -LiteralPath $InstallPath -PathType Leaf) {
                Remove-Item -LiteralPath $InstallPath -Force -ErrorAction SilentlyContinue
            }
            Move-Item -LiteralPath $BackupPath -Destination $InstallPath -Force -ErrorAction SilentlyContinue
            $BackupCreated = $false
        }
        catch {
            Write-Warning 'The previous installation could not be restored automatically.'
        }
    }

    Remove-TemporaryFiles

    Write-Host ''
    Write-Host 'Installation failed.' -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}
finally {
    Remove-TemporaryFiles
}
