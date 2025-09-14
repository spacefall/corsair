# `CORSair`: 🏴‍☠️ CORS HTTP Proxy

`CORSair` is a lightweight, configurable HTTP proxy that provides flexible request forwarding capabilities. It's designed primarily as a CORS proxy.

## Quick start

- [Installation](doc/install.md)
- [Configuration](doc/configuration.md)
- [usage](doc/usage.md)

## Overview

`CORSair` dynamically exposes endpoints based on a yaml configuration file.

Configuring an endpoint look like this:

```yaml
- path: /example
  remote_url: https://any.url/anything/
  headers:
    - X-Hello: "{{ WORLD }}" # Double braced expression are used for environment variable substitution.
  query_params:
    - foo: "bar"
```

Here, a query to `http://<your-corsair-url>/example` will be forwarded to `https://any.url.com/anything/?foo=bar` with configured headers. The global CORS header are added to the response (if configured).

See the [configuration reference](doc/configuration.md) for all options.

### Generic forwarding

It's also possible to use `/forward` endpoint to forward a request without any configuration, settings the remote url as query parameter: `http://localhost:8080/forward?url=https://any.url.com/anything/`

> [!CAUTION]
> **Forward Endpoint Risk**: The `/forward?url=` endpoint can potentially access internal networks. Consider disabling it (`forward_endpoint_enabled: false`) if you don't need it

## License

[MIT](https://opensource.org/licenses/MIT)
