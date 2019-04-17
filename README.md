# go-azuredevops

[![GoDoc](https://godoc.org/github.com/benmatselby/go-azuredevops/azuredevops?status.svg)](https://godoc.org/github.com/benmatselby/go-azuredevops/azuredevops)
[![Build Status](https://travis-ci.org/benmatselby/go-azuredevops.png?branch=master)](https://travis-ci.org/benmatselby/go-azuredevops)
[![codecov](https://codecov.io/gh/benmatselby/go-azuredevops/branch/master/graph/badge.svg)](https://codecov.io/gh/benmatselby/go-azuredevops)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/go-azuredevops?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/go-azuredevops)

This is a fork of [go-azuredevops](https://github.com/benmatselby/go-azuredevops), a Go client library for accessing the [Azure DevOps API](https://docs.microsoft.com/en-gb/rest/api/vsts/). This fork adapts a lot of code and style from the [go-github](https://github.com/google/go-github/) library.

## Services

There is partial implementation for the following services

* Boards
* Builds
* Favourites
* Iterations
* Pull Requests
* Service Events (webhooks)
* Work Items

## Usage

```go
import "github.com/benmatselby/go-azuredevops/azuredevops
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

## References 
* [Microsoft Azure Devops Rest API](https://github.com/MicrosoftDocs/vsts-rest-api-specs)
* [Microsoft NodeJS Azure Devops Client](https://github.com/Microsoft/azure-devops-node-api)

