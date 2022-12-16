package checkenv

import (
	"runtime"

	"golang.org/x/sys/cpu"
)

type OS string

const (
	OSUnknown OS = "unknown"
	OSLinux   OS = "linux"
	OSWindows OS = "windows"
	OSFreeBSD OS = "freebsd"
	OSDarwin  OS = "darwin"
	OSAndroid OS = "android"
)

type Arch string

const (
	ArchUnknown Arch = "unknown"
	Arch386     Arch = "386"
	ArchAMD64   Arch = "amd64"
	ArchARM64   Arch = "arm64"
	ArchARM     Arch = "arm"
	ArchPPC64   Arch = "ppc64"
	ArchPPC64LE Arch = "ppc64le"
	ArchWASM32  Arch = "wasm"
)

type Libc string

const (
	LibcUnknown Libc = "unknown"
	LibcGNU     Libc = "gnu"
	LibcMUSL    Libc = "musl"
	LibcMSVC    Libc = "msvc"
	LibcNone    Libc = "none"
)

type Environment struct {
	OS        OS   `json:"os"`
	Arch      Arch `json:"arch"`
	Libc      Libc `json:"libc"`
	Softfloat bool `json:"softfloat"`
}

var llvm_archmap = map[Arch]string{
	ArchARM64:   "aarch64",
	ArchARM:     "arm",
	Arch386:     "i386",
	ArchPPC64:   "powerpc64",
	ArchPPC64LE: "powerpc64le",
	ArchAMD64:   "x86_64",
	ArchWASM32:  "wasm32",
}

func (e Environment) String() string {
	switch e.OS {
	case OSAndroid:
		// Use linux with android libc
		if llvm_arch, ok := llvm_archmap[e.Arch]; ok {
			switch e.Arch {
			case ArchARM:
				return llvm_arch + "-" + string(OSLinux) + "-" + "androideabi"
			default:
				return llvm_arch + "-" + string(OSLinux) + "-" + "android"
			}
		}
	}

	switch e.Arch {
	case ArchARM:
		switch e.OS {
		case OSLinux:
			switch e.Libc {
			case LibcGNU:
				if e.Softfloat {
					return "arm-linux-gnueabi"
				}
				return "arm-linux-gnueabihf"
			case LibcMUSL:
				if e.Softfloat {
					return "arm-linux-musleabi"
				}
				return "arm-linux-musleabihf"
			case LibcMSVC:
				return "arm-windows-msvc"
			}
		case OSAndroid:
			return "arm-linux-androideabi"
		}
	default:
		if llvm_arch, ok := llvm_archmap[e.Arch]; ok {
			return llvm_arch + "-" + string(e.OS) + "-" + string(e.Libc)
		}
	}

	return string(e.Arch) + "-" + string(e.OS) + "-" + string(e.Libc)
}

func Check() (env Environment, err error) {
	env.Libc = LibcNone
	env.Softfloat = false

	env.Arch = Arch(runtime.GOARCH)
	env.OS = OS(runtime.GOOS)

	// Libc Check
	switch env.OS {
	case OSLinux:
		env.Libc = checkLibc()
	case OSWindows:
		env.Libc = LibcMSVC
	case OSFreeBSD, OSDarwin:
		env.Libc = LibcNone
	case OSAndroid:
		env.Libc = LibcUnknown
	}

	// Softfloat Check
	switch env.Arch {
	case ArchARM:
		env.Softfloat = !cpu.ARM.HasVFPv3 // VFPv3 is required for hardfloat (aka armhf)
	}

	return
}
