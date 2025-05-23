module sql-example

go 1.21.0

toolchain go1.24.3

replace github.com/gtoxlili/sql-adapter => ../../.

require (
	github.com/gtoxlili/sql-adapter v0.0.0-00010101000000-000000000000
	github.com/casbin/casbin/v2 v2.105.0
	github.com/go-sql-driver/mysql v1.9.2
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/casbin/govaluate v1.3.0 // indirect
)
