# Go Tax Calculator – API Docs

## Project

An HTTP API that follow REST style that returns Canadian **federal** income tax for a given **year** and **salary** using **marginal brackets**.

## Stack

* Go, `net/http`, `chi`
* Modules: `internal/api` (handlers), `internal/core` (Calculator), `internal/client` (tax Provider), `internal/models`, `internal/utils` (math utility functions)
* Makefile: `build`, `run`, `test`, `fmt`, `vet`

## Design

* Thin HTTP → pure `core.Calculator`
* `Provider` interface for tax brackets
* Money in **cents (int64)**; to prevent floating point arithmetic impresisions.
* Versioned routes: `/v1/...`

## Routes

* `GET /health` → `200 OK`
* `GET /v1/income-tax?year=[YEAR]&salary=[SALARY]`


| Name     | Type             | Required | Example    | Notes                                           |
| -------- | ---------------- | :------: | ---------- | ----------------------------------------------- |
| `year`   | integer          |     ✅    | `2019`     | Supported tax table year                        |
| `salary` | number or string |     ✅    | `50000.25` | Annual income; decimals allowed (dollars.cents) |

Example

`http://localhost:8080/v1/income-tax?year=2019&salary=50000.25`


  **200**:

```json
{
  "year": 2019,
  "salary_input": "50000.25",
  "salary_cents": 5000025,
  "total_tax_cents": 763040,
  "effective_rate": 0.1526,
  "per_bracket": [
    {
      "min_cents": 0,
      "max_cents": 4763000,
      "rate_basis_points": 1500,
      "tax_cents": 714450
    },
    {
      "min_cents": 4763000,
      "max_cents": 9525900,
      "rate_basis_points": 2050,
      "tax_cents": 48590
    }
  ],
  "message": "OK"
}

```

  **400/500**: `{"error":"..."}`

## Example

```bash
curl -G 'http://localhost:8080/v1/income-tax' \
  --data-urlencode year=2022 \
  --data-urlencode salary=85000.25
```

## Run & Test

```bash
make build && make run
make test
```
