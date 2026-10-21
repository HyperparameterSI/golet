// Author: HyperparameterSI

package main

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func tarData(t *testing.T, headers ...tar.Header) []byte {
	t.Helper()
	var b bytes.Buffer
	w := tar.NewWriter(&b)
	for _, h := range headers {
		if err := w.WriteHeader(&h); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestExtractRejectsEscapes(t *testing.T) {
	for _, header := range []tar.Header{
		{Name: "../outside", Typeflag: tar.TypeReg},
		{Name: "/absolute", Typeflag: tar.TypeReg},
		{Name: "escape", Typeflag: tar.TypeSymlink, Linkname: "../outside"},
		{Name: "escape", Typeflag: tar.TypeLink, Linkname: "../outside"},
		{Name: "device", Typeflag: tar.TypeChar},
	} {
		t.Run(header.Name, func(t *testing.T) {
			if err := extractTar(context.Background(), bytes.NewReader(tarData(t, header)), t.TempDir()); err == nil {
				t.Fatal("unsafe entry accepted")
			}
		})
	}
}

func TestExtractFrameworkLinks(t *testing.T) {
	dest := t.TempDir()
	archive := tarData(t, tar.Header{Name: "Framework/Versions/A/binary", Typeflag: tar.TypeReg, Mode: 0755}, tar.Header{Name: "Framework/Versions/Current", Typeflag: tar.TypeSymlink, Linkname: "A"}, tar.Header{Name: "Framework/binary", Typeflag: tar.TypeSymlink, Linkname: "Versions/Current/binary"})
	if err := extractTar(context.Background(), bytes.NewReader(archive), dest); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, "Framework/binary")); err != nil {
		t.Fatal(err)
	}
	copied := filepath.Join(t.TempDir(), "Framework")
	if err := copyTree(context.Background(), filepath.Join(dest, "Framework"), copied); err != nil {
		t.Fatal(err)
	}
	if target, err := os.Readlink(filepath.Join(copied, "binary")); err != nil || target != "Versions/Current/binary" {
		t.Fatalf("link=%q error=%v", target, err)
	}
}

func TestVerifyArchive(t *testing.T) {
	data := []byte("cef archive")
	sum := sha1.Sum(data)
	file := filepath.Join(t.TempDir(), "archive")
	if err := os.WriteFile(file, data, 0600); err != nil {
		t.Fatal(err)
	}
	entry := runtimeEntry{Size: int64(len(data)), SHA1: hex.EncodeToString(sum[:])}
	if err := verifyArchive(file, entry); err != nil {
		t.Fatal(err)
	}
	entry.Size--
	if err := verifyArchive(file, entry); err == nil {
		t.Fatal("wrong size accepted")
	}
	entry.Size++
	entry.SHA1 = "bad"
	if err := verifyArchive(file, entry); err == nil {
		t.Fatal("wrong digest accepted")
	}
}

func TestReplaceRuntimePreservesPreviousOnFailure(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "installed")
	if err := os.Mkdir(output, 0755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(output, "golet-chromium-runtime.json")
	if err := os.WriteFile(marker, []byte("previous"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := replaceRuntime(filepath.Join(dir, "missing"), output); err == nil {
		t.Fatal("missing source accepted")
	}
	data, err := os.ReadFile(marker)
	if err != nil || string(data) != "previous" {
		t.Fatalf("old installation lost: %q %v", data, err)
	}
}

func TestRefuseUnmanagedOutput(t *testing.T) {
	if err := validateOutput(t.TempDir()); err == nil {
		t.Fatal("unmanaged directory accepted")
	}
}

// This opt-in check exercises the real CEF archive's framework/link layout.
func TestExtractPinnedSDK(t *testing.T) {
	archive := os.Getenv("GOLET_CEF_SDK_ARCHIVE")
	if archive == "" {
		t.Skip("set GOLET_CEF_SDK_ARCHIVE to test the pinned SDK archive")
	}
	manifest, err := os.ReadFile("../../native/runtime.json")
	if err != nil {
		t.Fatal(err)
	}
	var entries map[string]runtimeEntry
	if err = json.Unmarshal(manifest, &entries); err != nil {
		t.Fatal(err)
	}
	entry, ok := entries[runtime.GOOS+"-"+runtime.GOARCH]
	if !ok {
		t.Skip("no SDK for this platform")
	}
	if err = verifyArchive(archive, entry); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(archive)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	directory := t.TempDir()
	if err = extractTar(context.Background(), bzip2.NewReader(file), directory); err != nil {
		t.Fatal(err)
	}
	matches, err := filepath.Glob(filepath.Join(directory, "*", "include", "cef_version.h"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("SDK header missing: %v", err)
	}
	header, err := os.ReadFile(matches[0])
	if err != nil || !strings.Contains(string(header), entry.Version) {
		t.Fatalf("version mismatch: %v", err)
	}
}

func TestProductBundleMetadata(t *testing.T) {
	root := t.TempDir()
	binary := filepath.Join(root, "host")
	if err := os.WriteFile(binary, []byte("fixture"), 0755); err != nil {
		t.Fatal(err)
	}
	app := filepath.Join(root, "app")
	if err := bundle(app, "Golet Browser", binary, "", false, options{appName: "feo & 测试", identifier: "dev.feo.browser"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(app, "Contents/Info.plist"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"<key>CFBundleDisplayName</key><string>feo &amp; 测试</string>", "<key>CFBundleExecutable</key><string>Golet Browser</string>", "<key>CFBundleIdentifier</key><string>dev.feo.browser</string>"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("missing %s", want)
		}
	}
}
