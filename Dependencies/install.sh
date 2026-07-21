#!/bin/bash

set -e

# This Bash script will automatically install all the dependencies needed to run the Calculate Pi Scripts.
# This was all written by Amonnium, and it should only be used on Bash.
# © 2026 Amonnium

clear

echo "Dependencies Installer by Calculate Pi Scripts"
sleep 1
echo "Hosted on GitHub since 2026"
echo "Made with dedication and love by Amonnium."
sleep 2
echo "Starting installation..."
sleep 2

# Check what system is.
source /etc/os-release

case "$ID" in
    ubuntu|debian|linuxmint|pop)
        PKG_MANAGER="apt"
        ;;
    fedora|rhel|rocky|almalinux)
        PKG_MANAGER="dnf"
        ;;
    arch|manjaro|endeavouros|cachyos)
        PKG_MANAGER="pacman"
        ;;
    opensuse*|opensuse-leap|opensuse-tumbleweed)
        PKG_MANAGER="zypper"
        ;;
    *)
        echo "Unsupported Linux distribution."
        exit 1
        ;;
esac

# Updating packages
echo
echo "----------------------------------------"
echo "Updating packages..."
echo "----------------------------------------"
sleep 0.5

update_system() {
    case "$PKG_MANAGER" in
        apt)
            sudo apt update && sudo apt upgrade -y
            ;;
        dnf)
            sudo dnf upgrade -y
            ;;
        pacman)
            sudo pacman -Syu --noconfirm
            ;;
        zypper)
            sudo zypper refresh
            sudo zypper update -y
            ;;
    esac
}

# Defining package manager settings.
install_package() {
    case "$PKG_MANAGER" in
        apt)
            sudo apt install -y "$@"
            ;;
        dnf)
            sudo dnf install -y "$@"
            ;;
        pacman)
            sudo pacman -S --noconfirm "$@"
            ;;
        zypper)
            sudo zypper install -y "$@"
            ;;
    esac
}

# Git installation
echo
echo "----------------------------------------"
echo "Installing Git..."
echo "----------------------------------------"
sleep 0.5

if command -v git >/dev/null 2>&1; then
    echo "Git is already installed."
else
    case "$PKG_MANAGER" in
        apt)
            install_package git
            ;;
        dnf)
            install_package git
            ;;
        pacman)
            install_package git
            ;;
        zypper)
            install_package git
            ;;
    esac
fi

# Python installation
echo
echo "----------------------------------------"
echo "Installing Python..."
echo "----------------------------------------"
sleep 0.5

if command -v python3 >/dev/null 2>&1; then
    echo "Python is already installed."
else
    case "$PKG_MANAGER" in

    apt)
        install_package python3 python-is-python3 python3-pip
        ;;

    dnf)
        install_package python3 python3-pip
        ;;

    pacman)
        install_package python python-pip
        ;;

    zypper)
        install_package python3 python3-pip
        ;;

    esac
fi

# GCC and G++ installation
echo
echo "----------------------------------------"
echo "Installing GCC and G++..."
echo "----------------------------------------"
sleep 0.5

if command -v gcc >/dev/null 2>&1 && command -v g++ >/dev/null 2>&1; then
    echo "GCC and G++ are already installed."
else
    case "$PKG_MANAGER" in

apt)
    install_package build-essential
    ;;

dnf)
    install_package gcc gcc-c++ make
    ;;

pacman)
    install_package base-devel
    ;;

zypper)
    install_package gcc gcc-c++ make
    ;;

    esac
fi

# Go installation
echo
echo "----------------------------------------"
echo "Installing Go..."
echo "----------------------------------------"
sleep 0.5

if command -v go >/dev/null 2>&1; then
    echo "Go is already installed."
else
    case "$PKG_MANAGER" in
        apt)
            install_package golang-go
            ;;
        dnf)
            install_package golang
            ;;
        pacman)
            install_package go
            ;;
        zypper)
            install_package go
            ;;
    esac
fi

# Java installation
echo
echo "----------------------------------------"
echo "Installing JDK..."
echo "----------------------------------------"
sleep 0.5

if command -v java >/dev/null 2>&1; then
    echo "Java is already installed."
