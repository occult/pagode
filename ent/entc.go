//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/occult/pagode/ent/admin"
)

func main() {
	err := entc.Generate("./schema",
		&gen.Config{},
		entc.Extensions(&admin.Extension{}),
	)
	if err != nil {
		log.Fatal("running ent codegen:", err)
	}
}
