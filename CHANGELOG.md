# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.1.1] - 2021-02-08

### Added

- Add "You may not request a new project via this API" to potential fetch errors [#9]

### Changed

- Updated Docker base image to `debian:10.7-slim`

## [v0.1.0] - 2020-10-16

Release as open source code

### Added

- Made list of expected and known to fail types configurable ([#4])

### Changed

- Clean up code base ([#1])
- Stream archive before exit ([#5])
- Included more files in output ([#6])

### Removed

- Droped creation of output directory ([#3])

[Unreleased]: https://github.com/projectsyn/component-k8s-object-dumper/compare/v0.1.1...HEAD
[v0.1.0]: https://github.com/projectsyn/component-k8s-object-dumper/releases/tag/v0.1.0
[v0.1.1]: https://github.com/projectsyn/component-k8s-object-dumper/releases/tag/v0.1.0

[#1]: https://github.com/projectsyn/k8s-object-dumper/pull/1
[#3]: https://github.com/projectsyn/k8s-object-dumper/pull/3
[#4]: https://github.com/projectsyn/k8s-object-dumper/pull/4
[#5]: https://github.com/projectsyn/k8s-object-dumper/pull/5
[#6]: https://github.com/projectsyn/k8s-object-dumper/pull/6
[#9]: https://github.com/projectsyn/k8s-object-dumper/pull/9
