@echo off
setlocal EnableDelayedExpansion

REM Default values
set "arch=32"
set "custom_name="
set "maxprocs="

REM Parse command-line arguments
:parse_args
if "%~1"=="" goto after_args

if "%~1"=="-arch" (
    set "arch=%~2"
    shift
    shift
    goto parse_args
)

if "%~1"=="-name" (
    set "custom_name=%~2"
    shift
    shift
    goto parse_args
)

if "%~1"=="-maxprocs" (
    set "maxprocs=%~2"
    shift
    shift
    goto parse_args
)

echo Unknown option: %~1
exit /b 1

:after_args
REM Check if -name option was provided
if "!custom_name!"=="" (
    echo Error: -name option is required.
    exit /b 1
)

REM Validate the architecture value
if not "!arch!"=="32" if not "!arch!"=="64" (
    echo Error: Invalid architecture "!arch!". Allowed values are 32 or 64.
    exit /b 1
)

REM Set environment variables based on architecture
if "!arch!"=="64" (
    set GOARCH=amd64
    set DLL_NAME=win64
    set PATH=C:\mingw64\bin;%PATH%
) else (
    set GOARCH=386
    set DLL_NAME=win32
    set PATH=C:\mingw32\bin;%PATH%
)

set GOOS=windows
set CGO_ENABLED=1
set GODEBUG=cgocheck=0

REM Set GOMAXPROCS if -maxprocs option was provided
if not "!maxprocs!"=="" (
    set GOMAXPROCS=!maxprocs!
)

REM Ensure the bin directory exists
if not exist bin (
    mkdir bin
)

if defined GOMAXPROCS (
    echo Building gmsv_%custom_name%_%DLL_NAME%.dll with GOARCH=%GOARCH% GOMAXPROCS=%GOMAXPROCS%
) else (
    echo Building gmsv_%custom_name%_%DLL_NAME%.dll with GOARCH=%GOARCH%
)

REM Build the Go program based on architecture and custom name
go build -buildmode=c-shared -o "bin\gmsv_!custom_name!_!DLL_NAME!.dll" 2>&1

if %errorlevel% neq 0 (
    echo Build failed.
) else (
    echo Build succeeded: gmsv_!custom_name!_!DLL_NAME!.dll
)

endlocal
