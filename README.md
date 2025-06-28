# rossn

[![Go Reference](https://pkg.go.dev/badge/github.com/ac999/rossn.svg)](https://pkg.go.dev/github.com/ac999/rossn)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**rossn** is a production-ready, permissively-licensed Go package for validating Romanian CNP (Personal Numeric Code) numbers.  
It checks CNP structure, birth date, county, serial number, and official checksum.  
Thoroughly tested against all official rules and edge cases.

---

## Features

- Validates CNP length, digits, date, county, serial (001–999), and checksum
- Full support for archival Bucharest codes 47/48 (historical), and for code 70 for foreign/stateless residents (S=7,8,9 only), as per [cnp-spec](https://github.com/vimishor/cnp-spec) and [Wikipedia](https://ro.wikipedia.org/wiki/Cod_numeric_personal).
- Full support for SIIEASC JJ=70 (2024+) as per the [official statement of the Romanian Ministry of Internal Affairs, May 2024](https://www.mai.gov.ro/precizari-in-ceea-ce-priveste-generarea-codului-numeric-personal/)
- 100% Go, no dependencies
- MIT licensed: Free for commercial and closed-source use
- Fast, robust, and tested

## Install

```bash
go get github.com/ac999/rossn
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/ac99/rossn"
)

func main() {
    err := rossn.Validate("1981214320015")
    if err != nil {
        fmt.Println("Invalid CNP:", err)
    } else {
        fmt.Println("Valid CNP!")
    }
}
```

## CNP Specification

A CNP is 13 digits: `SYYMMDDJJNNNC`

- `S`: Gender and century
- `YYMMDD`: Date of birth (with S determining the century)
- `JJ`: County code (01–52, official list)
- `NNN`: Serial number (001–999)
- `C`: Control digit (checksum, see [here](https://ro.wikipedia.org/wiki/Cod_numeric_personal))

All checks are performed according to the official Romanian government rules.

## Testing

Run all tests with:

```bash
go test ./...
```

Or for detailed output:

```bash
go test -v
```

## Contributing

Pull requests and issue reports are welcome!
Before submitting, make sure all tests pass and that new edge cases are covered by tests.

---

## License

MIT License — see [LICENSE](LICENSE)

---

## References

- [Personal Identification Number Specification (cnp-spec)](https://github.com/vimishor/cnp-spec)
- [Wikipedia - Cod numeric personal](https://ro.wikipedia.org/wiki/Cod_numeric_personal
- [Precizări în ceea ce privește generarea codului numeric personal (mai 2024, MAI)](https://www.mai.gov.ro/precizari-in-ceea-ce-priveste-generarea-codului-numeric-personal/)
