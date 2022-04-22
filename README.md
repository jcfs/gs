<h1 align="center">
GS
</h1>

<p align="center">
A fast and lightweight scanner
</p>

<p align="center">
<img width="598" alt="image" src="https://user-images.githubusercontent.com/1512601/164565360-8f4963b8-42fe-4a66-afac-182defc5472e.png">
</p>

<hr/>
<img alt="image" src="https://github.com/jcfs/gs/actions/workflows/build.yml/badge.svg">
<img alt="image" src="https://github.com/jcfs/gs/actions/workflows/test.yml/badge.svg">

**GS** is aims to be a fast and lightweight scanner (of several types; port/domain enum). **GS** intends to be a simple tool
that can be used on automated scripts and checks; it is possible to get the output formatted as json/xml to ease integration
with existing tooling.

# Installation
As of now, the only way to install **GS** is from source, so the repo must be cloned and the code compiled. In the future I intend to have multi-arch binary releases available.

```
$ git clone https://github.com/jcfs/gs ; cd gs
$ make
$ make install
```

# Usage
Full option list can be found by running the command `gs` without any arguments:
```
Usage: gs [options...] <host>
 -t, --type      <type>  The scan type (domain, port, ...)
 -v, --verbose           Show more information while scanning
 -f, --format    <type>  Report output format (text, json, xml)
 -o, --output    <file>  Write output to file
DOMAIN:
 -s, --subdomain <data>  The subdomain to test (ie: www,ns1,cloud)
 -w, --wordlist  <file>  The word list to use
PORT:
 -p, --port      <ports> The port(s) to scan (ex: "1", "1-10", "1,2,3")
 ```
#### General
There are two types of scan currently implemented, the port scan and the subdomain enumeration scan. The port scan will 
do a connect-scan to the specified ports on the target host. The domain enumeration scan will enumerate all the valid 
subdomains of a given domain based on an argument or file list.

There are several report output formats (`--format`) available: `text` (default), `json` and `xml`.

#### Examples
* `gs localhost`       - scans the common ports on localhost
* `gs -p 80 localhost` - scans port 80 on localhost
* `gs -p 80,81-90,443,8080-9090 localhost` scans port 80, 81 to 90, 443 and 8080 to 9090 (the range includes both limits)
* `gs -p 80 -f json localhost` - scans port 80 on localhost and outputs the report in json
* `gs -t domain -s www google.com` - checks if www is a valid subdomain for google.com
* `gs -t domain -w path/to/file google.com` - checks all the subdomains on SUBDOMAIN.google.com