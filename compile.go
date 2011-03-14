package main

import (
	"os"
	"io"
	"bufio"
	"fmt"
	"path"
	"exec"
	"runtime"
	"log"
)

func compile(source io.ReadSeeker, targetname string) os.Error {
	O, err := getArchSym()
	if err != nil {
		return err
	}
	gc := O + "g"
	gl := O + "l"
	tempobj := path.Join(os.TempDir(), targetname+"."+O)

	_, err = source.Seek(0, 0)
	if err != nil {
		return err
	}
	bufsource := bufio.NewReader(source)
	var insource io.Reader = bufsource
	if line, err := bufsource.ReadString('\n'); err != nil && err != os.EOF ||
		len(line) < 2 || line[:2] != "#!" {
		_, err := source.Seek(0, 0)
		if err != nil {
			return err
		}
		insource = source
	}

	err = run(gc, []string{gc, "-o", tempobj, "/dev/stdin"}, insource)
	if err != nil {
		return err
	}
	err = run(gl, []string{gl, "-o", path.Join(storedir, targetname),
		tempobj},
		nil)
	return err
}

func getArchSym() (sym string, err os.Error) {
	archsyms := map[string]string{
		"arm":   "5",
		"amd64": "6",
		"i386":  "8",
	}
	if sym, ok := archsyms[runtime.GOARCH]; ok {
		return sym, nil
	}
	return "", os.ErrorString("Unrecognized architecture")
}

func run(name string, argv []string, stdin io.Reader) os.Error {
	executable, err := exec.LookPath(name)
	if err != nil {
		return err
	}
	var sin int
	if stdin == nil {
		sin = exec.PassThrough
	} else {
		sin = exec.Pipe
	}
	proc, err := exec.Run(executable, argv, os.Environ(), "", sin,
		exec.PassThrough, exec.PassThrough)
	if err != nil {
		return err
	}
	defer proc.Close()
	if stdin != nil {
		if _, err := io.Copy(proc.Stdin, stdin); err != nil {
			return err
		}
		if err := proc.Stdin.Close(); err != nil {
			return err
		}
	}
	wm, err := proc.Wait(0)
	for err == os.EINTR {
		wm, err = proc.Wait(0)
	}
	if err != nil {
		return err
	}
	if wm.Exited() {
		if es := wm.ExitStatus(); es != 0 {
			return os.ErrorString(fmt.Sprintf("%s returned status %d", name,
				es))
		}
		return nil
	}
	return os.ErrorString("Something wierd happend with wait()")
}
