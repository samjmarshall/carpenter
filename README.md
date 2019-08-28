# Carpenter

## Prerequisites

Dependencies of drivers, testers and provisioners that must be pre-installed.

**Packer** `>= 1.4.0`

**Vagrant**

  * vagrant `>= 2.2.0`
  * virtualbox `>= 6.0.0`

**Docker** `>= 18.06`

## Installation

Ensure all listed prerequisites are installed for the driver, tester and provisioner of your choosing.

Requires:

  * go `>= 1.12`

```bash
git clone https://github.com/rokertech/carpenter.git

export GOPATH="$HOME/go" >> ~/.bash_profile
export PATH="$PATH:$GOPATH/bin" >> ~/.bash_profile

source ~/.bash_profile

cd ./carpenter
go install
```

## Bash completion

To configure your bash shell to load completions for each session run.

```bash
carpenter completion >$(brew --prefix)/etc/bash_completion.d/carpenter
```

## Usage

Example image build, test and destroy using the example `.carpenter.yaml` config. Carpenter searches for it's config in the current working directory. To run the carpenter examples navigate to the examples directory:

```bash
cd ./carpenter/examples
```

**Build**
```bash
carpenter image build php-nginx
```

**Test** - Test is run immediately after build by default. This allows you to run tests against your build at any time, however, only when the target is accessible after build i.e. local Docker or Vagrant builds.
```bash
carpenter image test php-nginx
```

**Destroy**
```bash
carpenter image destroy php-nginx
```

## Copyright

Copyright (c) 2019 Roker Labs. See [LICENSE](./LICENSE) for details.