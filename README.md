# xurl

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/xurl)](https://goreportcard.com/report/github.com/hueristiq/xurl) [![release](https://img.shields.io/github/release/hueristiq/xurl?style=flat&color=1E90FF)](https://github.com/hueristiq/xurl/releases) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/xurl.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/xurl/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/xurl.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/xurl/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/xurl/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/xurl/blob/master/CONTRIBUTING.md)

`xurl` is a command-line utility designed to extract specific parts of URLs.

## Resources

- [Features](#features)
- [Installation](#installation)
	- [Install release binaries (Without Go Installed)](#install-release-binaries-without-go-installed)
	- [Install source (With Go Installed)](#install-source-with-go-installed)
		- [`go install ...`](#go-install)
		- [`go build ...` the development Version](#go-build--the-development-version)
- [Usage](#usage)
	- [Examples](#examples)
		- [Domains](#domains)
		- [Apex Domains](#apex-domains)
		- [Paths](#paths)
		- [Query String Key/Value Pairs](#query-string-keyvalue-pairs)
		- [Query String Keys (Parameters)](#query-string-keys-parameters)
		- [Query String Values](#query-string-values)
		- [Custom Formats](#custom-formats)
- [Contributing](#contributing)
- [Licensing](#licensing)

## Features

- Multiple Extraction Modes
- Custom Formats
- Cross-Platform (Windows, Linux, and macOS)

## Installation

### Install release binaries (without Go installed)

Visit the [releases page](https://github.com/hueristiq/xurl/releases) and find the appropriate archive for your operating system and architecture. Download the archive from your browser or copy its URL and retrieve it with `wget` or `curl`:

- ...with `wget`:

	```bash
	wget https://github.com/hueristiq/xurl/releases/download/v<version>/xurl-<version>-linux-amd64.tar.gz
	```

- ...or, with `curl`:

	```bash
	curl -OL https://github.com/hueristiq/xurl/releases/download/v<version>/xurl-<version>-linux-amd64.tar.gz
	```

...then, extract the binary:

```bash
tar xf xurl-<version>-linux-amd64.tar.gz
```

> [!TIP]
> The above steps, download and extract, can be combined into a single step with this onliner
> 
> ```bash
> curl -sL https://github.com/hueristiq/xurl/releases/download/v<version>/xurl-<version>-linux-amd64.tar.gz | tar -xzv
> ```

> [!NOTE]
> On Windows systems, you should be able to double-click the zip archive to extract the `xurl` executable.

...move the `xurl` binary to somewhere in your `PATH`. For example, on GNU/Linux and OS X systems:

```bash
sudo mv xurl /usr/local/bin/
```

> [!NOTE]
> Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `xurl` to their `PATH`.

### Install source (with Go installed)

Before you install from source, you need to make sure that Go is installed on your system. You can install Go by following the official instructions for your operating system. For this, we will assume that Go is already installed.

#### `go install ...`

```bash
go install -v github.com/hueristiq/xurl/cmd/xurl@latest
```

#### `go build ...` the development version

- Clone the repository

	```bash
	git clone https://github.com/hueristiq/xurl.git 
	```

- Build the utility

	```bash
	cd xurl/cmd/xurl && \
	go build .
	```

- Move the `xurl` binary to somewhere in your `PATH`. For example, on GNU/Linux and OS X systems:

	```bash
	sudo mv xurl /usr/local/bin/
	```

	Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `xurl` to their `PATH`.


> [!CAUTION]
> While the development version is a good way to take a peek at `xurl`'s latest features before they get released, be aware that it may have bugs. Officially released versions will generally be more stable.

## Usage

To display help message for xurl use the `-h` flag:

```bash
xurl -h
```

help message:

```text

                 _                              _    _____
__  ___   _ _ __| |_   _ _ __  _ __   __ _  ___| | _|___ / _ __
\ \/ / | | | '__| | | | | '_ \| '_ \ / _` |/ __| |/ / |_ \| '__|
 >  <| |_| | |  | | |_| | | | | |_) | (_| | (__|   < ___) | |
/_/\_\\__,_|_|  |_|\__,_|_| |_| .__/ \__,_|\___|_|\_\____/|_|
                              |_|                         v0.1.0

USAGE:
 xurl [MODE] [FORMATSTRING] [OPTIONS]

MODES:
 domains                   the hostname (e.g. sub.example.com)
 apexes                    the apex domain (e.g. example.com from sub.example.com)
 paths                     the request path (e.g. /users)
 query                     `key=value` pairs from the query string (one per line)
 params                    keys from the query string (one per line)
 values                    query string values (one per line)
 format                    custom format (see below)

FORMAT DIRECTIVES:
  %%                       a literal percent character
  %s                       the request scheme (e.g. https)
  %u                       the user info (e.g. user:pass)
  %d                       the domain (e.g. sub.example.com)
  %S                       the subdomain (e.g. sub)
  %r                       the root of domain (e.g. example)
  %t                       the TLD (e.g. com)
  %P                       the port (e.g. 8080)
  %p                       the path (e.g. /users)
  %e                       the path's file extension (e.g. jpg, html)
  %q                       the raw query string (e.g. a=1&b=2)
  %f                       the page fragment (e.g. page-section)
  %@                       inserts an @ if user info is specified
  %:                       inserts a colon if a port is specified
  %?                       inserts a question mark if a query string exists
  %#                       inserts a hash if a fragment exists
  %a                       authority (alias for %u%@%d%:%P)

INPUT:
 -u, --url string[]        target URL
 -l, --list string         target URLs list file path

TIP: For multiple input URLs use comma(,) separated value with `-u`,
     specify multiple `-u`, load from file with `-l` or load from stdin.

OUTPUT:
     --unique bool         output unique values
     --monochrome bool     display no color output
 -s, --silent bool         stdout values only output
 -v, --verbose bool        stdout verbose output

```

### Examples

```bash
cat urls.txt
```

```
https://sub.example.com/users?id=123&name=Sam
https://sub.example.com/orgs?org=ExCo#about
http://example.net/about#contact
example.net/about#contact
```

#### Domains

You can extract the domains from the URLs with the `domains` mode:

```bash
cat urls.txt | xurl domains 
```

```
sub.example.com
sub.example.com
example.net
example.net
```

If you don't want to output duplicate values you can use the `--unique` flag:

```bash
cat urls.txt | xurl domains --unique
```

```
sub.example.com
example.net
```

The `-u`/`--unique` flag works for all modes.

#### Apex Domains

You can extract the apex part of the domain (e.g. the `example.com` in `http://sub.example.com`) using the `apexes` mode:

```bash
cat urls.txt | xurl apexes --unique 
```

```
example.com
example.net
```

#### Paths

```bash
cat urls.txt | xurl paths --unique 
```

```
/users
/orgs
/about
```

#### Query String Key/Value Pairs

```bash
cat urls.txt | xurl query
```

```
id=123
name=Sam
org=ExCo
```

#### Query String Keys (Parameters)

```bash
cat urls.txt | xurl params
```

```
id
name
org
```

#### Query String Values

```bash
cat urls.txt | xurl values
```

```
123
Sam
ExCo
```

#### Custom Formats

You can use the `format` mode to specify a custom output format:

```bash
cat urls.txt | xurl format %d%p
```

```
sub.example.com/users
sub.example.com/orgs
example.net/about
```

The available format directives are:

```
%%  A literal percent character
%s  The request scheme (e.g. https)
%u  The user info (e.g. user:pass)
%d  The domain (e.g. sub.example.com)
%S  The subdomain (e.g. sub)
%r  The root of domain (e.g. example)
%t  The TLD (e.g. com)
%P  The port (e.g. 8080)
%p  The path (e.g. /users)
%e  The path's file extension (e.g. jpg, html)
%q  The raw query string (e.g. a=1&b=2)
%f  The page fragment (e.g. page-section)
%@  Inserts an @ if user info is specified
%:  Inserts a colon if a port is specified
%?  Inserts a question mark if a query string exists
%#  Inserts a hash if a fragment exists
%a  Authority (alias for %u%@%d%:%P)
```

> For more format directives, checkout the help message `xurl -h` under `Format Directives`. 

Any characters that don't match a format directive remain untouched:

```bash
cat urls.txt | xurl format "%d (%s)"
```

```
sub.example.com (https)
example.net (http)
```

**Note** that if a URL does not include the data requested, there will be no output for that URL:

## Contributing

Feel free to submit [Pull Requests](https://github.com/hueristiq/xurl/pulls) or report [Issues](https://github.com/hueristiq/xurl/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/xurl/blob/master/CONTRIBUTING.md).

Huge thanks to the [contributors](https://github.com/hueristiq/xurl/graphs/contributors) thus far!

![contributors](https://contrib.rocks/image?repo=hueristiq/xurl&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/xurl/blob/master/LICENSE).