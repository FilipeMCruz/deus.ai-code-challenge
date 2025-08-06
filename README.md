# Go Coding Challenge v1

This document presents the needed information to install and run the solution given for the Code Challenge. To better
understand the code there are some comments in it.

This solution is a simple http server written in go.

Since the domain and requirements are simple there's no need to over-complexify the solution, therefore I choose to:

- not use any libraries, go already provides the basics needed to do what this challenge requires;
- not follow a commonly used architecture such as onion, layered or clean architecture for the sake of simplicity:
    - e.g. there's no need to have a "service/use case" layer when there's almost no business logic, rules or processes
      to run.

API details can be found [here](docs/API.md).

Architecture details can be found [here](docs/ARCHITCTURE.md).

Basic integration test can be found [here](tests/README.md).

Branch [feat/channels](https://github.com/FilipeMCruz/deus.ai-code-challenge/tree/feat/channels) contains a different
approach to data synchronization in the repository package/layer, using channels instead of a mutex.

Branch [feat/valid-pages-only](https://github.com/FilipeMCruz/deus.ai-code-challenge/tree/feat/valid-pages-only) adds a
business requirement (set of valid page urls are specified at startup) to justify adding a "service" layer to the mix.

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
./server -port 8080
```

### Docker

To run the solution in port 8080:

```shell
docker compose up -d
```

## Requirements assumed

- page urls and visitor id can't be represented as an empty string;
- all non-empty visitor ids and page urls received are valid and our service can assume that they exist within the
  deus.ai domain;
- a page url can be seen as a unique identifier, e.g.: https://example.org/page?query=x != http://example.org/page;
- page url and visitor ids are case sensitive, e.g.: visitor alex != AleX

## Possible improvements

This section details some of the thoughts that crossed my mind while working on the solution.

### Readability/Maintainability

For me "Readability" is deeply tied to the standards followed by the company and knowledge shared within the team.
Therefore, I always feel like I'm a bit in the dark when tackling this aspect.

Maintainability is also very dependent on the roadmap for the service and it's hard to know how to improve it,
nonetheless, imagining the following sprints would introduce much richer business rules/logic:

- introducing a "service/use case" layer to decouple input and output validation (api) from the actual business rules
  could bring improvements, it would also follow the clean architecture, known by most developers;
- currently the business contains a single "bounded context". In the future if the domain becomes richer, splitting the
  server by "bounded contexts" while keeping the api, domain and repository layers will lead to a Modular Monolith
  architecture that is much easier to split by teams and break into smaller services;

### Performance

Without proper test data it's near impossible to know how performance can be improved, but:

- if the number of unique pages is really big, one can shard the in memory map into multiple smaller ones, allowing
  locks to be more granular and improving concurrency;
- memory and IO resource usage may become a problem, if so, horizontal scaling is prefer in the cloud era we're living
  in, this would require the following changes to architecture:
    - load balancer in front of the multiple instances of our service;
    - shared redis instance where the required data is stored and accessed by our services (this would also ensure no
      data is lost if a service is shutdown);

### API Access

As it is, this service is prone to data theft and falsification even if it is not exposed to the public domain.
Since most communication will be made on a machine to machine basis, I'd suggest using API keys that are automatically
rotated every couple of months. Those keys would be shared by this service and all its clients.

If needed the keys could have their set of scopes defined as such:

- admin (full access);
- read (able to query page statistics);
- write (able to register page visits).

The key would be sent on every request under the header 'Authorization' and validated in this service before preforming
the requested action.

### Security

Currently this service doesn't encrypt its connections, serving the data in plain text. This makes it very susceptible
to man-in-the-middle attacks. This attack vector is not a concern if the service isn't exposed to the public domain,
since the company can trust its employees, and can therefore be ignored.

If it is exposed, the endpoint need to be secured with digital SSL/TLS certificates.
I'd suggest setting up a reverse proxy (such as nginx) or ingress controller that would handle TLS/SSL termination,
instead of trying to load and cycle SSL/TLS certificates within this service.

I've used 'lets encrypt' and 'certbot' in the past to handle this requirement.

### Data Retention

All data is currently stored in memory. This is far from ideal since, in the case of a shutdown everything would be
lost. With services running more and more frequently in stateless environments, that can be restarted/destroyed at
anytime, such as lambda functions or kubernetes this is a major concern.
Depending on the requirements, one would need to store the information in disk. Databases are the right tool for this,
if the functional requirements are kept, I'd use a key value store such as cassandra or DynamoDB where the key is the
page url.

Since this seems to be a service where performance matters I'd write a PoC where a cache, such as redis, is placed
between the service and the database to take some load from it (if the visitors are like me they'll spend most of the
time looking at the same set of pages, there's no need to repeatedly query the database with the same data).

## Notes

- the business rules of this code challenge are almost none existing, in a proper application I'd love for most business
  logic to live in the domain and have a proper 'services' layer, but with the given requirements there's no need to
  over-engineer the solution;
- the repository package is an agnostic package responsible for handling common needs such as logging and gracefully
  shutting down the server;
- CORS wasn't a concern since the challenge seems to indicate that there's a system that sits between the browser and
  this
  service.

This service is far from "production" ready, there's a lot of interesting topics to discuss here:

- what are the authentication/authorization needs;
- what tools are used to document the api surface (e.g. swagger/OpenAPI, simple API.md);
- what linting rules are used;
- should we define the requirements for this service test coverage? e.g. 75% of the lines covered by tests;
- what does the infrastructure tied to it looks like? for things like observability/monitoring, distributed logging,
  orchestration, etc;
- accessibility of the service, e.g. will it be publicly exposed? If so, how are TLS/SSL certificates normally used
  within the company?
- where it would run (single container, serverless function in the cloud, on-prem in a VM...);
- how are integration tests written and maintain now that the this service will be integrated in the overall
  architecture?
- what environments, besides prod, are set up?
- what does the SDLC looks like? How do teams approach the development? features branches, main branch for prod,
  develop for dev/uat environments?

As a final note, please understand that most of what's written in this section are just random ramblings made by me
while writing the solution, here to hopefully drive some future discussions.
