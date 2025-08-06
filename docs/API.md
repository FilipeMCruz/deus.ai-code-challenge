# API

This document lists the available endpoints and how to call them.

API Versioning is defined directly in the URLs. 

All unsuccessful requests return the following body:

```json
{
  "error": string
}
```

The error status codes can be:
- 400: invalid request (missing json field or query param)
- 404: para not found (given url is not valid)
- 500: internal server failure

## Number of unique visitors for given page

URL: '/api/v1/unique-visitors'
Body: none
Headers: none
Query:

- pageUrl: string

Successful response:

Status Code: 200 (ok)
Body:

```json
{
  "unique_visitors": number
}
```

Example:

```shell
curl "http://localhost:8080/api/v1/unique-visitors?pageUrl=u"
```

Other Status Codes: 400, 404, 500

## Stats

URL: '/api/v1/user-navigation'
Body:

```json
{
  "visitor_id": string
  "page_url": string
}
```

Headers: none
Query: none

Successful response:

Status Code: 200 (ok)

Example:

```shell
echo '{"visitor_id":"b", "page_url":"u"}' | curl -X POST "http://localhost:8080/api/v1/user-navigation" --data-binary @-
```

Other Status Codes: 400, 404, 500
