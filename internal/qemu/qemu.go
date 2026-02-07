package qemu

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Common QEMU image formats you might use in a libvirt lab.
const (
	FormatQCOW2 = "qcow2"
	FormatRAW   = "raw"
)

var (
	// Examples: "200G", "40GiB", "10240M", "50g", "8192K"
	sizeRe = regexp.MustCompile(`(?i)^\s*\+?\d+\s*(k|kb|kib|m|mb|mib|g|gb|gib|t|tb|tib)?\s*$`)
)

// CreateOptions controls how qemu-img create is invoked.
type CreateOptions struct {
	// Format defaults to qcow2 if empty.
	Format string

	// Size is required unless BackingFile is set AND you intend to inherit virtual size.
	// In practice, for predictable automation, you should always set Size.
	Size string

	// Optional qcow2 backing store.
	BackingFile string

	// Optional backing format (recommended if BackingFile is set).
	BackingFormat string

	// Additional qemu-img options string passed via -o, e.g.:
	// "cluster_size=2M,preallocation=metadata,lazy_refcounts=on"
	Options string
}

// Create creates a disk image at outPath.
// For qcow2 backing files, it will create a new qcow2 referencing BackingFile.
func Create(outPath string, opts CreateOptions) error {
	if strings.TrimSpace(outPath) == "" {
		return errors.New("outPath is required")
	}
	if err := ensureParentDir(outPath); err != nil {
		return err
	}

	format := strings.TrimSpace(opts.Format)
	if format == "" {
		format = FormatQCOW2
	}

	// Validate size if provided.
	if s := strings.TrimSpace(opts.Size); s != "" {
		if !isValidSize(s) {
			return fmt.Errorf("invalid size %q (examples: 40G, 200GiB, 10240M)", opts.Size)
		}
	} else {
		// Size omitted: allow only if using a backing file and caller accepts inherited size behavior.
		// For predictability in Anvil, you may choose to make Size mandatory.
		if strings.TrimSpace(opts.BackingFile) == "" {
			return errors.New("size is required when no backing file is specified")
		}
	}

	args := []string{"create", "-f", format}

	// Backing file handling (primarily for qcow2).
	if bf := strings.TrimSpace(opts.BackingFile); bf != "" {
		args = append(args, "-b", bf)
		if bff := strings.TrimSpace(opts.BackingFormat); bff != "" {
			args = append(args, "-F", bff)
		}
	}

	// -o options (format-specific)
	if o := strings.TrimSpace(opts.Options); o != "" {
		args = append(args, "-o", o)
	}

	args = append(args, outPath)

	// Size is last argument when provided.
	if s := strings.TrimSpace(opts.Size); s != "" {
		args = append(args, s)
	}

	if _, err := runQemuImg(args...); err != nil {
		return err
	}
	return nil
}

// Resize grows a disk image to newSize. By default, shrinking is blocked.
// newSize can be absolute (e.g. "200G") or relative growth (e.g. "+20G").
// Set allowShrink=true to permit shrinking (generally not recommended).
func Resize(imagePath string, newSize string, allowShrink bool) error {
	if strings.TrimSpace(imagePath) == "" {
		return errors.New("imagePath is required")
	}
	if strings.TrimSpace(newSize) == "" {
		return errors.New("newSize is required")
	}
	if !isValidSize(newSize) {
		return fmt.Errorf("invalid newSize %q (examples: 200G, +20G, 40GiB)", newSize)
	}
	if _, err := os.Stat(imagePath); err != nil {
		return fmt.Errorf("image does not exist: %s", imagePath)
	}

	// If it's not a relative increase (+X), enforce no shrink unless allowShrink.
	isRelativeGrow := strings.HasPrefix(strings.TrimSpace(newSize), "+")
	if !allowShrink && !isRelativeGrow {
		curBytes, err := VirtualSizeBytes(imagePath)
		if err != nil {
			return err
		}
		targetBytes, err := parseSizeToBytes(newSize)
		if err != nil {
			return err
		}
		if targetBytes < curBytes {
			return fmt.Errorf("refusing to shrink image: current=%d bytes target=%d bytes (pass allowShrink=true to override)", curBytes, targetBytes)
		}
	}

	_, err := runQemuImg("resize", imagePath, newSize)
	return err
}

// VirtualSizeBytes returns the virtual size in bytes using `qemu-img info --output=json`.
func VirtualSizeBytes(imagePath string) (int64, error) {
	out, err := runQemuImg("info", "--output=json", imagePath)
	if err != nil {
		return 0, err
	}

	// Minimal JSON parsing without extra deps:
	// We look for: "virtual-size": 12345
	// qemu-img emits this key in modern builds.
	key := `"virtual-size"`
	idx := strings.Index(out, key)
	if idx == -1 {
		return 0, fmt.Errorf("unable to find %s in qemu-img info output", key)
	}
	rest := out[idx+len(key):]
	colon := strings.Index(rest, ":")
	if colon == -1 {
		return 0, fmt.Errorf("malformed qemu-img info output near %s", key)
	}
	rest = rest[colon+1:]
	// take until comma or newline or }
	end := len(rest)
	for i, r := range rest {
		if r == ',' || r == '\n' || r == '}' {
			end = i
			break
		}
	}
	numStr := strings.TrimSpace(rest[:end])
	v, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse virtual-size %q: %w", numStr, err)
	}
	return v, nil
}

func runQemuImg(args ...string) (string, error) {
	cmd := exec.Command("qemu-img", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("qemu-img %s failed: %s", strings.Join(args, " "), msg)
	}

	return strings.TrimSpace(out.String()), nil
}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "/" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}

func isValidSize(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// Allow relative sizes like +20G
	if strings.HasPrefix(s, "+") {
		s = strings.TrimSpace(strings.TrimPrefix(s, "+"))
		s = "+" + s // keep format, regex allows optional + but we re-validate below
	}
	return sizeRe.MatchString(s)
}

// parseSizeToBytes converts human sizes like "200G", "40GiB", "10240M" to bytes.
// NOTE: This is used only for shrink-prevention comparisons.
// It uses binary units for KiB/MiB/GiB/TiB and decimal for KB/MB/GB/TB.
func parseSizeToBytes(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if strings.HasPrefix(s, "+") {
		return 0, fmt.Errorf("relative size %q cannot be parsed as absolute bytes", s)
	}

	// split number and suffix
	re := regexp.MustCompile(`^\s*(\d+)\s*([a-z]+)?\s*$`)
	m := re.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("invalid size %q", s)
	}

	n, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric size %q: %w", m[1], err)
	}

	suf := m[2]
	mult := int64(1)

	switch suf {
	case "", "b":
		mult = 1
	case "k", "kb":
		mult = 1000
	case "kib":
		mult = 1024
	case "m", "mb":
		mult = 1000 * 1000
	case "mib":
		mult = 1024 * 1024
	case "g", "gb":
		mult = 1000 * 1000 * 1000
	case "gib":
		mult = 1024 * 1024 * 1024
	case "t", "tb":
		mult = 1000 * 1000 * 1000 * 1000
	case "tib":
		mult = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unknown size suffix %q", suf)
	}

	// check overflow-ish
	if n > (1<<62)/mult {
		return 0, fmt.Errorf("size too large: %q", s)
	}
	return n * mult, nil
}
