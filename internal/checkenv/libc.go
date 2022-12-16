package checkenv

import (
	"bytes"
	"os/exec"
	"runtime"
)

func checkLibc() Libc {
	switch runtime.GOOS {
	case "linux":
		getconf, err := exec.LookPath("getconf")
		if err == nil {
			cmd := exec.Command(getconf, "GNU_LIBC_VERSION")
			var errout bytes.Buffer
			cmd.Stderr = &errout
			err = cmd.Run()
			if err == nil {
				return LibcGNU
			}
		}

		ldd, err := exec.LookPath("ldd")
		if err == nil {
			out, _ := exec.Command(ldd, "--version").CombinedOutput()
			data := bytes.ToLower(out)
			if bytes.Contains(data, []byte("musl")) {
				// High probability that we are using musl
				return LibcMUSL
			}
			if bytes.Contains(data, []byte("glibc")) {
				// High probability that we are using glibc
				return LibcGNU
			}
		}

		// Fallback to unknown
		return LibcUnknown
	case "windows":
		// TODO: check libc check for windows
		return LibcMSVC
	case "freebsd", "darwin":
		return LibcNone
	}
	return LibcUnknown
}
