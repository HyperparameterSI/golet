// Author: HyperparameterSI

package main

import (
	"archive/tar"
	"compress/bzip2"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func verifyArchive(file string, entry runtimeEntry) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha1.New() // CEF publishes SHA-1 in its distribution index; use the pinned digest.
	size, err := io.Copy(h, io.LimitReader(f, entry.Size+1))
	if err != nil {
		return err
	}
	if size != entry.Size || hex.EncodeToString(h.Sum(nil)) != entry.SHA1 {
		return fmt.Errorf("CEF archive size/checksum mismatch: %s", file)
	}
	return nil
}

func fetchSDK(ctx context.Context, entry runtimeEntry, cache string) (string, error) {
	u, err := url.Parse(entry.URL)
	if err != nil {
		return "", err
	}
	if u.Scheme != "https" || entry.Size <= 0 {
		return "", errors.New("invalid pinned CEF download")
	}
	name := path.Base(u.Path)
	if !strings.HasSuffix(name, ".tar.bz2") {
		return "", errors.New("expected a .tar.bz2 SDK")
	}
	archive := filepath.Join(cache, name)
	if _, err = os.Stat(archive); os.IsNotExist(err) {
		downloadCtx, cancel := context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()
		req, err := http.NewRequestWithContext(downloadCtx, http.MethodGet, entry.URL, nil)
		if err != nil {
			return "", err
		}
		client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Scheme != "https" || len(via) >= 10 {
				return errors.New("invalid CEF redirect")
			}
			return nil
		}}
		fmt.Printf("Downloading CEF (%.0f MB)...\n", float64(entry.Size)/1e6)
		res, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return "", fmt.Errorf("CEF download: %s", res.Status)
		}
		temp, err := os.CreateTemp(cache, ".cef-download-")
		if err != nil {
			return "", err
		}
		defer os.Remove(temp.Name())
		_, copyErr := io.Copy(temp, io.LimitReader(res.Body, entry.Size+1))
		closeErr := temp.Close()
		if err = errors.Join(copyErr, closeErr); err != nil {
			return "", err
		}
		if err = verifyArchive(temp.Name(), entry); err != nil {
			return "", err
		}
		if err = os.Rename(temp.Name(), archive); err != nil {
			return "", err
		}
	}
	if err = verifyArchive(archive, entry); err != nil {
		return "", err
	}
	directory := filepath.Join(cache, strings.TrimSuffix(name, ".tar.bz2"))
	if info, err := os.Stat(directory); err == nil && info.IsDir() {
		return directory, nil
	}
	temp, err := os.MkdirTemp(cache, ".cef-extract-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(temp)
	file, err := os.Open(archive)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err = extractTar(ctx, bzip2.NewReader(file), temp); err != nil {
		return "", err
	}
	if err = os.Rename(filepath.Join(temp, filepath.Base(directory)), directory); err != nil {
		return "", err
	}
	return directory, nil
}

func safeArchiveName(name string) (string, error) {
	if strings.Contains(name, "\\") || path.IsAbs(name) {
		return "", fmt.Errorf("unsafe archive path %q", name)
	}
	name = path.Clean(name)
	if name == ".." || strings.HasPrefix(name, "../") {
		return "", fmt.Errorf("unsafe archive path %q", name)
	}
	return filepath.FromSlash(name), nil
}

// os.Root confines every filesystem operation, including traversal through
// earlier symlinks. Links are deferred until regular contents are extracted.
func extractTar(ctx context.Context, reader io.Reader, destination string) error {
	root, err := os.OpenRoot(destination)
	if err != nil {
		return err
	}
	defer root.Close()
	tr := tar.NewReader(reader)
	var links []*tar.Header
	var total int64
	for count := 0; ; count++ {
		if err = ctx.Err(); err != nil {
			return err
		}
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if count >= 100000 {
			return errors.New("too many CEF archive entries")
		}
		name, err := safeArchiveName(h.Name)
		if err != nil {
			return err
		}
		total += h.Size
		if h.Size < 0 || total > 8<<30 {
			return errors.New("CEF archive exceeds 8 GiB extraction limit")
		}
		switch h.Typeflag {
		case tar.TypeDir:
			err = root.MkdirAll(name, 0755)
		case tar.TypeReg, tar.TypeRegA:
			if err = root.MkdirAll(filepath.Dir(name), 0755); err != nil {
				return err
			}
			var file *os.File
			file, err = root.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(h.Mode)&0777)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(file, tr)
			err = errors.Join(copyErr, file.Close())
		case tar.TypeSymlink, tar.TypeLink:
			copyHeader := *h
			links = append(links, &copyHeader)
		default:
			return fmt.Errorf("unsupported archive entry %q (type %d)", h.Name, h.Typeflag)
		}
		if err != nil {
			return err
		}
	}
	for _, h := range links {
		name, err := safeArchiveName(h.Name)
		if err != nil {
			return err
		}
		if err = root.MkdirAll(filepath.Dir(name), 0755); err != nil {
			return err
		}
		if h.Typeflag == tar.TypeSymlink {
			if path.IsAbs(h.Linkname) || strings.Contains(h.Linkname, "\\") {
				return errors.New("unsafe archive symlink")
			}
			if _, err = safeArchiveName(path.Join(path.Dir(h.Name), h.Linkname)); err != nil {
				return err
			}
			err = root.Symlink(h.Linkname, name)
		} else {
			target, e := safeArchiveName(h.Linkname)
			if e != nil {
				return e
			}
			err = root.Link(target, name)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
