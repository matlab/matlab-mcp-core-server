// Copyright 2026 The MathWorks, Inc.

package unpack_test

import (
	"archive/zip"
	"bytes"
	"path/filepath"
	"testing"

	mocks "github.com/matlab/matlab-mcp-server/tests/mocks/testutils/mcpbundle/unpack"
	"github.com/matlab/matlab-mcp-server/tests/testutils/mcpbundle/unpack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUnpack_ExtractsFiles(t *testing.T) {
	archive := createZip(t, map[string]string{
		"bin/launch.sh": "#!/bin/bash\necho hi",
		"manifest.json": `{"version":"1.0"}`,
	})

	wc := mocks.NewMockWriteCloser(t)
	wc.EXPECT().Write(mock.Anything).RunAndReturn(func(p []byte) (int, error) { return len(p), nil })
	wc.EXPECT().Close().Return(nil)

	fs := mocks.NewMockFileSystem(t)
	fs.EXPECT().MkdirAll(mock.Anything, mock.Anything).Return(nil)
	fs.EXPECT().OpenFile(mock.Anything, mock.Anything, mock.Anything).Return(wc, nil)

	u := unpack.NewUnpackerForTest(fs)

	err := u.Unpack(bytes.NewReader(archive), int64(len(archive)), "/dest")

	require.NoError(t, err)
}

func TestUnpack_CreatesParentDirectories(t *testing.T) {
	archive := createZip(t, map[string]string{
		"bin/launch.sh": "content",
	})

	wc := mocks.NewMockWriteCloser(t)
	wc.EXPECT().Write(mock.Anything).RunAndReturn(func(p []byte) (int, error) { return len(p), nil })
	wc.EXPECT().Close().Return(nil)

	fs := mocks.NewMockFileSystem(t)
	fs.EXPECT().MkdirAll(filepath.Join("/dest", "bin"), mock.Anything).Return(nil)
	fs.EXPECT().OpenFile(mock.Anything, mock.Anything, mock.Anything).Return(wc, nil)

	u := unpack.NewUnpackerForTest(fs)

	err := u.Unpack(bytes.NewReader(archive), int64(len(archive)), "/dest")

	require.NoError(t, err)
}

func TestUnpack_InvalidArchive(t *testing.T) {
	data := []byte("not a zip")

	fs := mocks.NewMockFileSystem(t)
	u := unpack.NewUnpackerForTest(fs)

	err := u.Unpack(bytes.NewReader(data), int64(len(data)), "/dest")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "reading archive")
}

func TestUnpack_MkdirError(t *testing.T) {
	archive := createZip(t, map[string]string{
		"bin/launch.sh": "content",
	})

	fs := mocks.NewMockFileSystem(t)
	fs.EXPECT().MkdirAll(mock.Anything, mock.Anything).Return(assert.AnError)

	u := unpack.NewUnpackerForTest(fs)

	err := u.Unpack(bytes.NewReader(archive), int64(len(archive)), "/dest")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating parent")
}

func TestUnpack_OpenFileError(t *testing.T) {
	archive := createZip(t, map[string]string{
		"file.txt": "content",
	})

	fs := mocks.NewMockFileSystem(t)
	fs.EXPECT().MkdirAll(mock.Anything, mock.Anything).Return(nil)
	fs.EXPECT().OpenFile(mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

	u := unpack.NewUnpackerForTest(fs)

	err := u.Unpack(bytes.NewReader(archive), int64(len(archive)), "/dest")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating file.txt")
}

func TestUnpack_RootDirectoryEntry(t *testing.T) {
	archive := createZipWithDirs(t,
		map[string]string{"file.txt": "content"},
		[]string{"./"},
	)

	wc := mocks.NewMockWriteCloser(t)
	wc.EXPECT().Write(mock.Anything).RunAndReturn(func(p []byte) (int, error) { return len(p), nil })
	wc.EXPECT().Close().Return(nil)

	fs := mocks.NewMockFileSystem(t)
	fs.EXPECT().MkdirAll(mock.Anything, mock.Anything).Return(nil)
	fs.EXPECT().OpenFile(mock.Anything, mock.Anything, mock.Anything).Return(wc, nil)

	u := unpack.NewUnpackerForTest(fs)

	err := u.Unpack(bytes.NewReader(archive), int64(len(archive)), "/dest")

	require.NoError(t, err)
}

func TestUnpack_ZipSlip_RejectsPathTraversal(t *testing.T) {
	archive := createZip(t, map[string]string{
		"../evil.txt": "malicious content",
	})

	fs := mocks.NewMockFileSystem(t)
	u := unpack.NewUnpackerForTest(fs)

	err := u.Unpack(bytes.NewReader(archive), int64(len(archive)), "/dest")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "illegal file path in archive")
}

func TestUnpack_ZipSlip_RejectsAbsolutePath(t *testing.T) {
	archive := createZip(t, map[string]string{
		"/etc/evil.txt": "malicious content",
	})

	fs := mocks.NewMockFileSystem(t)
	u := unpack.NewUnpackerForTest(fs)

	err := u.Unpack(bytes.NewReader(archive), int64(len(archive)), "/dest")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "illegal file path in archive")
}

func createZip(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	return createZipWithDirs(t, entries, nil)
}

func createZipWithDirs(t *testing.T, entries map[string]string, dirs []string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, dir := range dirs {
		_, err := w.Create(dir)
		require.NoError(t, err)
	}
	for name, content := range entries {
		f, err := w.Create(name)
		require.NoError(t, err)
		_, err = f.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, w.Close())
	return buf.Bytes()
}
