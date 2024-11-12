package dumper_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/projectsyn/k8s-object-dumper/internal/pkg/dumper"
)

func Test_DirDumper(t *testing.T) {
	tdir, err := os.MkdirTemp(".", "test")
	require.NoError(t, err)
	defer os.RemoveAll(tdir)

	subject, err := dumper.NewDirDumper(tdir)
	require.NoError(t, err)

	uls := []*unstructured.UnstructuredList{
		{
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
				{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name":      "test-pod",
							"namespace": "test-ns-2",
						},
					},
				},
			},
		},
		{
			Items: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name":      "test-pod-2",
							"namespace": "test-ns",
						},
					},
				},
			},
		},
		{
			Items: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"kind":       "Service",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name":      "test-svc",
							"namespace": "test-ns",
						},
					},
				},
			},
		},
		{
			Items: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"kind":       "ClusterRole",
						"apiVersion": "rbac.authorization.k8s.io/v1",
						"metadata": map[string]interface{}{
							"name": "cluster-scoped",
						},
					},
				},
			},
		},
	}

	for i, ul := range uls {
		require.NoErrorf(t, subject.Dump(ul), "failed to dump list %d", i)
	}
	defer func() {
		require.NoError(t, subject.Close())
	}()

	require.FileExists(t, tdir+"/objects-Pod.json")
	requireFileContains(t, tdir+"/objects-Pod.json", []ExpectedObject{
		{Kind: "Pod", Name: "test-pod", Namespace: "test-ns"},
		{Kind: "Pod", Name: "test-pod", Namespace: "test-ns-2"},
		{Kind: "Pod", Name: "test-pod-2", Namespace: "test-ns"},
	})
	require.FileExists(t, tdir+"/objects-Service.json")
	requireFileContains(t, tdir+"/objects-Service.json", []ExpectedObject{
		{Kind: "Service", Name: "test-svc", Namespace: "test-ns"},
	})
	require.FileExists(t, tdir+"/objects-ClusterRole.rbac.authorization.k8s.io.json")
	requireFileContains(t, tdir+"/objects-ClusterRole.rbac.authorization.k8s.io.json", []ExpectedObject{
		{Kind: "ClusterRole", Name: "cluster-scoped", Namespace: ""},
	})
	require.FileExists(t, tdir+"/split/test-ns/__all__.json")
	requireFileContains(t, tdir+"/split/test-ns/__all__.json", []ExpectedObject{
		{Kind: "Pod", Name: "test-pod", Namespace: "test-ns"},
		{Kind: "Pod", Name: "test-pod-2", Namespace: "test-ns"},
		{Kind: "Service", Name: "test-svc", Namespace: "test-ns"},
	})
	require.FileExists(t, tdir+"/split/test-ns/Pod.json")
	requireFileContains(t, tdir+"/split/test-ns/Pod.json", []ExpectedObject{
		{Kind: "Pod", Name: "test-pod", Namespace: "test-ns"},
		{Kind: "Pod", Name: "test-pod-2", Namespace: "test-ns"},
	})
	require.FileExists(t, tdir+"/split/test-ns/Service.json")
	requireFileContains(t, tdir+"/split/test-ns/Service.json", []ExpectedObject{
		{Kind: "Service", Name: "test-svc", Namespace: "test-ns"},
	})
	require.FileExists(t, tdir+"/split/test-ns-2/__all__.json")
	requireFileContains(t, tdir+"/split/test-ns-2/__all__.json", []ExpectedObject{
		{Kind: "Pod", Name: "test-pod", Namespace: "test-ns-2"},
	})
	require.FileExists(t, tdir+"/split/test-ns-2/Pod.json")
	requireFileContains(t, tdir+"/split/test-ns-2/Pod.json", []ExpectedObject{
		{Kind: "Pod", Name: "test-pod", Namespace: "test-ns-2"},
	})
}

type ExpectedObject struct {
	Kind, Name, Namespace string
}

func requireFileContains(t *testing.T, path string, expected []ExpectedObject) {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)

	var objs []unstructured.Unstructured
	raw = bytes.TrimSuffix(raw, []byte("\n"))
	rawObjs := bytes.Split(raw, []byte("\n"))
	for _, rawObj := range rawObjs {
		var obj unstructured.Unstructured
		require.NoError(t, json.Unmarshal(rawObj, &obj))
		objs = append(objs, obj)
	}

	actualObjects := make([]ExpectedObject, 0, len(objs))
	for _, obj := range objs {
		actualObjects = append(actualObjects, ExpectedObject{
			Kind:      obj.GetKind(),
			Name:      obj.GetName(),
			Namespace: obj.GetNamespace(),
		})
	}

	require.ElementsMatch(t, expected, actualObjects)
}
