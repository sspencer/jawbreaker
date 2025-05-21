package main

import (
	"bytes"
	"io"
	"io/fs"
	"sync"
	"time"
)

// ConcatFS implements fs.FS to present all files in a directory as a single concatenated file.
type ConcatFS struct {
	dirFS fs.FS
	name  string // Name of the virtual concatenated file
	once  sync.Once
	data  []byte
	err   error
}

// NewConcatFS creates a new ConcatFS that concatenates all files from the given directory
// into a single virtual file with the specified name.
func NewConcatFS(dirFS fs.FS, name string) *ConcatFS {
	return &ConcatFS{
		dirFS: dirFS,
		name:  name,
	}
}

// Open implements fs.FS, returning a file containing the concatenated content of all files
// if the requested name matches the virtual file name, or an error otherwise.
func (c *ConcatFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	if name != c.name {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	// Load content lazily
	c.once.Do(c.loadContent)

	if c.err != nil {
		return nil, c.err
	}

	return &concatFile{
		name:   c.name,
		reader: bytes.NewReader(c.data),
	}, nil
}

// loadContent walks the directory and concatenates all files.
func (c *ConcatFS) loadContent() {
	var buf bytes.Buffer

	err := fs.WalkDir(c.dirFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := c.dirFS.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(&buf, f)
		if err != nil {
			return err
		}

		// Add a newline between files for clarity
		buf.WriteByte('\n')
		return nil
	})

	if err != nil {
		c.err = err
		return
	}

	c.data = buf.Bytes()
}

// concatFile implements fs.File for the concatenated content.
type concatFile struct {
	name   string
	reader *bytes.Reader
}

// Stat returns the FileInfo for the concatenated file.
func (f *concatFile) Stat() (fs.FileInfo, error) {
	return &concatFileInfo{
		name: f.name,
		size: int64(f.reader.Len()),
	}, nil
}

// Read reads from the concatenated content.
func (f *concatFile) Read(b []byte) (int, error) {
	return f.reader.Read(b)
}

// Close is a no-op since we're reading from memory.
func (f *concatFile) Close() error {
	return nil
}

// concatFileInfo implements fs.FileInfo for the concatenated file.
type concatFileInfo struct {
	name string
	size int64
}

// Name returns the name of the concatenated file.
func (i *concatFileInfo) Name() string { return i.name }

// Size returns the total size of the concatenated content.
func (i *concatFileInfo) Size() int64 { return i.size }

// Mode returns the file mode (regular file).
func (i *concatFileInfo) Mode() fs.FileMode { return 0444 }

// ModTime returns a zero time since this is a virtual file.
func (i *concatFileInfo) ModTime() time.Time { return time.Time{} }

// IsDir returns false since this is a file.
func (i *concatFileInfo) IsDir() bool { return false }

// Sys returns nil as there's no underlying system info.
func (i *concatFileInfo) Sys() interface{} { return nil }

// Example usage
// 	//go:embed static/js/*.js
//	jawbreakerFS embed.FS
//
// func main() {
//	// Create a ConcatFS that presents all *.js files in a single directory
//  // as a single file
//	concatFS := NewConcatFS(jawbreakerFS, "all.js")
//
//	// Open the concatenated file
//	f, err := concatFS.Open("all.js")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	// Read the concatenated content
//	content, err := io.ReadAll(f)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Concatenated content:\n%s\n", content)
// }
