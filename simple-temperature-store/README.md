# A simple HTTP service that stores temperature update data

## Notes

- Uses `fiber` to design the HTTP API
  - Quite pleasant to use
  - Standalone library, does not integrate with `net/http`
- Data is stored in memory as a `map`, so a mutex is required to protect access
- Uses `viper` and `cobra` to have a CLI and get the http server configuration from flags, environment variables or an optional configuration file
  - In a bigger code base the `cobra` command would be moved to a `cmd/` folder / package
- Integration tests

## Usage

Get all data:

```text
$ http :3000/data
HTTP/1.1 200 OK
Content-Length: 74
Content-Type: application/json
Date: Tue, 01 Feb 2022 08:57:16 GMT

[
    {
        "sensorId": "123-abc",
        "value": 19.2
    },
    {
        "sensorId": "456-def",
        "value": -2.33
    }
]
```

Update some data:

```text
$ http POST :3000/record sensorId=1 value:=19.1
HTTP/1.1 200 OK
Content-Length: 29
Content-Type: application/json
Date: Tue, 01 Feb 2022 08:57:55 GMT

{
    "sensorId": "1",
    "value": 19.1
}
```

Get a specific sensor data:

```text
$ http :3000/data/1
HTTP/1.1 200 OK
Content-Length: 29
Content-Type: application/json
Date: Tue, 01 Feb 2022 08:58:41 GMT

{
    "sensorId": "1",
    "value": 19.1
}
```