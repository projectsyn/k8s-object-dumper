# K8s Object Dumper

K8s Object Dumper allows to collect all objects from Kubernetes and write them into files.
It is written to be used as a pre backup command for [K8up](https://k8up.io).

This repository is part of Project Syn.
For documentation on Project Syn, see https://syn.tools.

K8s Object Dumper consists of a shell script.
This script uses the Kubernetes API to list each and every API and Kind known to the targeted cluster.
It then dumps all those kinds to Json files (one file per kind).
For easier restore, those dumped objects then get split up by namespace and kind.

The resulting structure looks like the following:

```
├─ objects-<kind>.json
├─ …
└─ split/
   ├─ <namespace>/
   |  ├─ __all__.json
   |  ├─ <kind>.json
   |  └─ …
   └─ …
```

## Usage

Using Docker

```bash
docker run --rm -v "/path/to/kubeconfig:/kubeconfig" -e KUBECONFIG=/kubeconfig -v "${PWD}/data:/data" projectsyn/k8s-object-dumper:latest -d /data > objects.tar.gz
```

Using Kubernetes (with K8up)

See [Commodore Component: cluster-backup](https://github.com/projectsyn/component-cluster-backup).

## Configuration

The `dump-objects` scripts reads configuration from two files.

`/usr/local/share/k8s-object-dumper/must-exists` contains a list of types that must exist within the list of discovered types.
This is a safeguard helping to detect failure of the discover mechanism.
Types must be all lower case and plural.
One type per line.

Example:

```
configmaps
daemonsets
deployments
endpoints
ingresses
jobs
namespaces
nodes
persistentvolumeclaims
persistentvolumes
replicasets
roles
secrets
serviceaccounts
services
statefulsets
```

Some types can not be exported and the script will return an error for them.
Those errors can be suppressed by placing those types in `/usr/local/share/k8s-object-dumper/known-to-fail`.
Like `must-exist` types are listed line by line but in addition Bash regular expressions can be used.


Example:

```
.+mutators
.+reviews
.+validators
bindings
deploymentconfigrollbacks
imagesignatures
imagestream.+
mutations
useridentitymappings
validations
```

## Contributing and license

This library is licensed under [BSD-3-Clause](LICENSE).
For information about how to contribute see [CONTRIBUTING](CONTRIBUTING.md).
