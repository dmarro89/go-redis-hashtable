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
Benchmarking tests evaluate the performance of the Golang Redis-like hashtable implementation under different scenarios.

### Methodology
Benchmarks are conducted for setting, getting, and deleting key-value pairs, comparing the custom hashtable implementation with Go's native `map`.

### Results

| Benchmark                |   Num. Op. |   Time (ns/op) |   Mem (B/op) |   Mem.Op. |
|:-------------------------|-----------:|---------------:|-------------:|----------:|
| BenchmarkSet/1e1-4        |     563031 |           1948 |          896 |        17 |
| BenchmarkSet/1e2-4        |      60679 |          19587 |         7280 |       113 |
| BenchmarkSet/1e3-4        |       5817 |         194971 |        67312 |      1020 |
| BenchmarkGet/1e1-4        |    2811422 |           426.4|             0|         0 |
| BenchmarkGet/1e2-4        |     258082 |           4635 |             0|         0 |
| BenchmarkGet/1e3-4        |      24170 |          48600 |             0|         0 |
| BenchmarkDelete/1e1-4     |    1951491 |           660.7|             0|         0 |
| BenchmarkDelete/1e2-4     |     237284 |           5001 |             0|         0 |
| BenchmarkDelete/1e3-4     |      22318 |          53780 |             0|         0 |
| BenchmarkGoMapSet/1e1-4   |    1357090 |           884.9|           742|        11 |
| BenchmarkGoMapSet/1e2-4   |     103054 |          11401 |         11772|       109 |
| BenchmarkGoMapSet/1e3-4   |       7419 |         145960 |        182496|      1030 |
| BenchmarkGoMapGet/1e1-4   |    6300027 |           188.9|             0|         0 |
| BenchmarkGoMapGet/1e2-4   |     484138 |           2153 |             0|         0 |
| BenchmarkGoMapGet/1e3-4   |      38877 |          30688 |             0|         0 |
| BenchmarkGoMapDelete/1e1-4|    1498592 |           785.0|             0|         0 |
| BenchmarkGoMapDelete/1e2-4|     158318 |           7546 |             7|         0 |
| BenchmarkGoMapDelete/1e3-4|      15902 |          75495 |             0|         0 |


### Explanation

### Analysis

From the benchmark results (you can take a look also at the [actions](https://github.com/dmarro89/go-redis-hashtable/actions/workflows/go-main.yml)), we observe the following:

- **Insertion (`BenchmarkSet` vs `BenchmarkGoMapSet`)**: The Go Native `map` consistently outperforms the custom hashtable in terms of Time. For the memory consumption, the custom hash table performs better (so less memory consumed for operation) when the dataset dimension increases.
- **Retrieval (`BenchmarkGet` vs `BenchmarkGoMapGet`)**: The performance difference for retrievals is more moderate, but the Go Native `map` remains faster in terms of time (about the double of time) due to efficient key lookups.
- **Deletion (`BenchmarkDelete` vs `BenchmarkGoMapDelete`)**: Deletions in the custom hashtable are faster compared the Go Native `map` - both are consumin no memory for operation.

### Conclusion
While the Custom Hashtable provides a functional alternative to Go's native implementation, it is evident that the Go Native `map` in some cases are more efficient, particularly for large data sets.
This project is a work in progress, so there are constant iterations providing optimization for the data structure. 
Star the project and keep updated!


