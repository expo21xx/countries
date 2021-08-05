package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	data, err := loadData()
	if err != nil {
		return err
	}

	countries := jen.NewFile("countries")

	var cases []jen.Code

	for _, c := range data {
		if c.Name == "" {
			continue
		}

		name := normalizeName(c.Name)
		countries.Comment(name + " country data.")
		countries.Var().Id(name).Op("=").Id("Country").Values(jen.Dict{
			jen.Id("Name"):       jen.Lit(c.Name),
			jen.Id("Alpha2Code"): jen.Lit(c.Alpha2Code),
			jen.Id("Alpha3Code"): jen.Lit(c.Alpha3Code),
		})

		cases = append(cases, jen.Case(
			jen.Id(name).Dot("Name"),
			jen.Id(name).Dot("Alpha2Code"),
			jen.Id(name).Dot("Alpha3Code"),
		).Block(
			jen.Return(jen.Op("&").Id(name)),
		))
	}

	countries.Comment("ByNameOrCode returns the country by name, 3 or 2 letter code.")
	countries.Func().Id("ByNameOrCode").Params(jen.Id("c").Id("string")).Op("*").Id("Country").Block(
		jen.Switch(jen.Qual("strings", "ToUpper").Call(jen.Id("c"))).Block(
			cases...,
		),
		jen.Return(jen.Nil()),
	)

	return countries.Save("countries.go")
}

type country struct {
	Name       string `json:"CLDR display name"`
	Alpha2Code string `json:"ISO3166-1-Alpha-2"`
	Alpha3Code string `json:"ISO3166-1-Alpha-3"`
}

const dataJSONFile = "./country-codes.json"

func loadData() ([]country, error) {
	c, err := ioutil.ReadFile(dataJSONFile)
	if err != nil {
		return nil, err
	}

	var d []country

	err = json.Unmarshal(c, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func normalizeName(n string) string {
	re := strings.NewReplacer(
		" ", "",
		"&", "And",
		"(", "",
		")", "",
		"-", "",
		"’", "",
		".", "",
	)

	n = re.Replace(n)

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, n)
	return result
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
