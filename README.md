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
