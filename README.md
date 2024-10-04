# Bucket Proxy

A simple read bucket proxy. Currently only support for Google Cloud Storage. Feel free to contribute!

## Usage

```console
$ bucket-proxy --help

A read bucket proxy server

Usage:
  bucket-proxy [flags]

Flags:
      --bucket.gcs.credentials-file string   Path to the credentials file
      --bucket.name string                   Bucket name
      --bucket.provider string               Bucket provider must be one of: gcs (default "gcs")
  -h, --help                                 help for bucket-proxy
      --log.disable-quote                    Disable quote, this option only works for text format.
      --log.format string                    Log format (default "text")
      --log.level string                     Log level (default "info")
      --proxy.index-page string              Index page. Only used when web mode is enabled (default "index.html")
      --proxy.not-found-page string          Not found page. Only used when web mode is enabled (default "404.html")
      --proxy.web-mode                       Enable web mode. If enabled, the proxy will serve index and not found pages
      --server.address string                Server listen address (default "0.0.0.0:8080")
      --server.drain-timeout duration        Server drain timeout (default 15s)
  -v, --version                              Print version and exit

```
