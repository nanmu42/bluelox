package lox

func ExampleLox_logics() {
	const code = `
print "hi" or 2; // "hi".
print nil or "yes"; // "yes".
`

	l := NewLox()
	err := l.run([]byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// hi
	// yes
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

	l := NewLox()
	err := l.run([]byte(code))
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

	l := NewLox()
	err := l.run([]byte(code))
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

	l := NewLox()
	err := l.run([]byte(code))
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

	l := NewLox()
	err := l.run([]byte(code))
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

	l := NewLox()
	err := l.run([]byte(code))
	if err != nil {
		panic(err)
	}
	// Output:
	// 5050
}
