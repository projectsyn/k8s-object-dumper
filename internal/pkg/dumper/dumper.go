// Dumper provides means to dump a UnstructuredList returned by a dynamic client.
package dumper

import (
	"encoding/json"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Dumper is an interface for dumping a list of unstructured objects
type DumperFunc func(*unstructured.UnstructuredList) error

// DumpToWriter dumps the list of unstructured objects to the provided writer as JSON
func DumpToWriter(w io.Writer) DumperFunc {
	return func(l *unstructured.UnstructuredList) error {
		return json.NewEncoder(w).Encode(l)
	}
}
