package aa

// MyEnum is an example enumeration type
type MyEnum int

const (
	EnumValue1 MyEnum = iota // EnumValue1
	EnumValue2               // EnumValue2
	EnumValue3               // EnumValue3
)

// MyVar is an example variable
var MyVar = "Hello, World!" // MyVar

// SayHello prints a greeting message using MyVar.
func SayHello() {
	println(MyVar)
}
