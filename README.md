# Go Coding Challenge v1

This document presents the needed information to install and run the solution given for the Code Challenge. To better
understand the code there are some comments in it.

This solution is a simple http server written in go.

Since the domain and requirements are simple there's no need to over-complexify the solution, therefore I choose to:

- not use any libraries, go already provides the basics needed to do what this challenge requires;
- not follow a commonly used architecture such as onion, layered or clean architecture for the sake of simplicity.

## Build & Running

There's two different ways to run the solution:

- natively, requires golang v1.22+ to be installed;
- docker, requires docker to be installed.

Note that this has only been tested in linux.

### Natively

Ensure that the go compiler is available in your workspace.

To build the solution with go:

```shell
go build -o server .
```

To run the solution in port 8080:

```shell
./server
```

### Docker

To run the solution in port 8080:

```shell
docker compose up -d
```

## Requirements assumed

- page urls and visitor id can't be an empty string
- all non-empty visitor ids received are valid and exist within the deus.ai domain;
- a page url can be seen as a unique identifier, e.g.: https://example.org/page?query=x != http://example.org/page

## Possible improvements

## Notes

- the business rules of this code challenge are almost none existing, in a proper application I'd love for most business logic to
  live in the domain and have a proper 'services' layer, but with the given requirements there's no need to over-engineer the solution
