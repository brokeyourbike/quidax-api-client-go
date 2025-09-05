# quidax-api-client-go

[![Go Reference](https://pkg.go.dev/badge/github.com/brokeyourbike/quidax-api-client-go.svg)](https://pkg.go.dev/github.com/brokeyourbike/quidax-api-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/brokeyourbike/quidax-api-client-go)](https://goreportcard.com/report/github.com/brokeyourbike/quidax-api-client-go)

Quidax API Client for Go

## Installation

```bash
go get github.com/brokeyourbike/quidax-api-client-go
```

## Usage

```go
client := quidax.NewClient("token", signer)

_, err := client.FetchParentAccount(context.TODO())
require.NoError(t, err)
```

## Authors
- [Ivan Stasiuk](https://github.com/brokeyourbike) | [Twitter](https://twitter.com/brokeyourbike) | [LinkedIn](https://www.linkedin.com/in/brokeyourbike) | [stasi.uk](https://stasi.uk)

## License
[BSD-3-Clause License](https://github.com/brokeyourbike/quidax-api-client-go/blob/main/LICENSE)

