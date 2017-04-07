# Helm Template Plugin

This is a Helm plugin to help chart developers debug their charts. It works like
`helm install --dry-run --debug`, except that it runs locally, has more output
options, and is quite a bit faster.

<a href="https://asciinema.org/a/8kuehzpx5xyl8cm3cairica8z" target="_blank"><img src="https://asciinema.org/a/8kuehzpx5xyl8cm3cairica8z.png" width="589"/></a>

## Usage

Render chart templates locally and display the output.

This does not require Tiller. However, any values that would normally be
looked up or retrieved in-cluster will be faked locally. Additionally, none
of the server-side testing of chart validity (e.g. whether an API is supported)
is done.

```
$ helm template [flags] CHART
```

### Flags:

```
      --notes               show the computed NOTES.txt file as well.
      --set string          set values on the command line. See 'helm install -h'
  -f, --values valueFiles   specify one or more YAML files of values (default [])
  -v, --verbose             show the computed YAML values as well.
```


## Install

```
$ helm plugin install https://github.com/tehnosophos/helm-template
```

### Prerequisites

- You need to have [Go](http://golang.org) installed. Make sure to set `$GOPATH`
- If you don't have [Glide](http://glide.sh) installed, this will install it into
  `$GOPATH/bin` for you.
