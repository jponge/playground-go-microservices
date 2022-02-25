# Dashboard web app talking to the store HTTP API 

## Notes

- Uses the `net/http` server plus some Gorilla middleware
- Embeds and serves static content in the final executable
- Viper to load the configuration through some YAML configuration
- Proxy requests to a temperature store service API
- Graceful shutdowns

### Generating keys

```text
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits=4096 --host=localhost
```

## Usage

Just run then go to https://localhost:4000/