# Che-Inspector
Cli solution written in golang using cobra to allow running tests in any container platform.

# Specifications
* An instrumented cli with cobra
* Structured logging with logrus.
* Use client-go to talk with kubernetes API Server.
* Use docker-sdk to talk with docker socket to create test containers.
* Collect artifacts from test containers
* Send test results to slack. WORK IN PROGRESS

# Setup

A proper setup Go workspace using **Go 1.13+ is required**.

Install dependencies:
```
# Install dependencies
$ go mod tidy
# Copy the dependencies to vendor folder
$ go mod vendor
# Create che-inspector binary in bin folder. Please add the binary to the path or just execute ./bin/che-inspector
$ make build
```

# Installation
Linux
```bash

# Get Operating system and arch to get the right release binaries from github
export ARCH=$(case $(arch) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(arch) ;; esac)
export OS=$(uname | awk '{print tolower($0)}')

# Download the binaries and add che-inspector to /usr/local/bin
export CHE_INSPECTOR_URL=https://github.com/flacatus/che-inspector/releases/latest/download
curl -LO ${CHE_INSPECTOR_URL}/$OS-$ARCH-che-inspector
chmod +x $OS-$ARCH-che-inspector && sudo mv $OS-$ARCH-che-inspector /usr/local/bin/che-inspector
```
## Che-inspector binary

### Run tests using samples
In samples directory exists some Eclipse Che QE samples which can be run with the che-inspector cli

Run tests.
````
che-inspector run --file=samples/happy-path.yaml
````
## Config file tests spec
Work in progress
