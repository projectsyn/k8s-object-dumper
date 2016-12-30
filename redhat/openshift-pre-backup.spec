#
# Exports all OpenShift resources prior to backup
#
Summary: Exports all OpenShift resources prior to backup
Name: openshift-pre-backup
Version: 1.10.3
Release: 1
License: BSD 3-Clause License
Group: Applications/System
Source: Makefile
Source1: openshift-pre-backup
URL: https://git.vshn.net/vshn/openshift-pre-backup
Vendor: VSHN AG
Packager: Manuel Hutter <manuel.hutter@vshn.ch>
Requires: bash, jq, moreutils, krossa

%description
Simple pre-backup script for OpenShift. Exports all OpenShift resources to
/var/lib/openshift-backup.

%prep
%setup -cT
cp -a %SOURCE0 .
cp -a %SOURCE1 .

%build
make SBINDIR='%{_sbindir}' 'LIBDIR=%{_libdir}' 'DATADIR=%{_datadir}'

%install
%make_install SBINDIR='%{_sbindir}' 'LIBDIR=%{_libdir}' 'DATADIR=%{_datadir}'

%files
%{_sbindir}/openshift-pre-backup

%changelog
* Thu Nov 21 2019 Simon Gerber <simon.gerber@vshn.ch> 1.10.3-1
- Ignore unretrievable types mutations and validations.

* Wed Sep 27 2019 Christian Haeusler <christian.haeusler@vshn.ch> 1.10.2-1
- Expect CronJob type to be present.

* Wed Sep 26 2019 Christian Haeusler <christian.haeusler@vshn.ch> 1.10.1-1
- Up until now, only the preferred version was included. Now the backup
  includes all API group versions.

* Wed Jul 17 2019 Michael Hanselmann <hansmi@vshn.ch> 1.10.0-1
- The "extract-objects" program has been replaced by a separate program named
  "krossa" which no longer requires Python.

* Mon Jul 15 2019 Michael Hanselmann <hansmi@vshn.ch> 1.9.0-1
- Output from the "openshift-pre-backup" script is now also written to a file
  named "log" in the output directory.

* Wed Apr 12 2019 Gabriel Mainberger <gabriel.mainberger@vshn.ch> 1.8.0-1
- Add support for OpenShift 3.11
- Removed support for OpenShift <= 3.6

* Tue Apr 9 2019 Michael Hanselmann <hansmi@vshn.ch> 1.7.1-2
- Change "extract-objects" program to explicitly invoke "/usr/bin/python3.4"
  instead of "/usr/bin/python3" to handle differences in the installed default
  Python version.

* Tue Mar 12 2019 Michael Hanselmann <hansmi@vshn.ch> 1.7.1-0
- Ensure decimal values are transferred correctly when splitting JSON files

* Wed Jan 16 2019 Michael Hanselmann <hansmi@vshn.ch> 1.7.0-0
- Add support for OpenShift 3.10
- Increase delays on failures
- Retry API discovery requests

* Mon Mar 12 2018 Michael Hanselmann <hansmi@vshn.ch> 1.6.0-0
- Add support for OpenShift 3.7

* Thu Aug 31 2017 Michael Hanselmann <hansmi@vshn.ch> 1.5.0-0
- Versions older than OpenShift 3.5 are no longer supported
- Bugfixes for OpenShift 3.5 and 3.6

* Mon Jun 26 2017 Michael Hanselmann <hansmi@vshn.ch> 1.4.4-0
- Handle larger clusters by increasing allowed number of open files when
  splitting object lists.
- Skip "poddisruptionbudgets" objects only on OpenShift 1.5/3.5 and newer.

* Mon Jun 26 2017 Michael Hanselmann <hansmi@vshn.ch> 1.4.2-0
- Recognize unretrievable types from OpenShift 3.5.

* Mon Jun 26 2017 Michael Hanselmann <hansmi@vshn.ch> 1.4.1-0
- Split lists of K8s objects into per-namespace files stored in
  "/var/lib/openshift-backup/split".

* Mon Feb 27 2017 Michael Hanselmann <hansmi@vshn.ch> 1.3.2-0
- Recognize error message prefix only emitted by OpenShift 1.4 and newer.

* Mon Feb 27 2017 Michael Hanselmann <hansmi@vshn.ch> 1.3.1-0
- Remove output files from previous backup.

* Mon Feb 27 2017 Michael Hanselmann <hansmi@vshn.ch> 1.3.0-0
- Rewrite to make use of OpenShift/Kubernetes API to get a complete list
  of object types.

* Tue Feb 21 2017 Manuel Hutter <manuel.hutter@vshn.ch> 1.2.0
- Exit with exit code 1 if errors occured during backup

* Tue Feb 21 2017 Manuel Hutter <manuel.hutter@vshn.ch> 1.1.2
- APPU-271: Don't abort if a single project or resource cannot be exported

* Mon Feb 20 2017 Peter H. Ruegg <peter.ruegg@vshn.ch> 1.1.1-1
- Send errors to STDERR (VSHNOPS-509)

* Fri Jan 06 2017 Manuel Hutter <manuel.hutter@vshn.ch> 1.1.0
- Fixed: script was always exporting the same project
- Fixed: export more resource types
- Changed: use `oc export` instead of `oc get`
- Added: Delete old backups before exporting

* Tue Jan 03 2017 Manuel Hutter <manuel.hutter@vshn.ch> 1.0.0
- Initial release
