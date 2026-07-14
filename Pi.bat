@echo off
setlocal

:: 1. Create a temporary VBScript file
echo Dim pi > temp_pi.vbs
echo pi = 4 * Atn(1) >> temp_pi.vbs
echo WScript.Echo pi >> temp_pi.vbs

:: 2. Execute the VBScript and capture its output to a Batch variable
for /f "tokens=*" %%i in ('cscript //nologo temp_pi.vbs') do set "PI=%%i"

:: 3. Clean up the temporary VBScript file
del temp_pi.vbs

:: 4. Display the results
echo Calculated Pi is: %PI%

:: Example usage: calculate the area of a circle with radius 5
:: (Batch handles decimals in display only, not in /a calculations)
set "radius=5"
echo.
echo Radius: %radius%
cscript //nologo -e:vbs -b -nologo "ExecuteGlobal \"r = " & radius & ": pi = 4 * Atn(1) : WScript.Echo (pi * r^2)\""

endlocal
pause
