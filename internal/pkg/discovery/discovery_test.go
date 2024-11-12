package discovery_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/projectsyn/k8s-object-dumper/internal/pkg/discovery"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type objKey struct {
	apiVersion, kind, name, namespace string
}

func Test_DiscoverObjects(t *testing.T) {
	cfg, stop := setupEnvtestEnv(t)
	defer stop()

	objs := map[objKey]unstructured.Unstructured{}
	objTracker := func(obj *unstructured.UnstructuredList) error {
		for _, o := range obj.Items {
			objs[objKey{apiVersion: o.GetAPIVersion(), kind: o.GetKind(), name: o.GetName(), namespace: o.GetNamespace()}] = o
		}
		return nil
	}

	c, err := client.New(cfg, client.Options{})
	require.NoError(t, err)

	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-ns",
		},
	}
	sas := make([]client.Object, 0, 10)
	for i := 1; i <= cap(sas); i++ {
		sas = append(sas, &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("test-service-account-%d", i),
				Namespace: "test-ns",
			},
		})
	}
	cr := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster-role",
		},
	}
	r := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-role",
			Namespace: "test-ns",
		},
	}

	for _, obj := range append([]client.Object{&ns, &cr, &r}, sas...) {
		require.NoError(t, c.Create(context.Background(), obj))
	}

	require.NoError(t, discovery.DiscoverObjects(context.Background(), cfg, objTracker, discovery.DiscoveryOptions{
		ChunkSize: int64(cap(sas) / 2),
		IgnoreResources: []*regexp.Regexp{
			regexp.MustCompile(`^roles.rbac.authorization.k8s.io$`),
		},
		MustExistResources: []string{
			"clusterroles.rbac.authorization.k8s.io",
			"deployments.apps",
			"namespaces",
		},
	}))

	require.Contains(t, objs, objKey{apiVersion: "v1", kind: "Namespace", name: "test-ns", namespace: ""})
	require.Contains(t, objs, objKey{apiVersion: "rbac.authorization.k8s.io/v1", kind: "ClusterRole", name: "test-cluster-role", namespace: ""})
	for i := 1; i <= cap(sas); i++ {
		require.Contains(t, objs, objKey{apiVersion: "v1", kind: "ServiceAccount", name: fmt.Sprintf("test-service-account-%d", i), namespace: "test-ns"})
	}
	require.NotContains(t, objs, objKey{apiVersion: "rbac.authorization.k8s.io/v1", kind: "Role", name: "test-role", namespace: "test-ns"}, "Roles are ignored by regex")
}

func Test_DiscoverObjects_MustExistResources_NotSatisfied(t *testing.T) {
	cfg, stop := setupEnvtestEnv(t)
	defer stop()

	discard := func(obj *unstructured.UnstructuredList) error {
		return nil
	}

	require.ErrorContains(t, discovery.DiscoverObjects(context.Background(), cfg, discard, discovery.DiscoveryOptions{
		MustExistResources: []string{
			"fluxcapacitors.spaceship.io",
			"namespaces",
		},
	}), "missing resources: [fluxcapacitors.spaceship.io]")
}

func setupEnvtestEnv(t *testing.T) (cfg *rest.Config, stop func()) {
	t.Helper()

	testEnv := &envtest.Environment{}

	cfg, err := testEnv.Start()
	require.NoError(t, err)

	return cfg, func() {
		require.NoError(t, testEnv.Stop())
	}
}
