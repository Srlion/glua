# glua

## Installation

### Windows

To use `cgo` on Windows, the simplest method is to install a GCC compiler using [MSYS2](https://www.msys2.org/). Follow the steps below to set up both 32-bit and 64-bit GCC toolchains:

1. **Download and Install MSYS2:**

   - Visit the [MSYS2 website](https://www.msys2.org/) and download the installer.
   - Run the installer and follow the on-screen instructions to complete the installation.

2. **Update Package Database and Core Packages:**

   - Open the **MSYS2 MSYS** terminal from the Start Menu.
   - Update the package database and core system packages by executing:
     ```bash
     pacman -Syu
     ```
   - If prompted to close the terminal, do so, then reopen it and run the update command again to ensure all packages are up to date:
     ```bash
     pacman -Su
     ```

3. **Install the GCC Toolchains for 32-bit and 64-bit:**

   - In the MSYS2 terminal, install both 64-bit and 32-bit GCC compilers by running:
     ```bash
     pacman -S mingw-w64-x86_64-gcc mingw-w64-i686-gcc
     ```
   - This command installs both the 64-bit (`mingw-w64-x86_64-gcc`) and 32-bit (`mingw-w64-i686-gcc`) GCC compilers.

4. **Add MSYS2 GCC Binaries to Your PATH:**

   - To make `gcc` accessible from the Windows Command Prompt and other environments, add both MSYS2 `mingw64\bin` and `mingw32\bin` directories to your `PATH` environment variable.
     - **Typical Paths:**
       - `C:\msys64\mingw64\bin` (64-bit)
       - `C:\msys64\mingw32\bin` (32-bit)
   - **Steps to Add to PATH:**
     1. Press `Win + S`, type **Environment Variables**, and select **Edit the system environment variables**.
     2. In the **System Properties** window, click on the **Environment Variables...** button.
     3. Under **System variables**, find and select the `Path` variable, then click **Edit**.
     4. Click **New** and add both paths:
        - `C:\msys64\mingw64\bin`
        - `C:\msys64\mingw32\bin`
     5. Click **OK** to save the changes.

5. **Verify GCC Installations:**
   - Open the Command Prompt and run the following commands to verify both GCC versions:
     ```bash
     gcc --version
     ```
     - This should display the GCC version information for the default compiler in your PATH.
   - To specifically check the 64-bit GCC:
     ```bash
     gcc -m64 --version
     ```
   - To specifically check the 32-bit GCC:
     ```bash
     gcc -m32 --version
     ```
   - Ensure that both commands return the correct GCC version information, confirming that both 32-bit and 64-bit GCC compilers are installed and accessible.

### Linux

First, update your package list and install the required GCC packages:

```bash
sudo apt-get update
sudo apt-get install gcc-i686-linux-gnu
sudo apt-get install gcc-multilib
```

### Install glua

Once the necessary GCC setup is complete, you can install `glua` using the following command:

```bash
go get -u github.com/Srlion/glua
```

## Usage

```go
package main

import (
	"github.com/Srlion/glua"
)

func init() {
	glua.GMOD13_OPEN = gmod13_open
	glua.GMOD13_CLOSE = gmod13_close
}

func test(L glua.State) int {
	var str string = L.CheckString(1)
	println("string:", str)

	L.PushString("Hello from Go!")
	return 1
}

func gmod13_open(L glua.State) int {
	println("hello from gmod13_open!")

	L.PushGoFunc(test)
	L.SetGlobal("test")

	return 0
}

func gmod13_close(L glua.State) int {
	return 0
}

// Required by Go when using `-buildmode=c-shared`
func main() {}

```
