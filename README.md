# K8s Object Dumper

Discover and dump all listable objects from a Kubernetes cluster into JSON files.
Written to be used as a pre backup command for [K8up](https://k8up.io).

## Usage

The project uses controller-runtime's configuration discovery to find the Kubernetes API server.



### Dump to STDOUT

```bash
$ k8s-object-dumper
{"apiVersion":"v1","kind":"List","items":[{"apiVersion":"v1", ...}]}
{"apiVersion":"v1","kind":"List","items":[{"apiVersion":"apps/v1", ...}]}
```

### Dump to a directory

```bash
$ k8s-object-dumper -dir dir
```

Will result in the following directory structure:

```
└─ dir/
   ├─ objects-<kind>[.<group>].json
   ├─ …
   └─ split/
      ├─ <namespace>/
      |  ├─ __all__.json
      |  ├─ <kind>[.<group>].json
      |  └─ …
      └─ …
```

### Advanced usage

```bash
# Fail if a Pods, Deployments or AlertingRules are not found
$ k8s-object-dumper \
  -must-exist=pods \
  -must-exist=deployments.apps \
  -must-exist=alertingrules.monitoring.openshift.io
# Ignore all Secrets and all cert-manager objects
$ k8s-object-dumper \
  -ignore=secrets \
  -ignore=.+cert-manager.io
```

## Development

The project uses [envtest](https://book.kubebuilder.io/reference/envtest) to run tests against a real Kubernetes API server.

```bash
$ make test
```

## Differences to the original `bash` version `< 0.3.0`

- All APIs are fully qualified in both the options (`--must-exist=certificates.cert-manager.io`, `--ignore=deployment.apps`) and the output files (`objects-Certificate.cert-manager.io.json`).
  This makes it possible to distinguish between objects with the same kind but different groups. See https://github.com/projectsyn/k8s-object-dumper/issues/47.
- Resources without a list endpoint are ignored, do not cause an error, and don't need to be explicitly ignored.
- Ignore and must-exist options are now command line flags instead of files in `/usr/local/share`.

## Contributing and license

This library is licensed under [BSD-3-Clause](LICENSE).
For information about how to contribute see [CONTRIBUTING](CONTRIBUTING.md).
