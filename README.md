# go-azuredevops

*** Currently in development - may not be stable ***

[![GoDoc](https://godoc.org/github.com/mcdafydd/go-azuredevops/azuredevops?status.svg)](https://godoc.org/github.com/mcdafydd/go-azuredevops/azuredevops)
[![Build Status](https://travis-ci.org/mcdafydd/go-azuredevops.png?branch=master)](https://travis-ci.org/mcdafydd/go-azuredevops)
[![codecov](https://codecov.io/gh/mcdafydd/go-azuredevops/branch/master/graph/badge.svg)](https://codecov.io/gh/mcdafydd/go-azuredevops)
[![Go Report Card](https://goreportcard.com/badge/github.com/mcdafydd/go-azuredevops?style=flat-square)](https://goreportcard.com/report/github.com/mcdafydd/go-azuredevops)

This is a fork of [go-azuredevops](https://github.com/mcdafydd/go-azuredevops), a Go client library for accessing the [Azure DevOps API](https://docs.microsoft.com/en-gb/rest/api/vsts/). This fork adapts a lot of code and style from the [go-github](https://github.com/google/go-github/) library.

## Services

There is partial implementation for the following services

* Boards
* Builds
* Favourites
* Iterations
* Pull Requests
* Service Events (webhooks)
* Tests
* Work Items

## Usage

```go
import "github.com/mcdafydd/go-azuredevops/azuredevops
```

Construct a new Azure DevOps Client

```go
v := azuredevops.NewClient(account, project, token)
```

Get a list of iterations

```go
iterations, error := v.Iterations.List(team)
if error != nil {
    fmt.Println(error)
}

for index := 0; index < len(iterations); index++ {
    fmt.Println(iterations[index].Name)
}
```

## Contributing
This library is re-using a lot of the code and style from the [go-github](https://github.com/google/go-github/) library:

* Receiving struct members should be pointers [google/go-github#19](https://github.com/google/go-github/issues/19)
* Debate on whether nillable fields are important for receiving structs, especially when also using separate request structs.  Other popular libraries also use pointers approach, but it is often viewed as a big ugly. [google/go-github#537](https://github.com/google/go-github/issues/537)
* Large receiving struct types should return []*Type, not []Type [google/go-github#375](https://github.com/google/go-github/pull/375)
* Use omitempty in receiving struct JSON tags

An exception to the pointer members are the count/value receiving structs used for responses containing multiple entities.

May add separate request structs soon.

## References
* [Microsoft Azure Devops Rest API](https://github.com/MicrosoftDocs/vsts-rest-api-specs)
* [Microsoft NodeJS Azure Devops Client](https://github.com/Microsoft/azure-devops-node-api)

