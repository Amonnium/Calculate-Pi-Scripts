#!/bin/bash

set -e

# This Bash script will automatically install all the dependencies needed to run the Calculate Pi Scripts.
# This was all written by Amonnium, and is only recommended for Debian-based systems because it uses APT, cURL and Snap.
# © 2026 Amonnium

clear

echo "Dependencies Installer by Calculate Pi Scripts"
sleep 1
echo "Hosted on GitHub since 2026"
echo "Made with dedication and love by Amonnium."
sleep 2
echo "Starting installation..."
sleep 2

# Check if the system is Debian-based
if ! command -v apt >/dev/null 2>&1; then
    echo
    echo "This installer only supports Debian-based Linux distributions."
    exit 1
fi

# Updating packages
echo
echo "----------------------------------------"
echo "Updating packages..."
echo "----------------------------------------"
sleep 0.5

sudo apt update && sudo apt upgrade -y

# Python installation
echo
echo "----------------------------------------"
echo "Installing Python..."
echo "----------------------------------------"
sleep 0.5

if command -v python3 >/dev/null 2>&1; then
    echo "Python is already installed."
else
    sudo apt install -y python3 python-is-python3 python3-pip
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
    sudo apt install -y build-essential
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
    sudo apt install -y golang-go
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
    sudo apt install -y default-jdk
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
    sudo apt install -y lua5.4
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
    sudo apt install -y perl
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
    sudo apt install -y ruby-full
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
    sudo snap install zig --classic --edge
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
