# go-azuredevops

**6 May 2019 Currently in development - API should be stable soon**

[![GoDoc](https://godoc.org/github.com/mcdafydd/go-azuredevops/azuredevops?status.svg)](https://godoc.org/github.com/mcdafydd/go-azuredevops/azuredevops)
[![Build Status](https://travis-ci.org/mcdafydd/go-azuredevops.png?branch=master)](https://travis-ci.org/mcdafydd/go-azuredevops)
[![codecov](https://codecov.io/gh/mcdafydd/go-azuredevops/branch/master/graph/badge.svg)](https://codecov.io/gh/mcdafydd/go-azuredevops)
[![Go Report Card](https://goreportcard.com/badge/github.com/mcdafydd/go-azuredevops?style=flat-square)](https://goreportcard.com/report/github.com/mcdafydd/go-azuredevops)

This is a fork of https://github.com/benmatselby/go-azuredevops, a Go library for accessing the [Azure DevOps API](https://docs.microsoft.com/en-gb/rest/api/vsts/).  As it develops, the library is looking more like of a port of [go-github](https://github.com/google/go-github).

This is my first project in Go, so I'm learning as I code.

## Services

There is partial implementation for the following services:

* Boards
* Builds
* Favourites
* Git
* Iterations
* Pull Requests
* Service Events (webhooks)
* Tests
* Users
* Work Items

## Usage

For usage with a personal access token, create a token using the process described here:

https://docs.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/pats?view=azure-devops

Add the token to a basic auth transport with an empty username:

```go
tp := azuredevops.BasicAuthTransport{
		Username: "",
		Password: token,
	}
```

Supply this token in calls to NewClient():

```go
import "github.com/mcdafydd/go-azuredevops/azuredevops

client, _ := azuredevops.NewClient(tp.Client())
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

### OAuth
Instead of using a personal access token related to your personal user account, consider registering your app in Azure Devops:

https://docs.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/oauth?view=azure-devops

Supply a token generated from this process when you call NewClient().

## Contributing
This library is re-using a lot of the code and style from the [go-github](https://github.com/google/go-github/) library:

* Receiving struct members should be pointers [google/go-github#19](https://github.com/google/go-github/issues/19)
* Debate on whether nillable fields are important for receiving structs, especially when also using separate request structs.  Other popular libraries also use pointers approach, but it is often viewed as a big ugly. [google/go-github#537](https://github.com/google/go-github/issues/537)
* Large receiving struct types should return []*Type, not []Type [google/go-github#375](https://github.com/google/go-github/pull/375)
* Use omitempty in receiving struct JSON tags
* Pass context in API functions and Do() [google/go-github#526](https://github.com/google/go-github/issues/526#issuecomment-280985393)  **UNDER DEVELOPMENT AS OF 6 MAY 2019**

An exception to the pointer members are the count/value receiving structs used for responses containing multiple entities.

May add separate request structs soon.

## References
* [Microsoft Azure Devops Rest API](https://github.com/MicrosoftDocs/vsts-rest-api-specs)
* [Microsoft NodeJS Azure Devops Client](https://github.com/Microsoft/azure-devops-node-api)

