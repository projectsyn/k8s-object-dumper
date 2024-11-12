package dumper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// DirDumper writes objects to a directory.
// Must be initialized with newDirDumper.
// Must be closed after use.
type DirDumper struct {
	dir string

	openFiles map[string]*os.File
	sharedBuf *bytes.Buffer
}

// NewDirDumper creates a new dirDumper that writes objects to the given directory.
// The directory will be created if it does not exist.
// If the directory cannot be created, an error is returned.
func NewDirDumper(dir string) (*DirDumper, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %q: %w", dir, err)
	}
	return &DirDumper{
		dir:       dir,
		openFiles: make(map[string]*os.File),
		sharedBuf: new(bytes.Buffer),
	}, nil
}

// Close closes the dirDumper and all open files.
// The dirDumper cannot be used after it is closed.
func (d *DirDumper) Close() error {
	var errs []error
	for _, f := range d.openFiles {
		if err := f.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return multierr.Combine(errs...)
}

// Dump writes the objects in the list to the directory.
// The objects are written to the directory in two ways:
// - All objects are written to a file named objects-<kind>.json
// - Objects with a namespace are written to a directory named split/<namespace> with two files:
//   - __all__.json contains all objects in the namespace
//   - <kind>.json contains all objects of the kind in the namespace
//
// If an object cannot be written, an error is returned.
// This method is not safe for concurrent use.
func (d *DirDumper) Dump(l *unstructured.UnstructuredList) error {
	buf := d.sharedBuf
	var errs []error
	for _, o := range l.Items {
		buf.Reset()
		if err := json.NewEncoder(buf).Encode(o.Object); err != nil {
			errs = append(errs, fmt.Errorf("failed to encode object: %w", err))
			continue
		}
		p := buf.Bytes()
		gk := o.GroupVersionKind().GroupKind()

		if err := d.writeToFile(fmt.Sprintf("%s/objects-%s.json", d.dir, gk), p); err != nil {
			errs = append(errs, err)
		}

		if o.GetNamespace() == "" {
			continue
		}

		if err := d.writeToFile(fmt.Sprintf("%s/split/%s/__all__.json", d.dir, o.GetNamespace()), p); err != nil {
			errs = append(errs, err)
		}
		if err := d.writeToFile(fmt.Sprintf("%s/split/%s/%s.json", d.dir, o.GetNamespace(), gk), p); err != nil {
			errs = append(errs, err)
		}
	}
	return multierr.Combine(errs...)
}

func (d *DirDumper) writeToFile(path string, b []byte) error {
	f, err := d.file(path)
	if err != nil {
		return fmt.Errorf("failed to open file for copying: %w", err)
	}
	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("failed to copy to file: %w", err)
	}
	return nil
}

func (d *DirDumper) file(path string) (*os.File, error) {
	f, ok := d.openFiles[path]
	if ok {
		return f, nil
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %q: %w", dir, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %q: %w", path, err)
	}
	d.openFiles[path] = f
	return f, nil
}
