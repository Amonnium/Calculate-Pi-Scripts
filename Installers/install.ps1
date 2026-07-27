# Calculate Pi Scripts Installer
# This PowerShell script clones the Calculate Pi Scripts repository,
# installs all required dependencies and removes development files.
# Written by Amonnium
# © 2026 Amonnium

# Requires PowerShell 7+

$ErrorActionPreference = "Stop"

Clear-Host

Write-Host "Welcome to Calculate Pi Scripts Installer"
Start-Sleep -Seconds 1

Write-Host "Hosted on GitHub since 2026"
Write-Host "Made with dedication and love by Amonnium."

Start-Sleep -Seconds 2

Write-Host "Starting installation..."
Start-Sleep -Seconds 2


# --------------------------------------------------
# Check Administrator privileges
# --------------------------------------------------

$currentIdentity = [Security.Principal.WindowsIdentity]::GetCurrent()

$currentPrincipal = New-Object Security.Principal.WindowsPrincipal($currentIdentity)

if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {

    Write-Host ""
    Write-Host "This installer must be executed as Administrator." -ForegroundColor Red
    exit 1

}


# --------------------------------------------------
# Check PowerShell version
# --------------------------------------------------

if ($PSVersionTable.PSVersion.Major -lt 7) {

    Write-Host ""
    Write-Host "PowerShell 7 or newer is required." -ForegroundColor Red
    exit 1

}


# --------------------------------------------------
# Helper functions
# --------------------------------------------------

function Write-Section {

    param(
        [string]$Message
    )

    Write-Host ""
    Write-Host "----------------------------------------"
    Write-Host $Message
    Write-Host "----------------------------------------"

    Start-Sleep -Milliseconds 500

}


function Command-Exists {

    param(
        [string]$Command
    )

    return $null -ne (Get-Command $Command -ErrorAction SilentlyContinue)

}


function Refresh-Environment {

    $machinePath = [System.Environment]::GetEnvironmentVariable(
        "Path",
        "Machine"
    )

    $userPath = [System.Environment]::GetEnvironmentVariable(
        "Path",
        "User"
    )

    $env:Path = "$machinePath;$userPath"

}


function Install-Git {

    if (Command-Exists "git") {

        Write-Host "Git is already installed."
        return

    }

    Write-Host "Installing Git..."

    winget install `
        --id Git.Git `
        --silent `
        --accept-package-agreements `
        --accept-source-agreements

    Refresh-Environment

}


# --------------------------------------------------
# Check Winget
# --------------------------------------------------

if (-not (Command-Exists "winget")) {

    Write-Host ""
    Write-Host "Winget was not found." -ForegroundColor Red
    Write-Host "Please install App Installer from Microsoft Store."

    exit 1

}


# --------------------------------------------------
# Git installation
# --------------------------------------------------

Write-Section "Installing Git"

Install-Git

# --------------------------------------------------
# Clone the GitHub repository
# --------------------------------------------------

Write-Section "Cloning the GitHub repository"

$RepositoryUrl = "https://github.com/Amonnium/Calculate-Pi-Scripts.git"
$RepositoryName = "Calculate-Pi-Scripts"

try {

    if (Test-Path $RepositoryName) {

        Write-Host "The repository already exists."
        Write-Host "Using the existing directory."

    }
    else {

        git clone $RepositoryUrl

    }

}
catch {

    Write-Host ""
    Write-Host "Failed to clone the repository." -ForegroundColor Red
    Write-Host $_.Exception.Message

    exit 1

}


# --------------------------------------------------
# Enter the Dependencies directory
# --------------------------------------------------

Write-Section "Opening the Dependencies folder"

try {

    Push-Location (Join-Path $RepositoryName "Dependencies")

}
catch {

    Write-Host ""
    Write-Host "Unable to access the Dependencies directory." -ForegroundColor Red
    Write-Host $_.Exception.Message

    exit 1

}


# --------------------------------------------------
# Run the Dependencies Installer
# --------------------------------------------------

Write-Section "Running the Dependencies Installer"

try {

    $DependencyInstaller = Join-Path (Get-Location) "install.ps1"

    if (-not (Test-Path $DependencyInstaller)) {

        throw "Dependencies/install.ps1 was not found."

    }

    & pwsh -ExecutionPolicy Bypass -File $DependencyInstaller

}
catch {

    Write-Host ""
    Write-Host "The Dependencies Installer failed." -ForegroundColor Red
    Write-Host $_.Exception.Message

    Pop-Location

    exit 1

}


# --------------------------------------------------
# Return to the repository root
# --------------------------------------------------

Pop-Location


# --------------------------------------------------
# Preparing file cleanup
# --------------------------------------------------

Write-Section "Removing the development files"

$RepositoryRoot = Join-Path (Get-Location) $RepositoryName

$FilesToRemove = @(

    ".git",
    "README.md",
    "LICENSE",
    "Pi.bat",
    "Installers",
    "Pi_in_PowerShell",
    "Assets",
    "Dependencies\README.md",
    "Dependencies\install.sh",
    "Pi_in_C++_Folder\Pi_in_C++.exe",
    "Pi_in_C_Folder\Pi_in_C.exe",
    "Pi_in_Go_Folder\Math_Pi_in_Go.exe",
    "Pi_in_Go_Folder\Pi_in_Go.exe",
    "Pi_in_Rust_Folder\Math_Pi_in_Rust.exe",
    "Pi_in_Rust_Folder\Pi_in_Rust.exe"

)

# --------------------------------------------------
# Remove development files
# --------------------------------------------------

foreach ($Item in $FilesToRemove) {

    $Target = Join-Path $RepositoryRoot $Item

    if (Test-Path $Target) {

        try {

            Remove-Item `
                -Path $Target `
                -Recurse `
                -Force `
                -ErrorAction Stop

            Write-Host "Removed: $Item"

        }
        catch {

            Write-Host "Unable to remove: $Item" -ForegroundColor Yellow

        }

    }

}


Start-Sleep -Seconds 2


# --------------------------------------------------
# Installation completed
# --------------------------------------------------

Write-Host ""
Write-Host "----------------------------------------"
Write-Host "Installation completed successfully!"
Write-Host ""
Start-Sleep -Seconds 1

Write-Host "You may need to restart your terminal before using some of the newly installed languages."
Write-Host ""
Start-Sleep -Seconds 1

Write-Host "Thank you for installing Calculate Pi Scripts!"
Write-Host "- Amonnium"
Write-Host "----------------------------------------"

exit 0
