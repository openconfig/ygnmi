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
		log.Fatalf("Failed to Lookup: %v", err)
	}

	// Check if the value is present.
	child, ok := val.Val()
	if !ok {
		log.Fatal("No value at path")
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
		log.Fatalf("Failed to Get: %v", err)
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
		log.Fatalf("Failed to Watch: %v", err)
	}
	// Check if the value is present.
	val, ok := child.Val()
	if !ok {
		log.Fatal("No value at path")
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
		log.Fatalf("Failed to Await: %v", err)
	}
	// Check if the value is present.
	val, ok := child.Val()
	if !ok {
		log.Fatal("No value at path")
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
		log.Fatalf("Failed to Collect: %v", err)
	}

	fmt.Printf("Got %d values", len(vals))
}

func ExampleLookupAll() {
	c := initClient()
	// LookupAll on the keys of the SingleKey list
	path := root.New().Model().SingleKeyAny().Value()
	vals, err := ygnmi.LookupAll(context.Background(), c, path.State())
	if err != nil {
		log.Fatalf("Failed to Lookup: %v", err)
	}

	// Get the value of each list element.
	for _, val := range vals {
		if listVal, ok := val.Val(); ok {
			fmt.Printf("Got List Val %v at path %v", listVal, val.Path)
		}
	}
}

func ExampleGetAll() {
	c := initClient()
	// Get on the keys of the SingleKey list
	path := root.New().Model().SingleKeyAny().Value()
	vals, err := ygnmi.GetAll(context.Background(), c, path.State())
	if err != nil {
		log.Fatalf("Failed to Lookup: %v", err)
	}

	// Get the value of each list element.
	// With GetAll it is impossible to know the path associated with a value,
	// so use LookupAll or Batch with with wildcard path instead.
	for _, val := range vals {
		fmt.Printf("Got List Val %v", val)
	}
}

func ExampleWatchAll() {
	c := initClient()
	// Watch a path
	path := root.New().Model().SingleKeyAny().Value()
	// Use a context with a timeout, so that Watch doesn't continue indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var gotKeyFoo, gotKeyBar bool
	// Watch list values and wait for 2 different list elements to be equal to 10.
	// Set ExampleWatch_batch for another way to use Watch will a wildcard query.
	watcher := ygnmi.WatchAll(ctx, c, path.State(), func(v *ygnmi.Value[int64]) error {
		if val, ok := v.Val(); ok {
			gotKeyFoo = gotKeyFoo || v.Path.GetElem()[2].Key["key"] == "foo" && val == 10
			gotKeyBar = gotKeyBar || v.Path.GetElem()[2].Key["key"] == "bar" && val == 10
		}

		// Return nil to indicate success and stop evaluating predicate.
		if gotKeyFoo && gotKeyBar {
			return nil
		}
		// Return ygnmi.Continue to continue evaluating the func.
		return ygnmi.Continue
	})
	// Child is the last value received.
	child, err := watcher.Await()
	if err != nil {
		log.Fatalf("Failed to Watch: %v", err)
	}
	// Check if the value is present.
	val, ok := child.Val()
	if !ok {
		log.Fatal("No value at path")
	}
	fmt.Printf("Last Value: %v\n", val)
}

func ExampleCollectAll() {
	c := initClient()
	path := root.New().Model().SingleKeyAny().Value()

	// Use a context with a timeout, so that Collect doesn't continue indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Collect values at path until context deadline is reached.
	collector := ygnmi.CollectAll(ctx, c, path.State())
	vals, err := collector.Await()
	if err != nil {
		log.Fatalf("Failed to Collect: %v", err)
	}

	fmt.Printf("Got %d values", len(vals))
}

func ExampleGet_batch() {
	c := initClient()
	// Create a new batch object and add paths to it.
	b := new(root.Batch)
	// Note: the AddPaths accepts paths not queries.
	b.AddPaths(root.New().RemoteContainer().ALeaf(), root.New().Model().SingleKeyAny().Value())
	// Get the values of the paths in the paths, a Batch Query always returns the root object.
	val, err := ygnmi.Get(context.Background(), c, b.State())
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}
	fmt.Printf("Got aLeaf %v, SingleKey[foo]: %v", val.GetRemoteContainer().GetALeaf(), val.GetModel().GetSingleKey("foo").GetValue())

}

func ExampleWatch_batch() {
	c := initClient()
	// Create a new batch object and add paths to it.
	b := new(root.Batch)
	// Note: the AddPaths accepts paths not queries.
	b.AddPaths(root.New().RemoteContainer().ALeaf(), root.New().Model().SingleKeyAny().Value())
	// Watch all input path until they meet the desired condition.
	_, err := ygnmi.Watch(context.Background(), c, b.State(), func(v *ygnmi.Value[*exampleoc.Root]) error {
		val, ok := v.Val()
		if !ok {
			return ygnmi.Continue
		}
		if val.GetModel().GetSingleKey("foo").GetValue() == 1 && val.GetModel().GetSingleKey("bar").GetValue() == 2 &&
			val.GetRemoteContainer().GetALeaf() == "string" {

			return nil
		}
		return ygnmi.Continue
	}).Await()
	if err != nil {
		log.Fatalf("Failed to Watch: %v", err)
	}
}

func ExampleUpdate() {
	c := initClient()
	// Perform the Update request.
	res, err := ygnmi.Update(context.Background(), c, root.New().Parent().Child().Three().Config(), exampleoc.Child_Three_TWO)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	fmt.Printf("Timestamp: %v", res.Timestamp)
}

func ExampleReplace() {
	c := initClient()
	// Perform the Replace request.
	res, err := ygnmi.Replace(context.Background(), c, root.New().RemoteContainer().Config(), &exampleoc.RemoteContainer{ALeaf: ygot.String("foo")})
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	fmt.Printf("Timestamp: %v", res.Timestamp)
}

func ExampleDelete() {
	c := initClient()
	// Perform the Update request.
	res, err := ygnmi.Delete(context.Background(), c, root.New().Parent().Child().One().Config())
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	fmt.Printf("Timestamp: %v", res.Timestamp)
}

func ExampleSetBatch_Set() {
	c := initClient()
	b := new(ygnmi.SetBatch)

	// Add set operations to a batch request.
	ygnmi.BatchUpdate(b, root.New().Parent().Child().Three().Config(), exampleoc.Child_Three_TWO)
	ygnmi.BatchDelete(b, root.New().Parent().Child().One().Config())
	ygnmi.BatchReplace(b, root.New().RemoteContainer().Config(), &exampleoc.RemoteContainer{ALeaf: ygot.String("foo")})

	// Perform the gnmi Set request.
	res, err := b.Set(context.Background(), c)
	if err != nil {
		log.Fatalf("Set failed: %v", err)
	}
	fmt.Printf("Timestamp: %v", res.Timestamp)

	root.New().Model().SingleKeyAny().Config()
}

func initClient() *ygnmi.Client {
	// Initialize the client here.
	return nil
}
