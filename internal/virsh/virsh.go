package virsh

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func run(args ...string) (string, error) {
	cmd := exec.Command("virsh", args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s", stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func Start(vm string) error {
	_, err := run("start", vm)
	return err
}

func Stop(vm string) error {
	_, err := run("shutdown", vm)
	return err
}

func List(all bool) (string, error) {
	if all {
		return run("list", "--all")
	}
	return run("list")
}

func Clone(src, dst string) error {
	_, err := run(
		"clone",
		src,
		dst,
		"--auto-clone",
	)
	return err
}
