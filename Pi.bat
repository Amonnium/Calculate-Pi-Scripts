@echo off
set /a pi=4
set /a denom=3
set /a osc=-1
set /a count=0
:loop
    cls
    set /a count+=1
    set /a osc=-1^%count%
    if %osc% equ -1 (set /a pi=%pi%-4/%denom%)
    if %osc% equ 1 (set /a pi=%pi%+4/%denom%)
    set /a denom+=2
    echo %pi%
    pause
    goto loop   