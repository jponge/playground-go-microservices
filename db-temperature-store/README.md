# A HTTP service that stores temperature update data in a database

## Notes

- Uses `viper` and `cobra` for CLI and configuration
- Uses the `gorm` ORM with a `sqlite3` backend
- Uses `net/http` with [Chi](github.com/go-chi/chi) as a modern *muxer*.

## Usage

Same as `simple-temperature-store`.