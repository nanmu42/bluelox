package lox

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleLox_logics() {
	const code = `
print "hi" or 2; // "hi".
print nil or "yes"; // "yes".
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// hi
	// yes
}

func ExampleLox_binding_and_resolving() {
	const code = `
var a = "global";
{
  fun showA() {
    print a;
  }

  showA();
  var a = "block";
  showA();
}
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// global
	// global
}

func ExampleLox_closure() {
	const code = `fun makeCounter() {
  var i = 0;
  fun count() {
    i = i + 1;
    print i;
  }

  return count;
}

var counter = makeCounter();
counter(); // "1".
counter(); // "2".
counter(); // "3".
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleLox_fib() {
	const code = `
fun fib(n) {
  if (n <= 1) return n;
  return fib(n - 2) + fib(n - 1);
}

for (var i = 0; i < 20; i = i + 1) {
  print fib(i);
}
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// 0
	// 1
	// 1
	// 2
	// 3
	// 5
	// 8
	// 13
	// 21
	// 34
	// 55
	// 89
	// 144
	// 233
	// 377
	// 610
	// 987
	// 1597
	// 2584
	// 4181
}

func ExampleLox_if() {
	const code = `
if (true) {
print "me";
}
if (false) {
print "not me";
}
if (false) {
print "not me";
} else {
print "me again";
}
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// me
	// me again
}

func ExampleLox_for() {
	const code = `
var sum = 0;

for (var i = 1; i <= 100; i = i + 1) {
  sum = sum + i;
}

print sum;
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// 5050
}

func ExampleLox_while() {
	const code = `
var i = 1;
var sum = 0;

while (i <= 100) {
  sum = sum + i;
  i = i + 1;
}

print sum;
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// 5050
}

func ExampleLox_inheritance() {
	const code = `
class Doughnut {
  cook() {
    print "Fry until golden brown.";
  }
}

class BostonCream < Doughnut {
  cook() {
    super.cook();
    print "Pipe full of custard and coat with chocolate.";
  }
}

BostonCream().cook();
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// Fry until golden brown.
	// Pipe full of custard and coat with chocolate.
}

func ExampleLox_inheritance2() {
	const code = `
class A {
  method() {
    print "A method";
  }
}

class B < A {
  method() {
    print "B method";
  }

  test() {
    super.method();
  }
}

class C < B {}

C().test();
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// A method
}

func ExampleLox_resolving() {
	const code = `
class Cell {
    init(field) {
        // on or off
        this.s = false;

        // which field does this cell belongs to
        this.field = field;

        // neighbors, Cell
        this.up = nil;
        this.right = nil;
        this.down = nil;
        this.left = nil;
    }
}

class Field {
    // weight and height
    init(w, h) {
        this.w = w;
        this.h = h;

        // upper-left cell
        this.root = Cell(this);

        // weaving cells
        // Phase 1:
        // O ↔ O ↔ O
        // ↕
        // O ↔ O ↔ O
        // ↕
        // O ↔ O ↔ O
        var head = this.root;
        var tail = head;
        for (var col = 1; col < this.w; col = col+1) {
            var newTail = Cell(this);
            newTail.left = tail;
            tail.right = newTail;
            tail = newTail;
        }

        for (var row = 1; row < this.h; row = row+1) {
            var newHead = Cell(this);
            newHead.up = head;
            head.down = newHead;
            head = newHead;

            tail = head;
            for (var col = 1; col < this.w; col = col+1) {
                var newTail = Cell(this);
                newTail.left = tail;
                tail.right = newTail;
                tail = newTail;
            }
        }



        // Phase 2:
        // O - O - O
        // |   ↕   ↕
        // O - O - O
        // |   ↕   ↕
        // O - O - O
        var rowEnds = this.root;
        for (var row = 1; row < this.h; row = row+1) {
            var head = rowEnds;
            rowEnds = rowEnds.down;
            var tail = rowEnds;
			print this.root.down.right;
			print tail;
			print rowEnds;
			print rowEnds.right;
			print tail.right;
		}
	}
}

var l = Field(2, 2);
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// Cell instance
	// Cell instance
	// Cell instance
	// Cell instance
	// Cell instance
}

func ExampleLox_print_class() {
	const code = `
class DevonshireCream {
  serveOn() {
    return "Scones";
  }
}

print DevonshireCream; // Prints "DevonshireCream".

class Bagel {}
var bagel = Bagel();
print bagel; // Prints "Bagel instance".

class Bacon {
  eat() {
    print "Crunch crunch crunch!";
  }
}

Bacon().eat(); // Prints "Crunch crunch crunch!".

class Cake {
  taste() {
    var adjective = "delicious";
    print "The " + this.flavor + " cake is " + adjective + "!";
  }
}

var cake = Cake();
cake.flavor = "German chocolate";
cake.taste(); // Prints "The German chocolate cake is delicious!".

class Thing {
  getCallback() {
    fun localFunction() {
      print this;
    }

    return localFunction;
  }
}

var callback = Thing().getCallback();
callback();

class Foo {
  init() {
    print "foo initialized";
  }
}

var foo = Foo();
print foo.init();
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// DevonshireCream
	// Bagel instance
	// Crunch crunch crunch!
	// The German chocolate cake is delicious!
	// Thing instance
	// foo initialized
	// foo initialized
	// Foo instance
}

func Test_Lox_no_local_duplicated(t *testing.T) {
	const code = `
var a = "outer";
{
  var a = a;
}
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	require.Error(t, err)
}

func Test_Lox_no_invalid_super(t *testing.T) {
	const code = `
class Eclair {
  cook() {
    super.cook();
    print "Pipe full of crème pâtissière.";
  }
}
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	require.Error(t, err)
}

func Test_Lox_no_invalid_super2(t *testing.T) {
	const code = `
super.notEvenInAClass();
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	require.Error(t, err)
}

func Test_Lox_no_returning_from_init(t *testing.T) {
	const code = `
class Foo {
  init() {
    return "something else";
  }
}
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	require.Error(t, err)
}

func Test_Lox_invalid_use_of_this(t *testing.T) {
	const code = `
fun notAMethod() {
  print this;
}
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	require.Error(t, err)
}

func Test_Lox_no_top_level_return(t *testing.T) {
	const code = `
return "surprise!";
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	require.Error(t, err)
}

func Test_Lox_no_duplicated_declaring(t *testing.T) {
	const code = `
fun bad() {
  var a = "first";
  var a = "second";
}
`

	l := NewLox(os.Stdout)
	err := l.Run(context.TODO(), []byte(code))
	require.Error(t, err)
}
