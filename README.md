# xurlbits

![made with go](https://img.shields.io/badge/made%20with-Go-0000FF.svg) [![release](https://img.shields.io/github/release/hueristiq/xurlbits?style=flat&color=0000FF)](https://github.com/hueristiq/xurlbits/releases) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=0000FF)](https://github.com/hueristiq/xurlbits/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0000FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/xurlbits.svg?style=flat&color=0000FF)](https://github.com/hueristiq/xurlbits/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/xurlbits.svg?style=flat&color=0000FF)](https://github.com/hueristiq/xurlbits/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-0000FF.svg)](https://github.com/hueristiq/xurlbits/blob/master/CONTRIBUTING.md)

`xurlbits` is a command-line interface (CLI) utility to pull out bits of URLs.

## Resources

* [Installation](#installation)
	* [Install release binaries (Without Go Installed)](#install-release-binaries-without-go-installed)
	* [Install source (With Go Installed)](#install-source-with-go-installed)
		* [`go install ...`](#go-install)
		* [`go build ...` the development Version](#go-build--the-development-version)
* [Usage](#usage)
	* [Examples](#examples)
		* [Domains](#domains)
		* [Apex Domains](#apex-domains)
		* [Paths](#paths)
		* [Query String Key/Value Pairs](#query-string-keyvalue-pairs)
		* [Query String Keys (Parameters)](#query-string-keys-parameters)
		* [Query String Values](#query-string-values)
		* [Custom Formats](#custom-formats)
* [Contributing](#contributing)
* [Licensing](#licensing)
* [Credits](#credits)
	* [Contributors](#contributors)
	* [Similar Projects](#similar-projects)

## Installation

### Install release binaries (without Go installed)

Visit the [releases page](https://github.com/hueristiq/xtee/releases) and find the appropriate archive for your operating system and architecture. Download the archive from your browser or copy its URL and retrieve it with `wget` or `curl`:

* ...with `wget`:

	```bash
	wget https://github.com/hueristiq/xtee/releases/download/v<version>/xtee-<version>-linux-amd64.tar.gz
	```

* ...or, with `curl`:

	```bash
	curl -OL https://github.com/hueristiq/xtee/releases/download/v<version>/xtee-<version>-linux-amd64.tar.gz
	```

...then, extract the binary:

```bash
tar xf xtee-<version>-linux-amd64.tar.gz
```

> [!TIP]
> The above steps, download and extract, can be combined into a single step with this onliner
> 
> ```bash
> curl -sL https://github.com/hueristiq/xtee/releases/download/v<version>/xtee-<version>-linux-amd64.tar.gz | tar -xzv
> ```

> [!NOTE]
> On Windows systems, you should be able to double-click the zip archive to extract the `xtee` executable.

...move the `xtee` binary to somewhere in your `PATH`. For example, on GNU/Linux and OS X systems:

```bash
sudo mv xtee /usr/local/bin/
```

> [!NOTE]
> Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `xtee` to their `PATH`.

### Install source (with Go installed)

Before you install from source, you need to make sure that Go is installed on your system. You can install Go by following the official instructions for your operating system. For this, we will assume that Go is already installed.

#### `go install ...`

```bash
go install -v github.com/hueristiq/xtee/cmd/xtee@latest
```

#### `go build ...` the development version

* Clone the repository

	```bash
	git clone https://github.com/hueristiq/xtee.git 
	```

* Build the utility

	```bash
	cd xtee/cmd/xtee && \
	go build .
	```

* Move the `xtee` binary to somewhere in your `PATH`. For example, on GNU/Linux and OS X systems:

	```bash
	sudo mv xtee /usr/local/bin/
	```

	Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `xtee` to their `PATH`.


> [!CAUTION]
> While the development version is a good way to take a peek at `xtee`'s latest features before they get released, be aware that it may have bugs. Officially released versions will generally be more stable.

## Usage

To display help message for xurlbits use the `-h` flag:

```bash
$ xurlbits -h
```

help message:

```text
                 _ 
__  ___   _ _ __| |
\ \/ / | | | '__| |
 >  <| |_| | |  | |
/_/\_\\__,_|_|  |_| v0.0.0

A CLI utility to pull out bits of URLs.

USAGE:
  xurlbits [MODE] [FORMATSTRING] [OPTIONS]

INPUT:
  -i, --input       input file (use `-` to get from stdin)

OUTPUT:
  -m, --monochrome  disable output content coloring
  -u, --unique      output unique values
  -v, --verbosity   debug, info, warning, error, fatal or silent (default: info)

MODE:
  domains           the hostname (e.g. sub.example.com)
  apexes            the apex domain (e.g. example.com from sub.example.com)
  paths             the request path (e.g. /users)
  query             `key=value` pairs from the query string (one per line)
  params            keys from the query string (one per line)
  values            values from the query string (one per line)
  format            custom format (see below)

FORMAT DIRECTIVES:
  %%                a literal percent character
  %s                the request scheme (e.g. https)
  %u                the user info (e.g. user:pass)
  %d                the domain (e.g. sub.example.com)
  %S                the subdomain (e.g. sub)
  %r                the root of domain (e.g. example)
  %t                the TLD (e.g. com)
  %P                the port (e.g. 8080)
  %p                the path (e.g. /users)
  %e                the path's file extension (e.g. jpg, html)
  %q                the raw query string (e.g. a=1&b=2)
  %f                the page fragment (e.g. page-section)
  %@                inserts an @ if user info is specified
  %:                inserts a colon if a port is specified
  %?                inserts a question mark if a query string exists
  %#                inserts a hash if a fragment exists
  %a                authority (alias for %u%@%d%:%P)

EXAMPLES:
  cat urls.txt | xurlbits params -i -
  cat urls.txt | xurlbits format %s://%h%p?%q -i -
```

### Examples

```
$ cat urls.txt

https://sub.example.com/users?id=123&name=Sam
https://sub.example.com/orgs?org=ExCo#about
http://example.net/about#contact
```

#### Domains

You can extract the domains from the URLs with the `domains` mode:

```
$ cat urls.txt | xurlbits domains -i -

sub.example.com
sub.example.com
example.net
```

If you don't want to output duplicate values you can use the `-u` or `--unique` flag:

	```
	$ cat urls.txt | xurlbits domains  -i - --unique
	sub.example.com
	example.net
	```

The `-u`/`--unique` flag works for all modes.

#### Apex Domains

You can extract the apex part of the domain (e.g. the `example.com` in `http://sub.example.com`) using the `apexes` mode:

```
$ cat urls.txt | unfurl apexes -i - -u
example.com
example.net
```

#### Paths

```
$ cat urls.txt | xurlbits paths -i -

/users
/orgs
/about
```

#### Query String Key/Value Pairs

```
$ cat urls.txt | xurlbits query -i -

id=123
name=Sam
org=ExCo
```

#### Query String Keys (Parameters)

```
$ cat urls.txt | xurlbits params -i -

id
name
org
```

#### Query String Values

```
$ cat urls.txt | xurlbits values -i -

123
Sam
ExCo
```

#### Custom Formats

You can use the `format` mode to specify a custom output format:

```
$ cat urls.txt | xurlbits format %d%p -i -

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

> For more format directives, checkout the help message `xurlbits -h` under `Format Directives`. 

Any characters that don't match a format directive remain untouched:

```
$ cat urls.txt | xurlbits format "%d (%s)"  -i - -u

sub.example.com (https)
example.net (http)
```

**Note** that if a URL does not include the data requested, there will be no output for that URL:

```
$ echo http://example.com | xurlbits format "%P"  -i -

$ echo http://example.com:8080 | xurlbits format "%P" -i -
8080
```

## Contributing

We welcome contributions! Feel free to submit [Pull Requests](https://github.com/hueristiq/xtee/pulls) or report [Issues](https://github.com/hueristiq/xtee/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/xtee/blob/master/CONTRIBUTING.md).

## Licensing

This utility is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/xtee/blob/master/LICENSE).

## Credits

### Contributors

A huge thanks to all the contributors who have helped make `xtee` what it is today!

[![contributors](https://contrib.rocks/image?repo=hueristiq/xtee&max=500)](https://github.com/hueristiq/xtee/graphs/contributors)

### Similar Projects

If you're interested in more utilities like this, check out:

* [Tom Hudson](https://github.com/tomnomnom), we took the initial code from his [unfurl](https://github.com/tomnomnom/unfurl).