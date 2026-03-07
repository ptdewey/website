package generator

import (
	"io"
	"os"
	"path/filepath"
)

func CopyDirIncremental(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		dstInfo, err := os.Stat(target)
		if err == nil {
			srcMod := info.ModTime().UnixNano()
			dstMod := dstInfo.ModTime().UnixNano()
			if dstMod >= srcMod {
				return nil
			}
		}

		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp := dst + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmp)
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.Chmod(tmp, srcInfo.Mode()); err != nil {
		return err
	}

	return os.Rename(tmp, dst)
}
