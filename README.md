# Go-Redis-Hashtable

**Golang Redis Hashtable Implementation**

## Overview

This project provides a Golang implementation of a hashtable, inspired by the hashtable used in Redis. Redis employs a hash table to store its key-value pairs efficiently in-memory.

## Table of Contents

1. [Introduction](#introduction)
2. [Hashtable in Redis](#hashtable-in-redis)
3. [Golang Implementation](#golang-implementation)
   - [Dict Struct](#dict-struct)
   - [Operations](#operations)
4. [Usage](#usage)
   - [Getting Started](#getting-started)
   - [Example](#example)
5. [Benchmarking](#benchmarking)
   - [Introduction](#introduction-1)
   - [Methodology](#methodology)
   - [Results](#results)
      - [Insertion Benchmark](#insertion-benchmark)
      - [Retrieval Benchmark](#retrieval-benchmark)
      - [Deletion Benchmark](#deletion-benchmark)
6. [Contributing](#contributing)
7. [License](#license)
## Introduction

### Project Objective

This project implements a Redis-like hashtable in Golang, providing a high-performance and reliable data structure for key-value storage.

## Hashtable in Redis

Redis, an in-memory data structure store, uses a hashtable as its main indexing mechanism. The hashtable is responsible for mapping keys to values efficiently.

### Redis Hashtable

#### Structure and Characteristics

- **Separate Chaining:** Redis uses a form of separate chaining to handle collisions. In this approach, each bucket in the hashtable contains a linked list of key-value pairs that hash to the same index. If multiple keys collide, they are stored in the same bucket as a linked list.

- **Dynamic Resizing:** Redis employs dynamic resizing to adapt to the number of elements it contains. When the hashtable reaches a certain load factor, Redis increases its size and rehashes the existing elements to maintain a balanced distribution and optimize performance.

- **Rehashing:** During resizing, a new larger hashtable is created, and all existing elements are rehashed and redistributed. This process ensures that the new hashtable has enough space to accommodate the growing number of elements.

- **Hashing Algorithm:** Redis uses a reliable hashing algorithm to distribute keys uniformly across the hashtable. This minimizes the likelihood of collisions, but when they occur, the linked list structure efficiently handles them.

#### Link
Here's some interesting links to read about Redis hashmap.

- https://codeburst.io/a-closer-look-at-redis-dictionary-implementation-internals-3fd815aae535
- https://blog.wjin.org/posts/redis-internal-data-structure-dictionary.html

## Golang Implementation

### Dict Struct

The `Dict` struct is the main struct of the project, named after the Redis data structure. It closely follows the Redis model with key components:

- `HashTable[2]`: The `Dict` struct contains two hash tables: a main table (index 0) and a rehashing table (index 1). Usually, items are stored in the main hashtable, and the rehashing one is used during the expanding and rehashing process.

- `DictEntry`: Each hashtable has a linked list of `DictEntry`. Each `DictEntry` has a key, a value, and a pointer to the next element.

### Operations

The `Set` operation is responsible for adding a key-value pair to the hashtable. It determines the index using the SipHash algorithm and stores the element in the main table at the specified index. If multiple key/value pairs share the same index, the hashtable utilizes a linked list to manage collisions.

The `Get` operation retrieves an entry based on the provided key. It considers potential collisions and rehashing if necessary, returning the corresponding value.

The `Delete` operation removes an entry from the hashtable based on the specified key, maintaining the integrity of the hashtable structure.

Initially, each hashtable starts with a small size (4). Upon exceeding this size, the main hashtable undergoes expansion. The expansion involves using a rehashing table, which is twice the size of the mainTable. Linked lists are transferred to the expanded table during this process. Once migration is complete, the rehashing table becomes the main table, and the rehashing table is reset as an empty one

## Usage

### Getting Started

To use this hashtable implementation, follow these steps:

1. Import the `datastructures` package.
2. Create a new `Dict` instance using `NewDict()`.
3. Use the provided methods for adding, retrieving, and deleting key-value pairs.

### Example

```go
package main

import (
	"fmt"
	"go_db/datastructures"
)

func main() {
	// Create a new hashtable
	myDict := datastructures.NewDict()

	// Add key-value pairs
	err := myDict.Set("key1", "value1")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Retrieve a value
	result := myDict.Get("key1")
	fmt.Println("Value for key1:", result)

	// Delete an entry
	err = myDict.Delete("key1")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

## Benchmarking
### Introduction
Benchmarking tests were conducted to evaluate the performance of the Golang Redis-like hashtable implementation in various scenarios.

The benchmark results provide insights into the execution time and resource utilization of the hashtable.
The BenchmarkSet,BenchmarkGet and BenchmarkDelete are referring to the hashtable implementation of the project while the GoMap ones are referring to the native golang implementation.

| Benchmar                   | Num. Op.       |   Time           |     Mem         |     Mem.Op. |
| -------------------------- |:--------------:|:----------------:|:---------------:|:-----------:|
| BenchmarkSet-4       	 	 | 582559	      | 2077 ns/op	     | 321 B/op	       | 4 allocs/op |
| BenchmarkGet-4       	 	 | 525124	      | 1984 ns/op	     | 206 B/op	       | 2 allocs/op |
| BenchmarkDelete-4    	 	 | 1198183	      | 1695 ns/op	     | 173 B/op	       | 2 allocs/op |
| BenchmarkGoMapSet-4        | 684192	      | 2262 ns/op	     | 132 B/op	       | 1 allocs/op |
| BenchmarkGoMapGet-4        | 2281518	      | 574.4 ns/op	     | 0 B/op	       | 0 allocs/op |
| BenchmarkGoMapDelete-4     | 2994595	      | 595.4 ns/op	     | 0 B/op	       | 0 allocs/op |


