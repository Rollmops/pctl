# PCTL

[![Coverage Status](https://coveralls.io/repos/github/Rollmops/pctl/badge.svg?branch=master)](https://coveralls.io/github/Rollmops/pctl?branch=master)

## Ideas

- Configuration File (YAML)
  - dependencies between processes (e.g. depends_on)
    - check of dependencies (e.g. circle deps, ...)
        https://pm2.keymetrics.io/docs/usage/quick-start/#list-managed-applications
    - visualization as a dependency-tree
  - readiness-probe
    - default readiness-probe: check for running
    - http readiness-probe: define endpoint, ...
    - custom readiness probe: script
  - lifeness-probe
    - default readiness-probe: check for running
    - http readiness-probe: define endpoint, ...
    - custom readiness probe: script
  - stop-strategy
    - default: stop by signal (SIGTERM as default)
    - custom: stop by script (e.g. KM)
    - graceful shutdown timeout → kill
  - executable md5sum check
    - take the first fraction of the command, save md5sum
    - take /proc/<pid>/exe and save md5sum
    - this will make sure keeping track of script starts (e.g. python scripts)

