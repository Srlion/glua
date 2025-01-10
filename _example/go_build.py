#!/usr/bin/env python3
import argparse
import os
import sys
import subprocess
import multiprocessing

def main():
    # ----------------------------------------
    # Default values
    # ----------------------------------------
    DEFAULT_ARCH = 32
    DEFAULT_CFLAGS = ""
    DEFAULT_LDFLAGS = ""
    DEFAULT_DIRECTORY = "."
    DEFAULT_MAXPROCS = 1

    # ----------------------------------------
    # Argument parser
    # ----------------------------------------
    parser = argparse.ArgumentParser(
        description="Build a Go shared library with configurable flags."
    )

    parser.add_argument(
        "-C", "--directory",
        type=str,
        default=DEFAULT_DIRECTORY,
        help="Change to directory before executing commands (default: current directory)"
    )
    parser.add_argument(
        "-a", "--arch",
        type=int,
        default=DEFAULT_ARCH,
        choices=[32, 64],
        help="Architecture (32 or 64). Default: 32"
    )
    parser.add_argument(
        "-f", "--cflags",
        type=str,
        default=DEFAULT_CFLAGS,
        help="C compiler flags. Default: empty"
    )
    parser.add_argument(
        "-l", "--ldflags",
        type=str,
        default=DEFAULT_LDFLAGS,
        help="Linker flags. Default: empty"
    )
    parser.add_argument(
        "-n", "--name",
        type=str,
        required=True,
        help="Output executable/DLL name (required)."
    )
    parser.add_argument(
        "-p", "--maxprocs",
        type=int,
        default=DEFAULT_MAXPROCS,
        help="Maximum parallel processes. Default: number of CPU cores"
    )
    parser.add_argument(
        "-o", "--outdir",
        type=str,
        default="bin",
        help="Output directory for the built DLL (default: ./bin)"
    )

    args = parser.parse_args()

    # ----------------------------------------
    # Extract parsed arguments
    # ----------------------------------------
    directory = args.directory
    arch = args.arch
    cflags = args.cflags
    ldflags = args.ldflags
    name = args.name
    maxprocs = args.maxprocs
    outdir = args.outdir
    if not os.path.exists(outdir):
        os.makedirs(outdir, exist_ok=True)

    # ----------------------------------------
    # Additional validations (some already handled by argparse choices)
    # ----------------------------------------
    if not name.strip():
        print("Error: Output name cannot be empty.")
        parser.print_help()
        sys.exit(1)

    # Make sure the directory exists
    if not os.path.isdir(directory):
        print(f"Error: Directory '{directory}' does not exist.")
        sys.exit(1)

    # ----------------------------------------
    # Change to the target directory
    # ----------------------------------------
    old_pwd = os.getcwd()
    try:
        os.chdir(directory)
    except Exception as e:
        print(f"Error: Failed to change directory to '{directory}': {e}")
        sys.exit(1)

    # ----------------------------------------
    # Set environment variables
    # ----------------------------------------
    os.environ["GOOS"] = "windows" if os.name == 'nt' else "linux"
    os.environ["CGO_ENABLED"] = "1"
    os.environ["GOAMD64"] = "v3"
    os.environ["GOMAXPROCS"] = str(maxprocs)
    os.environ["CGO_CFLAGS"] = f"-Ofast -fvisibility=hidden -flto {cflags}"
    LDFLAGS = f"-s -w -extldflags '-Ofast -fvisibility=hidden -flto {ldflags}'"

    path_sep = ';' if os.name == 'nt' else ':'

    if arch == 32:
        if os.name == 'nt':
            os.environ["PATH"] = f"C:/mingw32/bin" + path_sep + os.environ["PATH"]
        os.environ["GOARCH"] = "386"
        dll_arch = "win32" if os.name == 'nt' else "linux"
    else:
        if os.name == 'nt':
            os.environ["PATH"] = f"C:/mingw64/bin" + path_sep + os.environ["PATH"]
        os.environ["GOARCH"] = "amd64"
        dll_arch = "win64" if os.name == 'nt' else "linux64"

    dll_name = os.path.join(outdir, f"gmsv_{name}_{dll_arch}.dll")

    # ----------------------------------------
    # Display configuration
    # ----------------------------------------
    print("========================================")
    print("Build Configuration:")
    print("----------------------------------------")
    print(f"Architecture  : {arch}-bit")
    print(f"CFLAGS        : {os.environ['CGO_CFLAGS']}")
    print(f"LDFLAGS       : {LDFLAGS}")
    print(f"DLL File      : {dll_name}")
    print(f"Max Procs     : {maxprocs}")
    print("========================================")
    print()

    # ----------------------------------------
    # Build Process
    # ----------------------------------------
    os.makedirs("bin", exist_ok=True)
    build_cmd = [
        "go", "build",
        "-buildmode=c-shared",
        "-o", dll_name,
        "-ldflags", LDFLAGS,
    ]

    print("Running build command:", " ".join(build_cmd))
    try:
        subprocess.run(build_cmd, check=True)
    except subprocess.CalledProcessError as e:
        print("Compilation failed.")
        sys.exit(1)

    print(f"Compilation succeeded. Output: {dll_name}")

    # ----------------------------------------
    # Return to the old directory
    # ----------------------------------------
    try:
        os.chdir(old_pwd)
    except Exception as e:
        print(f"Error: Failed to change directory back to '{old_pwd}': {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
