#!/bin/sh

exec > overflow_impl.go

echo "package overflow

// This is generated code, created by overflow_template.sh executed
// by \"go generate\"

import \"math\"

"


for SIZE in 8 16 32 64
do
echo "
// Add${SIZE} performs + operation on two int${SIZE} operands
// returning a result and status
func Add${SIZE}(a, b int${SIZE}) (int${SIZE}, bool) {
        c := a + b
        if (c > a) == (b > 0) {
                return c, true
        }
        return c, false
}

// Add${SIZE}p is the unchecked panicking version of Add${SIZE}
func Add${SIZE}p(a, b int${SIZE}) int${SIZE} {
        r, ok := Add${SIZE}(a, b)
        if !ok {
                panic(\"addition overflow\")
        }
        return r
}

// Sub${SIZE} performs - operation on two int${SIZE} operands
// returning a result and status
func Sub${SIZE}(a, b int${SIZE}) (int${SIZE}, bool) {
        c := a - b
        if (c < a) == (b > 0) {
                return c, true
        }
        return c, false
}

// Sub${SIZE}p is the unchecked panicking version of Sub${SIZE}
func Sub${SIZE}p(a, b int${SIZE}) int${SIZE} {
        r, ok := Sub${SIZE}(a, b)
        if !ok {
                panic(\"subtraction overflow\")
        }
        return r
}

// Mul${SIZE} performs * operation on two int${SIZE} operands
// returning a result and status
func Mul${SIZE}(a, b int${SIZE}) (int${SIZE}, bool) {
        if a == 0 || b == 0 {
                return 0, true
        }
        c := a * b
        if (c < 0) == ((a < 0) != (b < 0)) {
                if c/b == a {
                        return c, true
                }
        }
        return c, false
}

// Mul${SIZE}p is the unchecked panicking version of Mul${SIZE}
func Mul${SIZE}p(a, b int${SIZE}) int${SIZE} {
        r, ok := Mul${SIZE}(a, b)
        if !ok {
                panic(\"multiplication overflow\")
        }
        return r
}

// Div${SIZE} performs / operation on two int${SIZE} operands
// returning a result and status
func Div${SIZE}(a, b int${SIZE}) (int${SIZE}, bool) {
        q, _, ok := Quotient${SIZE}(a, b)
        return q, ok
}

// Div${SIZE}p is the unchecked panicking version of Div${SIZE}
func Div${SIZE}p(a, b int${SIZE}) int${SIZE} {
        r, ok := Div${SIZE}(a, b)
        if !ok {
                panic(\"division failure\")
        }
        return r
}

// Quotient${SIZE} performs + operation on two int${SIZE} operands
// returning a quotient, a remainder and status
//
// Deprecated: Quotient${SIZE} is deprecated, use Quo${SIZE} instead
func Quotient${SIZE}(a, b int${SIZE}) (int${SIZE}, int${SIZE}, bool) {
        if b == 0 {
                return 0, 0, false
        }
        c := a / b
        status := (c < 0) == ((a < 0) != (b < 0))
        return c, a % b, status
}

// Quo${SIZE} performs + operation on two int${SIZE} operands
// returning a quotient, a remainder and status
func Quo${SIZE}(a, b int${SIZE}) (int${SIZE}, int${SIZE}, bool) {
        if b == 0 {
                return 0, 0, false
        } else if b == -1 && a == math.MinInt${SIZE} {
                return 0, 0, false
        }
        c := a / b
        return c, a % b, true
}
"

# generate unsigned integer operations
echo "
// UAdd${SIZE} performs + operation on two uint${SIZE} operands
// returning a result and status
func UAdd${SIZE}(a, b uint${SIZE}) (uint${SIZE}, bool) {
        c := a + b
        if c >= a {
                return c, true
        }
        return c, false
}

// UAdd${SIZE}p is the unchecked panicking version of UAdd${SIZE}
func UAdd${SIZE}p(a, b uint${SIZE}) uint${SIZE} {
        r, ok := UAdd${SIZE}(a, b)
        if !ok {
                panic(\"addition overflow\")
        }
        return r
}

// USub${SIZE} performs - operation on two uint${SIZE} operands
// returning a result and status
func USub${SIZE}(a, b uint${SIZE}) (uint${SIZE}, bool) {
        if b <= a {
                return a - b, true
        }
        return 0, false
}

// USub${SIZE}p is the unchecked panicking version of USub${SIZE}
func USub${SIZE}p(a, b uint${SIZE}) uint${SIZE} {
        r, ok := USub${SIZE}(a, b)
        if !ok {
                panic(\"subtraction overflow\")
        }
        return r
}

// UMul${SIZE} performs * operation on two uint${SIZE} operands
// returning a result and status
func UMul${SIZE}(a, b uint${SIZE}) (uint${SIZE}, bool) {
        if a == 0 || b == 0 {
                return 0, true
        }
        c := a * b
        if (c / b) == a {
                return c, true
        }
        return c, false
}

// UMul${SIZE}p is the unchecked panicking version of UMul${SIZE}
func UMul${SIZE}p(a, b uint${SIZE}) uint${SIZE} {
        r, ok := UMul${SIZE}(a, b)
        if !ok {
                panic(\"multiplication overflow\")
        }
        return r
}

// UDiv${SIZE} performs / operation on two uint${SIZE} operands
// returning a result and status
func UDiv${SIZE}(a, b uint${SIZE}) (uint${SIZE}, bool) {
        q, _, ok := UQuotient${SIZE}(a, b)
        return q, ok
}

// UDiv${SIZE}p is the unchecked panicking version of UDiv${SIZE}
func UDiv${SIZE}p(a, b uint${SIZE}) uint${SIZE} {
        r, ok := UDiv${SIZE}(a, b)
        if !ok {
                panic(\"division failure\")
        }
        return r
}

// UQuotient${SIZE} performs / operation on two uint${SIZE} operands
// returning a quotient, a remainder and status
func UQuotient${SIZE}(a, b uint${SIZE}) (uint${SIZE}, uint${SIZE}, bool) {
        if b == 0 {
                return 0, 0, false
        }
        return a / b, a % b, true
}
"
done
