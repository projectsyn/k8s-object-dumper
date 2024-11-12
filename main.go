package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/projectsyn/k8s-object-dumper/internal/pkg/discovery"
	"github.com/projectsyn/k8s-object-dumper/internal/pkg/dumper"
)

func main() {
	var dir string
	var chunkSize int64
	mustExistResources := new(repeatableStringFlag)
	ignoreResources := new(repeatableRegexpFlag)

	flag.StringVar(&dir, "dir", "", "Directory to dump objects into")
	flag.Int64Var(&chunkSize, "chunk-size", 500, "Chunk size for listing objects")
	flag.Var(mustExistResources, "must-exist", "Resource that must exist in the cluster. Can be used multiple times.")
	flag.Var(ignoreResources, "ignore", "Resource to ignore during discovery. Regexp, anchored by default. Can be used multiple times.")

	flag.Parse()

	df := dumper.DumpToWriter(os.Stdout)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directory %s: %v\n", dir, err)
			os.Exit(1)
		}
		d, err := dumper.NewDirDumper(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directory dumper: %v\n", err)
			os.Exit(1)
		}
		defer d.Close()
		df = d.Dump
	}

	conf, err := ctrl.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get Kubernetes config: %v", err)
	}

	if err := discovery.DiscoverObjects(context.Background(), conf, df, discovery.DiscoveryOptions{
		ChunkSize:          chunkSize,
		LogWriter:          os.Stderr,
		MustExistResources: *mustExistResources,
		IgnoreResources:    *ignoreResources,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to dump some or all objects: %+v\n", err)
		os.Exit(1)
	}
}

type repeatableStringFlag []string

func (i *repeatableStringFlag) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *repeatableStringFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type repeatableRegexpFlag []*regexp.Regexp

func (i *repeatableRegexpFlag) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *repeatableRegexpFlag) Set(value string) error {
	value = fmt.Sprintf("^%s$", value)
	r, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("failed to compile regexp %q: %w", value, err)
	}
	*i = append(*i, r)
	return nil
}
