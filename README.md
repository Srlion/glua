# glua

> [!NOTE]
> If you are using an editor, you may have lots of errors, this is due to how cgo works. You can ignore these errors, as they will not affect the final build.
> You can use `_example/go_build.py` to build your module and it will just work. I asumme you have python installed.

## Installation

### Windows

> [!NOTE]
> Previously, I suggested using MySys2 to install, however I have found that it causes crashes when restarting the server. I recommend using WinLibs MinGW-w64 instead.

1. **Download WinLibs MinGW-w64:**

   - Visit the [WinLibs website](https://winlibs.com/) and download the latest standalone MinGW-w64 builds:
     - **64-bit GCC**: Download the `mingw64` archive.
     - **32-bit GCC**: Download the `mingw32` archive.

2. **Extract the Archives:**

   - Extract each `.zip` archive to:
     - `C:\mingw64`.
     - `C:\mingw32`.

3. **Add MinGW-w64 Binaries to Your PATH:**

   - **Steps to Add to PATH:**
     1. Press `Win + S`, type **Environment Variables**, and select **Edit the system environment variables**.
     2. In the **System Properties** window, click on the **Environment Variables...** button.
     3. Under **System variables**, find and select the `Path` variable, then click **Edit**.
     4. Click **New** and add both paths:
        - `C:\mingw64\bin`
        - `C:\mingw32\bin`
     5. Click **OK** to save the changes.

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
