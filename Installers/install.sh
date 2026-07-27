#!/bin/bash

set -e

# This is the Calculate Pi Scripts Installer. It's an all-in-one script that clones the repo with Git, goes into the directory, executes the Dependencies Installer and removes .git and the files only compatible on Windows.
# This is the recommended way for those who only want to have the languages and the files installed (normal users), without the bloat of a .git directory, README.md, LICENSE, and batch files, PowerShell files and .exe executables.
# Only use this on Bash shell, present on most Unix-like systems.
# © 2026 Amonnium

clear

echo "Welcome to Calculate Pi Scripts Installer"
sleep 1
echo "Hosted on GitHub since 2026"
echo "Made with dedication and love by Amonnium."
sleep 2
echo "Starting installation..."
sleep 2

# Clones the repository from GitHub with Git and goes into the directory.
echo
echo "----------------------------------------"
echo "Cloning the GitHub repository"
echo "----------------------------------------"
sleep 0.5

git clone https://github.com/Amonnium/Calculate-Pi-Scripts.git
cd Calculate-Pi-Scripts/Dependencies/

# Running the Dependencies Installer.

echo
echo "----------------------------------------"
echo "Running the Dependencies Installer."
echo "----------------------------------------"
sleep 0.5

chmod +x install.sh
./install.sh
sleep 0.2
cd ..

# Removing the bloat files.

echo
echo "----------------------------------------"
echo "Removing the bloat files."
echo "----------------------------------------"
sleep 0.5

rm -rf .git
rm README.md
rm LICENSE
rm Pi.bat
rm -r Pi_in_PowerShell
rm -r Assets
rm -r Dependencies
rm -r Installers
rm Pi_in_C++_Folder/Pi_in_C++.exe
rm Pi_in_C_Folder/Pi_in_C.exe
rm Pi_in_Go_Folder/Math_Pi_in_Go.exe
rm Pi_in_Go_Folder/Pi_in_Go.exe
rm Pi_in_Rust_Folder/Pi_in_Rust.exe
rm Pi_in_Rust_Folder/Math_Pi_in_Rust.exe

echo
echo "----------------------------------------"
echo "Installation completed successfully!"
echo
sleep 1
echo "You may need to restart your terminal before using the Terminal."
echo
sleep 1
echo "Thank you for installing Calculate Pi Scripts!"
echo "- Amonnium"
echo "----------------------------------------"
