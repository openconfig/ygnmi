// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ygnmi_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/openconfig/ygnmi/internal/exampleoc"
	"github.com/openconfig/ygnmi/internal/exampleoc/root"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
)

func ExampleLookup() {
	c := initClient()
	// Lookup a value at path.
	path := root.New().Parent().Child()
	val, err := ygnmi.Lookup(context.Background(), c, path.State())
	if err != nil {
		log.Fatalf("failed to lookup: %v", err)
	}

	// Check if the value is present.
	child, ok := val.Val()
	if !ok {
		log.Fatal("no value at path")
	}
	fmt.Printf("Child: %v\n, RecvTimestamo %v", child, val.RecvTimestamp)
}

func ExampleGet() {
	c := initClient()
	// Get a value at path.
	// No metadata is returned, and an error is returned if not present.
	path := root.New().Parent().Child()
	child, err := ygnmi.Get(context.Background(), c, path.State())
	if err != nil {
		log.Fatalf("failed to lookup: %v", err)
	}

	fmt.Printf("Child: %v\n", child)
}

func ExampleWatch() {
	c := initClient()
	// Watch a path
	path := root.New().Parent().Child()
	// Use a context with a timeout, so that Watch doesn't continue indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	watcher := ygnmi.Watch(ctx, c, path.State(), func(v *ygnmi.Value[*exampleoc.Parent_Child]) error {
		if c, ok := v.Val(); ok {
			if c.GetOne() == "good" {
				// Return nil to indicate success and stop evaluating predicate.
				return nil
			} else if c.GetTwo() == "bad" {
				// Return a non-nil to indicate failure and stop evaluation predicate.
				// This error is returned by Await().
				return fmt.Errorf("unexpected value")
			}
		}

		// Return ygnmi.Continue to continue evaluating the func.
		return ygnmi.Continue
	})
	// Child is the last value received.
	child, err := watcher.Await()
	if err != nil {
		log.Fatalf("failed to lookup: %v", err)
	}
	// Check if the value is present.
	val, ok := child.Val()
	if !ok {
		log.Fatal("no value at path")
	}
	fmt.Printf("Last Value: %v\n", val)
}

func ExampleAwait() {
	c := initClient()
	path := root.New().Parent().Child()

	// Use a context with a timeout, so that Await doesn't continue indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	want := &exampleoc.Parent_Child{
		Two: ygot.String("good"),
	}

	// Wait until the values at the path is equal to want.
	child, err := ygnmi.Await(ctx, c, path.State(), want)
	if err != nil {
		log.Fatalf("failed to lookup: %v", err)
	}
	// Check if the value is present.
	val, ok := child.Val()
	if !ok {
		log.Fatal("no value at path")
	}
	fmt.Printf("Last Value: %v\n", val)
}

func ExampleCollect() {
	c := initClient()
	path := root.New().Parent().Child()

	// Use a context with a timeout, so that Collect doesn't continue indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Collect values at path until context deadline is reached.
	collector := ygnmi.Collect(ctx, c, path.State())
	vals, err := collector.Await()
	if err != nil {
		log.Fatalf("failed to lookup: %v", err)
	}

	fmt.Printf("Get %d values", len(vals))
}

func initClient() *ygnmi.Client {
	return nil
}
