package util

import "fmt"

// Stack represents a generic LIFO (Last-In, First-Out) stack.
// It uses a slice to store elements of type T.
type Stack[T any] struct {
	items []T
}

// NewStack creates and returns a pointer to a new, empty Stack for elements of type T.
// It initializes the internal slice.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: []T{}}
	// Alternative: return &Stack[T]{} // nil slice works fine with append
}

// Push adds an item of type T to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the item from the top of the stack.
// It returns the item and true if the stack was not empty.
// If the stack is empty, it returns the zero value of T and false.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T // Create the zero value for type T
		return zero, false
	}

	// Get the index of the top item
	index := len(s.items) - 1
	// Get the item
	item := s.items[index]
	// Remove the item by slicing (doesn't keep reference, allows GC)
	s.items = s.items[:index]

	return item, true
}

// Peek returns the item from the top of the stack without removing it.
// It returns the item and true if the stack was not empty.
// If the stack is empty, it returns the zero value of T and false.
// (Optional, but often useful)
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	index := len(s.items) - 1
	item := s.items[index]
	return item, true
}

// Size returns the number of items currently in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// IsEmpty returns true if the stack contains no items, false otherwise.
// (Optional helper)
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all elements from the stack
func (s *Stack[T]) Clear() {
	s.items = make([]T, 0)
}

func (s *Stack[T]) Debug() {
	for i, item := range s.items {
		fmt.Printf("Stack[%d] = %v\n", i, item)
	}
}

// --- Example Usage (can be in a separate main.go file) ---

/*
package main

import (
	"fmt"
	"your_module_path/stack" // Import the stack package
)

func main() {
	// --- Integer Stack ---
	fmt.Println("--- Integer Stack ---")
	intStack := stack.New[int]()

	fmt.Printf("Initial size: %d, IsEmpty: %t\n", intStack.Size(), intStack.IsEmpty())

	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)

	fmt.Printf("Size after pushes: %d\n", intStack.Size())

	peekVal, ok := intStack.Peek()
	if ok {
		fmt.Printf("Peeked: %d\n", peekVal)
	}
    fmt.Printf("Size after peek: %d\n", intStack.Size())


	val, ok := intStack.Pop()
	if ok {
		fmt.Printf("Popped: %d\n", val)
	}

	val, ok = intStack.Pop()
	if ok {
		fmt.Printf("Popped: %d\n", val)
	}

	fmt.Printf("Size after pops: %d\n", intStack.Size())

	intStack.Push(40)
	fmt.Printf("Size after push: %d\n", intStack.Size())

	// Pop remaining items
	for !intStack.IsEmpty() {
		val, ok := intStack.Pop()
		if ok {
			fmt.Printf("Popped: %d\n", val)
		}
	}

	fmt.Printf("Final size: %d, IsEmpty: %t\n", intStack.Size(), intStack.IsEmpty())

	// Try popping from empty stack
	val, ok = intStack.Pop()
	if !ok {
		fmt.Printf("Pop from empty stack failed as expected. Value: %d\n", val) // val will be 0 (zero value for int)
	}


	// --- String Stack ---
	fmt.Println("\n--- String Stack ---")
	stringStack := stack.New[string]()
	stringStack.Push("hello")
	stringStack.Push("world")
	stringStack.Push("generics")

	fmt.Printf("String stack size: %d\n", stringStack.Size())

	strVal, ok := stringStack.Pop()
	if ok {
		fmt.Printf("Popped string: %s\n", strVal)
	}
	strVal, ok = stringStack.Pop()
	if ok {
		fmt.Printf("Popped string: %s\n", strVal)
	}
	fmt.Printf("String stack size: %d\n", stringStack.Size())
}
*/
