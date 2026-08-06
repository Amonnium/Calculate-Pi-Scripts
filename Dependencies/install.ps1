# Calculate Pi Scripts - Windows Dependencies Installer
# This PowerShell script automatically installs all dependencies required to run Calculate Pi Scripts.
# Written by Amonnium
# © 2026 Amonnium

# Requires PowerShell 7+

$ErrorActionPreference = "Stop"

Clear-Host

Write-Host "Dependencies Installer by Calculate Pi Scripts"
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
    Write-Host "Please restart PowerShell as Administrator and try again."
    exit 1
}


# --------------------------------------------------
# Check PowerShell version
# --------------------------------------------------

if ($PSVersionTable.PSVersion.Major -lt 7) {

    Write-Host ""
    Write-Host "PowerShell 7 or newer is required." -ForegroundColor Red
    Write-Host "Current version: $($PSVersionTable.PSVersion)"
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


function Install-Winget {

    param(
        [string]$Name,
        [string]$Id,
        [string]$Command
    )


    if (Command-Exists $Command) {

        Write-Host "$Name is already installed."

    }
    else {

        Write-Host "Installing $Name..."

        winget install `
            --id $Id `
            --silent `
            --accept-package-agreements `
            --accept-source-agreements

        Refresh-Environment
    }
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
# Scoop installation
# --------------------------------------------------

Write-Section "Installing Scoop"

if (Command-Exists "scoop") {

    Write-Host "Scoop is already installed."

}
else {

    Write-Host "Installing Scoop..."

    Set-ExecutionPolicy Bypass -Scope Process -Force
    iwr -useb get.scoop.sh | iex

    Refresh-Environment
}


# --------------------------------------------------
# Git installation
# --------------------------------------------------

Write-Section "Installing Git"

if (Command-Exists "git") {

    Write-Host "Git is already installed."

}
else {

    Install-Winget `
        -Name "Git" `
        -Id "Git.Git" `
        -Command "git"

}


# --------------------------------------------------
# Python installation
# --------------------------------------------------

Write-Section "Installing Python"

if (Command-Exists "python") {

    Write-Host "Python is already installed."

}
else {

    Install-Winget `
        -Name "Python" `
        -Id "Python.Python.3" `
        -Command "python"

}


# --------------------------------------------------
# MSYS2 installation
# --------------------------------------------------

Write-Section "Installing GCC, G++ and Make (MSYS2)"

if (
    (Command-Exists "gcc") -and
    (Command-Exists "g++") -and
    (Command-Exists "make")
) {

    Write-Host "GCC, G++ and Make are already installed."

}
else {

    Install-Winget `
        -Name "MSYS2" `
        -Id "MSYS2.MSYS2" `
        -Command "pacman"


    Refresh-Environment


    $msysBash = "C:\msys64\usr\bin\bash.exe"


    if (Test-Path $msysBash) {

        Write-Host "Updating MSYS2 packages..."

        & $msysBash -lc "pacman -Syu --noconfirm"


        Write-Host "Installing build tools..."

        & $msysBash -lc `
            "pacman -S --needed --noconfirm mingw-w64-ucrt-x86_64-gcc mingw-w64-ucrt-x86_64-g++ make"

    }
    else {

        Write-Host "MSYS2 was installed but the installation path was not found." -ForegroundColor Yellow

    }
}
# --------------------------------------------------
# Go installation
# --------------------------------------------------

Write-Section "Installing Go"

if (Command-Exists "go") {

    Write-Host "Go is already installed."

}
else {

    Install-Winget `
        -Name "Go" `
        -Id "GoLang.Go" `
        -Command "go"

}


# --------------------------------------------------
# Java installation
# --------------------------------------------------

Write-Section "Installing Java OpenJDK"

if (Command-Exists "java") {

    Write-Host "Java is already installed."

}
else {

    Install-Winget `
        -Name "OpenJDK" `
        -Id "Microsoft.OpenJDK.21" `
        -Command "java"

}


# --------------------------------------------------
# .NET installation
# --------------------------------------------------

Write-Section "Installing .NET SDK"

if (Command-Exists "dotnet") {

    Write-Host ".NET is already installed."

}
else {

    Install-Winget `
        -Name ".NET SDK" `
        -Id "Microsoft.DotNet.SDK.10" `
        -Command "dotnet"

}


# --------------------------------------------------
# Nim installation
# --------------------------------------------------

Write-Section "Installing Nim"

if (Command-Exists "nim") {

    Write-Host "Nim is already installed."

}
else {

    Write-Host "Installing Nim..."

    scoop install nim

    Refresh-Environment
}


# --------------------------------------------------
# Julia installation
# --------------------------------------------------

Write-Section "Installing Julia"

if (Command-Exists "julia") {

    Write-Host "Julia is already installed."

}
else {

    Install-Winget `
        -Name "Julia" `
        -Id "JuliaLang.Julia" `
        -Command "julia"

}


# --------------------------------------------------
# Kotlin installation
# --------------------------------------------------

Write-Section "Installing Kotlin"

if (Command-Exists "kotlin") {

    Write-Host "Kotlin is already installed."

}
else {

    Install-Winget `
        -Name "Kotlin" `
        -Id "JetBrains.Kotlin" `
        -Command "kotlin"

}


# --------------------------------------------------
# Rust installation
# --------------------------------------------------

Write-Section "Installing Rust"

if (Command-Exists "rustc") {

    Write-Host "Rust is already installed."

}
else {

    Install-Winget `
        -Name "Rust" `
        -Id "Rustlang.Rustup" `
        -Command "rustc"

}


# --------------------------------------------------
# Node.js installation
# --------------------------------------------------

Write-Section "Installing Node.js"

if (Command-Exists "node") {

    Write-Host "Node.js is already installed."

}
else {

    Install-Winget `
        -Name "Node.js" `
        -Id "OpenJS.NodeJS.LTS" `
        -Command "node"

}


# --------------------------------------------------
# Lua installation
# --------------------------------------------------

Write-Section "Installing Lua"

if (Command-Exists "lua") {

    Write-Host "Lua is already installed."

}
else {

    Install-Winget `
        -Name "Lua" `
        -Id "DEVCOM.Lua" `
        -Command "lua"

}


# --------------------------------------------------
# Perl installation
# --------------------------------------------------

Write-Section "Installing Perl"

if (Command-Exists "perl") {

    Write-Host "Perl is already installed."

}
else {

    Install-Winget `
        -Name "Perl" `
        -Id "StrawberryPerl.StrawberryPerl" `
        -Command "perl"

}


# --------------------------------------------------
# Ruby installation
# --------------------------------------------------

Write-Section "Installing Ruby"

if (Command-Exists "ruby") {

    Write-Host "Ruby is already installed."

}
else {

    Install-Winget `
        -Name "Ruby" `
        -Id "RubyInstallerTeam.RubyWithDevKit.3.3" `
        -Command "ruby"

}


# --------------------------------------------------
# Zig installation
# --------------------------------------------------

Write-Section "Installing Zig"

if (Command-Exists "zig") {

    Write-Host "Zig is already installed."

}
else {

    Install-Winget `
        -Name "Zig" `
        -Id "zig.zig" `
        -Command "zig"

}
# --------------------------------------------------
# Raku installation
# --------------------------------------------------

Write-Section "Installing Raku"

if (Command-Exists "raku") {

    Write-Host "Raku is already installed."

}
else {

    Install-Winget `
        -Name "Raku" `
        -Id "Rakudo.Rakudo" `
        -Command "raku"

}


# --------------------------------------------------
# Final environment refresh
# --------------------------------------------------

Refresh-Environment


Start-Sleep -Seconds 2


# --------------------------------------------------
# Installation completed
# --------------------------------------------------

Write-Host ""

Write-Host "----------------------------------------"

Write-Host "Installation completed successfully!"

Write-Host ""

Write-Host "You may need to restart your terminal before using some of the newly installed languages."

Write-Host ""

Write-Host "Thank you for using this utility!"

Write-Host "- Amonnium"

Write-Host "----------------------------------------"