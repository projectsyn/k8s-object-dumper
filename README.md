# openshift-pre-backup

Export all OpenShift objects before backup.


## Developing

1. Hack
2. `./openshift-pre-backup -v -d ./tmp`


### Building

    docker run -it --rm \
      -v "$(pwd)":/home/builder/rpmbuild/SOURCES \
      mhutter/rpmbuild \
      rpmbuild -bb rpmbuild/SPECS/openshift-pre-backup.spec