else
    case "$PKG_MANAGER" in
        apt)
            install_package default-jdk
            ;;
        dnf)
            install_package java-latest-openjdk-devel
            ;;
        pacman)
            install_package jdk-openjdk
            ;;
        zypper)
            install_package java-21-openjdk-devel
            ;;
    esac
fi

# Julia installation
echo
echo "----------------------------------------"
echo "Installing Julia..."
echo "----------------------------------------"
sleep 0.5

if command -v julia >/dev/null 2>&1; then
    echo "Julia is already installed."
else
    curl -fsSL https://install.julialang.org | sh
fi

# Kotlin installation (without IntelliJ)
echo
echo "----------------------------------------"
echo "Installing Kotlin (without IntelliJ)..."
echo "----------------------------------------"
sleep 0.5

if command -v kotlin >/dev/null 2>&1; then
    echo "Kotlin is already installed."
else
    curl -s https://get.sdkman.io | bash

    source "$HOME/.sdkman/bin/sdkman-init.sh"

    sdk install kotlin
fi

# Rust installation
echo
echo "----------------------------------------"
echo "Installing Rust..."
echo "----------------------------------------"
sleep 0.5

if command -v rustc >/dev/null 2>&1; then
    echo "Rust is already installed."
else
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
fi

# Node.js installation
echo
echo "----------------------------------------"
echo "Installing Node.js..."
echo "----------------------------------------"
sleep 0.5

if command -v node >/dev/null 2>&1; then
    echo "Node.js is already installed."
else
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

    export NVM_DIR="$HOME/.nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"

    nvm install --lts
fi

# Lua installation
echo
echo "----------------------------------------"
echo "Installing Lua..."
echo "----------------------------------------"
sleep 0.5

if command -v lua >/dev/null 2>&1; then
    echo "Lua is already installed."
else
    case "$PKG_MANAGER" in
        apt)
            install_package lua5.4
            ;;
        dnf)
            install_package lua
            ;;
        pacman)
            install_package lua
            ;;
        zypper)
            install_package lua
            ;;
    esac
fi

# Perl installation
echo
echo "----------------------------------------"
echo "Installing Perl..."
echo "----------------------------------------"
sleep 0.5

if command -v perl >/dev/null 2>&1; then
    echo "Perl is already installed."
else
    case "$PKG_MANAGER" in
        apt)
            install_package perl
            ;;
        dnf)
            install_package perl
            ;;
        pacman)
            install_package perl
            ;;
        zypper)
            install_package perl
            ;;
    esac
fi

# Ruby installation
echo
echo "----------------------------------------"
echo "Installing Ruby..."
echo "----------------------------------------"
sleep 0.5

if command -v ruby >/dev/null 2>&1; then
    echo "Ruby is already installed."
else
    case "$PKG_MANAGER" in
        apt)
            install_package ruby-full
            ;;
        dnf)
            install_package ruby
            ;;
        pacman)
            install_package ruby
            ;;
        zypper)
            install_package ruby
            ;;
    esac
fi

# Zig installation
echo
echo "----------------------------------------"
echo "Installing Zig..."
echo "----------------------------------------"
sleep 0.5

if command -v zig >/dev/null 2>&1; then
    echo "Zig is already installed."
else
    case "$PKG_MANAGER" in
        apt)
            install_package zig
            ;;
        dnf)
            install_package zig
            ;;
        pacman)
            install_package zig
            ;;
        zypper)
            install_package zig
            ;;
    esac
fi

# Raku installation
echo
echo "----------------------------------------"
echo "Installing Raku..."
echo "----------------------------------------"
sleep 0.5

if command -v raku >/dev/null 2>&1; then
    echo "Raku is already installed."
else
    case "$PKG_MANAGER" in
        apt)
            install_package rakudo
            ;;
        dnf)
            install_package rakudo
            ;;
        pacman)
            install_package rakudo
            ;;
        zypper)
            install_package raku
            ;;
    esac
fi

sleep 2

echo
echo "----------------------------------------"
echo "Installation completed successfully!"
echo
sleep 1
echo "You may need to restart your terminal before using some of the newly installed languages."
echo
sleep 1
echo "Thank you for using this utility!"
echo "- Amonnium"
echo "----------------------------------------"
