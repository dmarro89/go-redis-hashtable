package main

import (
	"fmt"

	"github.com/dmarro89/go-redis-hashtable/structure"
)

func main() {
	database := structure.NewSipHashDict()

	database.Set("key1", "value1")
	database.Set("key2", "value2")

	value := database.Get("key1")
	fmt.Printf("key1: %s", value)

	value = database.Get("key2")
	fmt.Printf("key2: %s", value)

	database.Delete("key1")
	value = database.Get("key1")
	fmt.Printf("key1 after delete: %s", value)

	database.Set("key2", "value3")
	value = database.Get("key2")
	fmt.Printf("key2 after update: %s", value)
}
