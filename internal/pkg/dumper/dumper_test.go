package dumper_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/projectsyn/k8s-object-dumper/internal/pkg/dumper"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_DumpToWriter(t *testing.T) {
	var b bytes.Buffer

	subject := dumper.DumpToWriter(&b)

	require.NoError(t,
		subject(&unstructured.UnstructuredList{
			Object: map[string]interface{}{
				"kind": "List",
			},
			Items: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name":      "test-pod",
							"namespace": "test-ns",
						},
					},
				},
			},
		}),
	)

	var got unstructured.UnstructuredList
	require.NoError(t, json.NewDecoder(&b).Decode(&got))

	require.Len(t, got.Items, 1)
	require.Equal(t, "Pod", got.Items[0].GetKind())
}
