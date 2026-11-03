# Calculate Pi Scripts CLI Uninstaller - Windows

[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$AssetName = 'calculate-pi-cli.exe'

if ($env:OS -ne 'Windows_NT') {
    throw 'This uninstaller supports Windows only.'
}

if ([string]::IsNullOrWhiteSpace($env:LOCALAPPDATA)) {
    throw 'LOCALAPPDATA is not available. The CLI installation location cannot be determined safely.'
}

$InstallDirectory = Join-Path $env:LOCALAPPDATA 'Programs\Calculate-Pi-Scripts\bin'
$InstallPath = Join-Path $InstallDirectory $AssetName

function Test-InstallDirectoryPathEntry {
    param(
        [Parameter(Mandatory = $true)][AllowEmptyString()][string]$Entry,
        [Parameter(Mandatory = $true)][string]$Directory
    )

    $normalizedEntry = [Environment]::ExpandEnvironmentVariables($Entry.Trim().Trim('"')).TrimEnd([char[]]@('\', '/'))
    $normalizedDirectory = $Directory.TrimEnd([char[]]@('\', '/'))

    return [string]::Equals($normalizedEntry, $normalizedDirectory, [StringComparison]::OrdinalIgnoreCase)
}

if (-not (Test-Path -LiteralPath $InstallPath -PathType Leaf)) {
    throw "CLI executable not found at '$InstallPath'. Nothing was uninstalled."
}

Remove-Item -LiteralPath $InstallPath -Force
Write-Host ("Removed CLI executable: {0}" -f $InstallPath) -ForegroundColor Green

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$pathEntryRemoved = $false

if (-not [string]::IsNullOrEmpty($userPath)) {
    $userPathEntries = @($userPath -split ';')
    $remainingUserPathEntries = @(
        $userPathEntries | Where-Object {
            -not (Test-InstallDirectoryPathEntry -Entry $_ -Directory $InstallDirectory)
        }
    )

    if ($remainingUserPathEntries.Count -lt $userPathEntries.Count) {
        [Environment]::SetEnvironmentVariable('Path', ($remainingUserPathEntries -join ';'), 'User')
        $pathEntryRemoved = $true
    }
}

$processPathEntries = @($env:Path -split ';')
$remainingProcessPathEntries = @(
    $processPathEntries | Where-Object {
        -not (Test-InstallDirectoryPathEntry -Entry $_ -Directory $InstallDirectory)
    }
)

if ($remainingProcessPathEntries.Count -lt $processPathEntries.Count) {
    $env:Path = $remainingProcessPathEntries -join ';'
}

if ($pathEntryRemoved) {
    Write-Host ("Removed the CLI installation directory from the user PATH: {0}" -f $InstallDirectory) -ForegroundColor Green
}
else {
    Write-Host 'The user PATH did not contain the CLI installation directory.'
}

Write-Host 'Calculate Pi Scripts CLI uninstallation complete.' -ForegroundColor Green
