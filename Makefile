

.PHONY: build gen

build: gen
	go build .

gen: country-codes.json
	go run ./gen

country-codes.json:
	curl -L https://pkgstore.datahub.io/core/country-codes/country-codes_json/data/616b1fb83cbfd4eb6d9e7d52924bb00a/country-codes_json.json -o country-codes.json
