package main

import (
	"fmt"
	"math"
)

func doubleValue(n int) int {
	n *= 2
	return n
}

func doublePointer(n *int) int {
	*n *= 2
	return *n
}

func pointersVsValue(x int) {
	fmt.Println("--- Pointer vs Value Drill ---")

	doubleValue(x)
	fmt.Printf("After doubleValue(x): %d (No change to original)\n", x)

	doublePointer(&x)
	fmt.Printf("After doublePointer(&x): %d (Original mutated)\n", x)

	/*
	   WHY IT MATTERS:
	   1. Scope Control: Pass-by-value protects the original data from side effects.
	   2. Efficiency: Pointers avoid "stack" bloat. If 'x' were a 10MB struct,
	      passing by value would copy all 10MB into memory for the function.
	      Passing a pointer only copies 8 bytes (the address).
	*/
}

type Counter struct {
	count int
}

func (c Counter) Value() int  { return c.count }
func (c *Counter) Increment() { c.count++ }

func drill_02_whyPointerReceiver() {
	fmt.Println("--- DRILL 02: Pointer vs Value Receiver ---")

	c := Counter{}
	c.Increment()
	c.Increment()
	fmt.Println("Counter value:", c.Value())

	/*
	   POINTER VS VALUE RECEIVERS:
	   In web development, we often use pointer receivers (*) to mutate data
	   coming from a web server, such as in middleware or when processing
	   HTTP requests/responses.

	   MEMORY ALLOCATION:
	   Every time you construct this struct, you tell the computer to allocate
	   enough space for this group of data. Ensure data types are aligned
	   for memory efficiency.

	   type BadExample struct {
	       a bool  // 1 byte
	       b int64 // 8 bytes
	       c bool  // 1 byte
	   } // Total size: 24 bytes (due to padding)

	   type GoodExample struct {
	       b int64 // 8 bytes
	       a bool  // 1 byte
	       c bool  // 1 byte
	   } // Total size: 16 bytes (less padding)

	   STRUCTS VS CLASSES:
	   Structs are value types; when you copy a struct, you create a new,
	   independent copy of the data in memory. In contrast, classes (in other
	   languages) are reference types; when you copy a class, both variables
	   point to the same memory address.

	   QUIRK:
	   You can call pointer methods on addressable values (variables),
	   but you CANNOT call pointer methods on non-addressable values.
	   Example: Counter{}.Increment() will cause a compile error because
	   a literal is temporary and has no address.

	   INTERFACE TRAP:
	   If a struct has a pointer receiver method, then *Counter satisfies
	   the interface, but the Counter value does not.
	*/
}

// NIL POINTERS
// Calling a method on a nil pointer depends on whether you handle the nil case.
// You should only return a value if the receiver is not nil; otherwise, it is safe.
// If you try to access data on a nil pointer without handling it, the program will panic.
type Node struct {
	n    int
	next *Node
}

func (n *Node) SafeVal() int {
	if n == nil {
		return -1
	}
	return n.n
}

func (n *Node) unSafeVal() *Node {
	// If 'n' is nil, this will panic because it's trying to access the 'next' field
	// of a memory address that doesn't exist (0x0).
	if n == nil {
		return nil
	}
	return n.next
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

// INTERFACE INTERNALS:
// An interface value has two fields internally: type and value (often called 'tab' and 'ptr').
// When you assign a value to an interface, Go stores:
// 1. A pointer to the type information (method table).
// 2. A pointer to the actual data.
// An interface check validates BOTH TYPE AND VALUE.
type Logger interface {
	Log(msg string)
}

type MyLogger struct{}

// Satisfying the Logger interface
func (m *MyLogger) Log(msg string) {
	if m == nil {
		fmt.Println("Log skipped: Logger is nil")
		return
	}
	fmt.Println("Log:", msg)
}

func drill_equality() {
	fmt.Println("--- DRILL: Interface Equality ---")
	var a, b interface{}

	a = int(5)
	b = int64(5)

	fmt.Printf("Is a == b? %v\n", a == b)
	// Result: false. The values are both 5, but the types (int vs int64) differ.
}

func drill_pointer_equality() {
	fmt.Println("--- DRILL: Pointer Equality ---")
	var a, b any

	x, y := 10, 10

	a = &x
	b = &y

	fmt.Printf("Is a == b? %v\n", a == b)
	// Result: false
	// Why? Both are type *int, but the Pointers (the 'ptr' field)
	// point to different memory addresses in the stack/heap.
}

type Circle struct{ Radius float64 }
type Rect struct{ W, H float64 }

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func (r Rect) Area() float64        { return r.W * r.H }
func (r Rect) Perimeter() float64   { return 2 * (r.W + r.H) }

func PrintShapeInfo(s Shape) {
	fmt.Printf("Area : %.2f | Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func lobya() {
	fmt.Println("--- DRILL: Interface Slice (lobya) ---")
	shapes := []Shape{
		Circle{Radius: 5},
		Rect{W: 10, H: 5},
		Circle{Radius: 2},
	}
	for _, v := range shapes {
		PrintShapeInfo(v)
	}
}

func main() {
	fmt.Println("=== RUNNING ALL DRILLS ===")
	drill_02_whyPointerReceiver()

	fmt.Println("\n=== POINTER DRILLS ===")
	y := 42
	pointersVsValue(y)

	fmt.Println()
	drill_equality()
	drill_pointer_equality()

	fmt.Println()
	lobya()

	fmt.Println("\n=== NIL RECEIVER CHECK ===")
	var n *Node
	fmt.Println("SafeVal on nil Node:", n.SafeVal())
}
