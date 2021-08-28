# GctSubdomains

Tool for discover subdomains in Certificate Transparency logs using Google's Transparency Report.

### Requirements
[Go](https://golang.org/) environment ready to go.

### Install

```bash
go get github.com/lord3ver/gctsubdomains
```

### Usage

```txt

     ______     __  _____       __        __                      _
    / ____/____/ /_/ ___/__  __/ /_  ____/ /___  ____ ___  ____ _(_)___  _____
   / / __/ ___/ __/\__ \/ / / / __ \/ __  / __ \/ __ '__ \/ __ '/ / __ \/ ___/
  / /_/ / /__/ /_ ___/ / /_/ / /_/ / /_/ / /_/ / / / / / / /_/ / / / / (__  )
  \____/\___/\__//____/\__,_/_.___/\__,_/\____/_/ /_/ /_/\__,_/_/_/ /_/____/


        Google Transparencyreport subdomains finder

        Version:        0.9.0

Usage:

  -d string
        Target domain. E.g. bing.com
  -lookout
  	Do DNS lookups for the domains to see which ones exist (default true)
  -out
        Print results to stdout (default true)
  -outfile string
        Specify an output file when completed. Create or append if exists.
  -rmd
        Remove duplicates (default true)
  -rme
        Remove external domains, like xyz.com for uber.com domain (default false)
  -rmw
        Remove wildcard domains, ex. *.uber.com (default true)
```

At least one option among `out` and `outfile` must be specified. 

### Example
```
gctsubdomains -d uber.com --rmd=true --rme=true --out=true --outfile=output.txt
```
```
accessibility.uber.com
assets.uber.com
base.uber.com
assets-share.uber.com
beacon.uber.com
bizblog.uber.com
blog.uber.com
...
```

### Why?
- Why not? I needed it and on a summer night I did it.
- Search CT logs for unknown subdomains.
- Other subdomain scanners companion.
