## Configuration

### Server Configuration

```yaml
server:
  address: "localhost"              # Interface to bind to (default: localhost)  
  port: 8080                        # Server port (default: 8080)
  forward_endpoint_enabled: true    # Enable/disable /forward endpoint (default: true)
  default_timeout: "10s"            # Default timeout for external requests (default: 10s)
```

**Address Options:**
- `"localhost"` - Local access only (default, secure)
- `"0.0.0.0"` - All IPv4 interfaces (for containers/external access)
- `"127.0.0.1"` - IPv4 localhost only
- `"192.168.1.100"` - Specific IPv4 address
- `"::1"` - IPv6 localhost only
- `"::"` - All IPv6 interfaces
- `"2001:db8::1"` - Specific IPv6 address

### Logging Configuration

```yaml
logging:
  level: info       # Log level: debug, info, warn, error (default: info)
  format: text      # Log format: json, text, pretty (default: text)
```

**Format Options:**
`json`: structured logging using json format
`text`: structured logging using Logfmt format.
`pretty`: Human friendly colored structured logging using Logfmt format. (using [tint](https://github.com/lmittmann/tint) 🌈). Colors only enabled on tty.

### CORS Configuration

General [CORS Guide](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS)

```yaml
cors:
  allow_origins: ["*"]                 # Allowed origins (* for all)
  allow_methods: "GET, POST, OPTIONS"  # Allowed HTTP methods
  allow_headers: "*"                   # Allowed headers (* for all)
  allow_credentials: true              # Allow credentials
```

**Origin Examples:**
- `["*"]` - Allow all origins
- `["https://example.com"]` - Specific origin
- `["*.example.com"]` - Subdomain wildcard
- `["http://localhost:3000", "https://app.com"]` - Multiple origins

### Endpoint Configuration

```yaml
endpoints:
  - path: /api                           # Proxy path prefix
    remote_url: https://api.example.com  # Target server
    headers:                             # Additional headers
      - X-My-Header: "my header value"
    query_params:                        # Additional query parameters
      - name: "value"
      - foo: "bar"
```

### Template Variables

Use `{{ VARIABLE_NAME }}` syntax in remote url, headers or query params values to inject environment variables:

```yaml
remote_url: "https://{{ SERVER_HOST }}/api"
headers:
  - Authorization: "Basic {{ AUTH_TOKEN }}"
query_params:
  - api_key: "{{ API_KEY }}"
```
