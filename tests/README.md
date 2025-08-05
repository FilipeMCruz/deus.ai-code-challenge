# Tests

This section presents an easy way to test that the solution is free from data races.
The tool used was [k6](https://k6.io/), the script present in [here](k6/script.js) helps to test the solution under
various basic load scenarios.

These tests could, after careful considerations, run in the CI pipeline and profile the application in terms
of performance, helping us detect future performance regressions.

## Running

To run the test one needs to install k6.
Once installed, ensure that the service is running and, from the root folder, run:

```shell
k6 run ./tests/k6/script.js
```

Tweak the scenario by changing the parameters/variables:

```shell
k6 run ./tests/k6/script.js --env DOMAIN=localhost:8080 --env N_VISITORS=15 --env N_PAGES=10 --env N_VISITS=20
```

## Scenario

Given a number of:

- unique visitors
- unique pages
- visits made per unique visitor

Generate fake data to flood the service, compute the expected results, send the data and verify if the service report
back unique page visits correctly.

Be aware that the generated page urls between runs can collide, resulting in test failures in the server isn't restarted
between runs.

## Parameters/Variables

| Parameter/Variable | Default        | Description                                                      |
|--------------------|----------------|------------------------------------------------------------------|
| DOMAIN             | localhost:8080 | domain where the solution is hosted                              |
| N_VISITORS         | 10             | number of unique visitors (concurrent users calling the service) |
| N_PAGES            | 10             | number of unique pages                                           |
| N_VISITS           | 10             | number of visits made by each visitor                            |
