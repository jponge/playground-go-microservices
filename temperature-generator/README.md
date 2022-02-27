# A temperature data producer

## Notes

- Uses bare `net/http` to issue HTTP client request
- Uses `cobra` for the CLI but not `viper` for configuration
  - Uses an idiomatic `cobra` app packages layout
  - ...so commands are package-level global variables and `init()` functions
- The main thread handles `SIGINT` / `SIGTERM` signals, in a more elaborated app this would allow for graceful shutdowns

## Usage

```
Run the generator

Usage:
  temperature-generator run [flags]

Flags:
      --count int     The number of generators (default 5)
  -h, --help          help for run
      --host string   The target temperature store host (default "localhost")
      --port int      The target temperature store port (default 3000)
```