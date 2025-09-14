# Installation

Create the `config.yaml` file base on the [example in this repository](config.yaml). See the [configuration reference](doc/configuration.md) for more information on options. If the `-c` option is not provided, `corsair` will use `/etc/corsair/config.yaml` as default configuration file.

### Using Docker

```console
git clone https://github.com/bastienwirtz/corsair
docker build . -t corsair
docker run [-e MY_ENV='my-value'] -p 8080:8080 corsair
```

Use `-e` to configure env vars if used in configuration.

### Using docker compose (recommended)

```yaml
services:
  corsair:
    build: .
    volumes:
      - /path/to/config/dir:/etc/corsair
    ports:
      - 8080:8080
# Optionnal: Configure env vars if used in configuration:
    environment: 
      - MY_ENV-ENV=my-value 
    restart: unless-stopped
```

```console
docker compose up -d
```

### Local build

Build the server.

```bash
git clone https://github.com/bastienwirtz/corsair
cd corsair
go build -o corsair ./cmd
```

Start the server

```bash
# Optionnal: Configure env vars if used in configuration:
export MY_ENV="my-value"
./corsair [-c config.yaml]    # Uses custom config file
```
