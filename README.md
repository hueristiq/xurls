# xurls

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/xurls)](https://goreportcard.com/report/github.com/hueristiq/xurls) [![release](https://img.shields.io/github/release/hueristiq/xurls?style=flat&color=1E90FF)](https://github.com/hueristiq/xurls/releases) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/xurls.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/xurls/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/xurls.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/xurls/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/xurls/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/xurls/blob/master/CONTRIBUTING.md)

`xurls` is a command-line utility for working with URLs.

## Resources

- [Features](#features)
- [Installation](#installation)
	- [Install release binaries (Without Go Installed)](#install-release-binaries-without-go-installed)
	- [Install source (With Go Installed)](#install-source-with-go-installed)
		- [`go install ...`](#go-install)
		- [`go build ...` the development Version](#go-build--the-development-version)
- [Usage](#usage)
- [Contributing](#contributing)
- [Licensing](#licensing)

## Features

- Extract URLs from text
- Decompose URLs into granular component
- Cross-Platform (Windows, Linux, and macOS)

## Installation

### Install release binaries (without Go installed)

Visit the [releases page](https://github.com/hueristiq/xurls/releases) and find the appropriate archive for your operating system and architecture. Download the archive from your browser or copy its URL and retrieve it with `wget` or `curl`:

- ...with `wget`:

	```bash
	wget https://github.com/hueristiq/xurls/releases/download/v<version>/xurls-<version>-linux-amd64.tar.gz
	```

- ...or, with `curl`:

	```bash
	curl -OL https://github.com/hueristiq/xurls/releases/download/v<version>/xurls-<version>-linux-amd64.tar.gz
	```

...then, extract the binary:

```bash
tar xf xurls-<version>-linux-amd64.tar.gz
```

> [!TIP]
> The above steps, download and extract, can be combined into a single step with this onliner
> 
> ```bash
> curl -sL https://github.com/hueristiq/xurls/releases/download/v<version>/xurls-<version>-linux-amd64.tar.gz | tar -xzv
> ```

> [!NOTE]
> On Windows systems, you should be able to double-click the zip archive to extract the `xurls` executable.

...move the `xurls` binary to somewhere in your `PATH`. For example, on GNU/Linux and OS X systems:

```bash
sudo mv xurls /usr/local/bin/
```

> [!NOTE]
> Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `xurls` to their `PATH`.

### Install source (with Go installed)

Before you install from source, you need to make sure that Go is installed on your system. You can install Go by following the official instructions for your operating system. For this, we will assume that Go is already installed.

#### `go install ...`

```bash
go install -v github.com/hueristiq/xurls/cmd/xurls@latest
```

#### `go build ...` the development version

- Clone the repository

	```bash
	git clone https://github.com/hueristiq/xurls.git 
	```

- Build the utility

	```bash
	cd xurls/cmd/xurls && \
	go build .
	```

- Move the `xurls` binary to somewhere in your `PATH`. For example, on GNU/Linux and OS X systems:

	```bash
	sudo mv xurls /usr/local/bin/
	```

	Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `xurls` to their `PATH`.


> [!CAUTION]
> While the development version is a good way to take a peek at `xurls`'s latest features before they get released, be aware that it may have bugs. Officially released versions will generally be more stable.

## Usage

To display help message for xurls use the `-h` flag:

```bash
xurls -h
```

help message:

```text

                 _
__  ___   _ _ __| |___
\ \/ / | | | '__| / __|
 >  <| |_| | |  | \__ \
/_/\_\\__,_|_|  |_|___/
                 v0.0.0

Usage:
  xurls [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  extract     Command for extracting URLs from text.
  help        Help about any command
  parse       Command for parsing URLs into granular components.

Flags:
  -h, --help         help for xurls
      --monochrome   display no color output
  -s, --silent       stdout values only output
  -v, --verbose      stdout verbose output

Use "xurls [command] --help" for more information about a command.

```

## Contributing

Contributions are welcome and encouraged! Feel free to submit [Pull Requests](https://github.com/hueristiq/xurls/pulls) or report [Issues](https://github.com/hueristiq/xurls/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/xurls/blob/master/CONTRIBUTING.md).

A big thank you to all the [contributors](https://github.com/hueristiq/xurls/graphs/contributors) for your ongoing support!

![contributors](https://contrib.rocks/image?repo=hueristiq/xurls&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/xurls/blob/master/LICENSE).