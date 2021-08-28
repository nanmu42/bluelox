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
