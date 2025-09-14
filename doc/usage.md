# Usage

Corsair provides two ways to forward requests:

- Configured requests
- Generic Forwarding

## Configured requests

Configure corsair endpoints in `config.yaml` to define specific proxy routes with custom headers, query parameters, and timeout settings.

### Endpoint Configuration

```yaml
endpoints:
  - 
    path: /api
    remote_url: https://api.example.com
    timeout: "30s"  # Optional: Override global default timeout
    headers:
      - X-Custom-Header: "value"
    query_params:
      - version: "v1"
```

> [!NOTE]
> You can use a request bin / webhook tester service as remote_url to view the generated request details, and test your configuration.

### Sub-path Forwarding

Sub-paths are automatically preserved and appended to the remote URL:

- Request: `GET /api/users/123`
- Endpoint path: `/api`
- Remote URL: `https://api.example.com`
- **Forwarded to**: `https://api.example.com/users/123`

The endpoint path is stripped from the request path, and the remaining sub-path is appended to the remote URL.

### Headers Handling

**Original headers are preserved**: All headers from the original request are forwarded to the upstream server.

**Custom headers override**: Headers defined in the endpoint configuration will override any matching headers from the original request.

**Template variables**: Use `{{ ENV_VAR }}` syntax in header values to inject environment variables at runtime:

```yaml
headers:
  - Authorization: "Bearer {{ API_TOKEN }}"
  - X-Request-ID: "{{ REQUEST_ID }}"
```

### Query Parameters Handling

**Original query parameters are preserved**: All query parameters from the original request are forwarded.

**Custom parameters override**: Query parameters defined in the endpoint configuration will override any matching parameters from the original request using `Set()`.

**Template variables supported**: Use `{{ ENV_VAR }}` syntax in parameter values:

```yaml
query_params:
  - api_key: "{{ API_KEY }}"
  - version: "v2"
```

## Generic Forwarding

The `/forward` endpoint allows ad-hoc proxying to any URL without pre-configuration.

### Usage

```bash
curl "http://localhost:8080/forward?url=https://api.example.com/data"
```

> [!NOTE]
>
> - **HTTPS default**: URLs without a scheme default to HTTPS for security
> - **Protocol restriction**: Only HTTP and HTTPS protocols are allowed
> - **URL validation**: Malformed URLs are rejected with `400 Bad Request`

### Example

```bash
# Forward to external API
curl "http://localhost:8080/forward?url=https://jsonplaceholder.typicode.com/posts/1"

# Auto-HTTPS (url becomes https://example.com)
curl "http://localhost:8080/forward?url=example.com/api"
```

## Configuration Options

### Server Settings

```yaml
server:
  address: "localhost"              # Bind address (default: localhost)
  port: 8080                       # Server port (default: 8080)
  forward_endpoint_enabled: true   # Enable/disable /forward endpoint
  default_timeout: "10s"           # Default timeout for all requests
```

### Timeout Configuration

- **Global default**: Set `server.default_timeout` for all endpoints
- **Per-endpoint override**: Set `timeout` in individual endpoint configurations
- **Format**: Duration strings like `"10s"`, `"1m30s"`, `"2m"`
- **Forward endpoint**: Uses global default timeout

### Template Variables

Template processing is applied on remote url, headers and query params, enabling secure secret injection:

```yaml
endpoints:
  - path: /secure
    remote_url: "{{ API_BASE_URL }}"
    headers:
      - Authorization: "Bearer {{ JWT_TOKEN }}"
      - X-Environment: "{{ ENVIRONMENT }}"
    query_params:
      - client_id: "{{ CLIENT_ID }}"
```

## Error Handling

- **Invalid configuration**: Server fails to start with validation errors
- **Upstream errors**: HTTP status codes and response bodies forwarded as-is
- **Network errors**: Returns `502 Bad Gateway` for connection failures
