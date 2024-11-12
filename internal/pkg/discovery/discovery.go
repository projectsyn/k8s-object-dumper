package discovery

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"go.uber.org/multierr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type DiscoveryOptions struct {
	BatchSize int64
	LogWriter io.Writer

	// MustExistResources is a list of resources that must exist in the cluster.
	// This can be used as a sanity check to ensure that the discovery process is working as expected.
	// If a resource does not exist, the discovery process will fail.
	// If the list is empty, no resources are required to exist.
	MustExistResources []string

	// IgnoreResources is a list of resources to ignore during discovery.
	IgnoreResources []*regexp.Regexp
}

// GetBatchSize returns the set batch size for listing objects or the default.
func (opts DiscoveryOptions) GetBatchSize() int64 {
	if opts.BatchSize == 0 {
		return 500
	}
	return opts.BatchSize
}

// GetLogWriter returns the set batch size for listing objects or io.Discard as default.
func (opts DiscoveryOptions) GetLogWriter() io.Writer {
	if opts.LogWriter == nil {
		return io.Discard
	}
	return opts.LogWriter
}

// DiscoverObjects discovers all objects in the cluster and calls the provided callback for each list of objects.
// The callback can be called multiple times with the same object kind.
// Objects are unique in general, but the callback should be able to handle duplicates.
// Some API servers do not implement list batching correctly and thus might introduce duplicates.
func DiscoverObjects(ctx context.Context, conf *rest.Config, cb func(*unstructured.UnstructuredList) error, opts DiscoveryOptions) error {
	batchSize := opts.GetBatchSize()
	logWriter := opts.GetLogWriter()

	dc, err := discovery.NewDiscoveryClientForConfig(conf)
	if err != nil {
		return fmt.Errorf("failed to create discovery client: %w", err)
	}
	dynClient, err := dynamic.NewForConfig(conf)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	sprl, err := dc.ServerPreferredResources()
	if err != nil {
		return fmt.Errorf("failed to get server preferred resources: %w", err)
	}

	fmt.Fprintln(logWriter, "Discovered resources:")
	for _, re := range sprl {
		fmt.Fprintln(logWriter, re.GroupVersion)
		for _, r := range re.APIResources {
			fmt.Fprintln(logWriter, "  ", r.Kind)
		}
	}

	if len(opts.MustExistResources) > 0 {
		want := sets.New(opts.MustExistResources...)
		have := sets.New[string]()
		for _, re := range sprl {
			for _, r := range re.APIResources {
				res := formatGVRForComparison(groupVersionFromString(re.GroupVersion).WithResource(r.Name))
				have.Insert(res)
			}
		}
		missing := want.Difference(have)
		if missing.Len() > 0 {
			return fmt.Errorf("missing resources: %s", sets.List(missing))
		}
	}

	var errors []error
	for _, re := range sprl {
		for _, r := range re.APIResources {
			res := groupVersionFromString(re.GroupVersion).WithResource(r.Name)
			if !slices.Contains(r.Verbs, "list") {
				fmt.Fprintf(logWriter, "skipping %s: no list verb\n", res)
				continue
			}

			if i := slices.IndexFunc(opts.IgnoreResources, func(re *regexp.Regexp) bool {
				return re.MatchString(formatGVRForComparison(res))
			}); i > -1 {
				fmt.Fprintf(logWriter, "skipping %s: ignored by regex %q\n", res, opts.IgnoreResources[i].String())
				continue
			}

			continueKey := ""
			for {
				l, err := dynClient.Resource(res).List(ctx, metav1.ListOptions{
					Limit:    batchSize,
					Continue: continueKey,
				})
				if err != nil {
					errors = append(errors, fmt.Errorf("failed to list %s: %w", res, err))
					break
				}
				if err := cb(l); err != nil {
					errors = append(errors, fmt.Errorf("failed to dump %s: %w", res, err))
				}
				if l.GetContinue() == "" {
					break
				}
				continueKey = l.GetContinue()
			}
		}
	}

	return multierr.Combine(errors...)
}

func groupVersionFromString(s string) schema.GroupVersion {
	parts := strings.Split(s, "/")
	if len(parts) == 1 {
		return schema.GroupVersion{Version: parts[0]}
	}
	return schema.GroupVersion{Group: parts[0], Version: parts[1]}
}

func formatGVRForComparison(gvr schema.GroupVersionResource) string {
	if gvr.Group == "" {
		return gvr.Resource
	}
	return fmt.Sprintf("%s.%s", gvr.Resource, gvr.Group)
}
