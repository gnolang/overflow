package overflow

import (
	"math"
	"testing"
)

// sample all possibilities of 8 bit numbers
// by checking against 64 bit numbers

func TestAlgorithms(t *testing.T) {

	errors := 0

	for a64 := int64(math.MinInt8); a64 <= int64(math.MaxInt8); a64++ {

		for b64 := int64(math.MinInt8); b64 <= int64(math.MaxInt8) && errors < 10; b64++ {

			a8 := int8(a64)
			b8 := int8(b64)

			if int64(a8) != a64 || int64(b8) != b64 {
				t.Fatal("LOGIC FAILURE IN TEST")
			}

			// ADDITION
			{
				r64 := a64 + b64

				// now the verification
				result, ok := Add8(a8, b8)
				if ok && int64(result) != r64 {
					t.Errorf("failed to fail on %v + %v = %v instead of %v\n",
						a8, b8, result, r64)
					errors++
				}
				if !ok && int64(result) == r64 {
					t.Fail()
					errors++
				}
			}

			// SUBTRACTION
			{
				r64 := a64 - b64

				// now the verification
				result, ok := Sub8(a8, b8)
				if ok && int64(result) != r64 {
					t.Errorf("failed to fail on %v - %v = %v instead of %v\n",
						a8, b8, result, r64)
				}
				if !ok && int64(result) == r64 {
					t.Fail()
					errors++
				}
			}

			// MULTIPLICATION
			{
				r64 := a64 * b64

				// now the verification
				result, ok := Mul8(a8, b8)
				if ok && int64(result) != r64 {
					t.Errorf("failed to fail on %v * %v = %v instead of %v\n",
						a8, b8, result, r64)
					errors++
				}
				if !ok && int64(result) == r64 {
					t.Fail()
					errors++
				}
			}

			// DIVISION
			if b8 != 0 {
				r64 := a64 / b64

				// now the verification
				result, _, ok := Quotient8(a8, b8)
				if ok && int64(result) != r64 {
					t.Errorf("failed to fail on %v / %v = %v instead of %v\n",
						a8, b8, result, r64)
					errors++
				}
				if !ok && result != 0 && int64(result) == r64 {
					t.Fail()
					errors++
				}
			}
		}
	}

}

func TestQuo8(t *testing.T) {
	tests := []struct {
		name     string
		a        int8
		b        int8
		wantQuot int8
		wantRem  int8
		wantBool bool
	}{
		{"simple division", 10, 3, 3, 1, true},
		{"exact division", 12, 4, 3, 0, true},
		{"zero dividend", 0, 5, 0, 0, true},
		{"one divisor", 42, 1, 42, 0, true},

		{"positive / positive", 7, 3, 2, 1, true},
		{"positive / negative", 7, -3, -2, 1, true},
		{"negative / positive", -7, 3, -2, -1, true},
		{"negative / negative", -7, -3, 2, -1, true},

		{"max value / 1", math.MaxInt8, 1, math.MaxInt8, 0, true},
		{"min value / 1", math.MinInt8, 1, math.MinInt8, 0, true},
		{"max value / 2", math.MaxInt8, 2, 63, 1, true},
		{"min value / 2", math.MinInt8, 2, -64, 0, true},

		{"min value / -1", math.MinInt8, -1, 0, 0, false},
		{"min value / min value", math.MinInt8, math.MinInt8, 1, 0, true},
		{"-1 / min value", -1, math.MinInt8, 0, -1, true},

		{"division by zero", 42, 0, 0, 0, false},
		{"zero division by zero", 0, 0, 0, 0, false},

		{"max value with remainder", math.MaxInt8, 3, 42, 1, true},   // 127 / 3 = 42 remainder 1
		{"min value with remainder", math.MinInt8, 3, -42, -2, true}, // -128 / 3 = -42 remainder -2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quot, rem, ok := Quo8(tt.a, tt.b)
			if ok != tt.wantBool {
				t.Errorf("Quo8(%v, %v) status = %v, want %v",
					tt.a, tt.b, ok, tt.wantBool)
			}
			if ok {
				if quot != tt.wantQuot {
					t.Errorf("Quo8(%v, %v) quotient = %v, want %v",
						tt.a, tt.b, quot, tt.wantQuot)
				}
				if rem != tt.wantRem {
					t.Errorf("Quo8(%v, %v) remainder = %v, want %v",
						tt.a, tt.b, rem, tt.wantRem)
				}

				if tt.b != 0 {
					// a = b * quot + rem
					if tt.a != tt.b*tt.wantQuot+tt.wantRem {
						t.Errorf("Quo8(%v, %v) relation check failed: %v != %v * %v + %v",
							tt.a, tt.b, tt.a, tt.b, tt.wantQuot, tt.wantRem)
					}

					// |rem| < |b| (except MinInt8)
					if tt.b != math.MinInt8 && uint8(abs8(t, tt.wantRem)) >= uint8(abs8(t, tt.b)) {
						t.Errorf("Quo8(%v, %v) remainder %v larger than or equal to divisor %v",
							tt.a, tt.b, tt.wantRem, tt.b)
					}

					// rem's sign must be the same as dividend's sign
					if (tt.wantRem < 0) != (tt.a < 0) && tt.wantRem != 0 {
						t.Errorf("Quo8(%v, %v) remainder %v has wrong sign",
							tt.a, tt.b, tt.wantRem)
					}
				}
			}
		})
	}
}

func TestDiv8Panic(t *testing.T) {
	tests := []struct {
		name      string
		a         int8
		b         int8
		wantPanic bool
	}{
		{"normal division", 10, 2, false},
		{"division by zero", 42, 0, true},
		{"min value by -1", math.MinInt8, -1, true},
		{"valid negative division", -128, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("%v: panic = %v, wantPanic = %v",
						tt.name, r != nil, tt.wantPanic)
				}
			}()
			_ = Div8p(tt.a, tt.b)
		})
	}
}

func TestUnsignedAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        uint8
		b        uint8
		want     uint8
		wantBool bool
	}{
		{"zero addition", 0, 0, 0, true},
		{"simple addition", 1, 2, 3, true},
		{"max value", math.MaxUint8, 0, math.MaxUint8, true},
		{"overflow", math.MaxUint8, 1, 0, false},
		{"near overflow", math.MaxUint8 - 1, 1, math.MaxUint8, true},
		{"half values", math.MaxUint8 / 2, math.MaxUint8 / 2, math.MaxUint8 - 1, true},
		{"overflow large", 200, 100, 44, false}, // 300 % 256 = 44
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := UAdd8(tt.a, tt.b)
			if ok != tt.wantBool {
				t.Errorf("UAdd8(%v, %v) overflow check = %v, want %v",
					tt.a, tt.b, ok, tt.wantBool)
			}
			if got != tt.want {
				t.Errorf("UAdd8(%v, %v) = %v, want %v",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestUnsignedSub(t *testing.T) {
	tests := []struct {
		name     string
		a        uint8
		b        uint8
		want     uint8
		wantBool bool
	}{
		{"zero subtraction", 0, 0, 0, true},
		{"simple subtraction", 3, 2, 1, true},
		{"max value", math.MaxUint8, 1, math.MaxUint8 - 1, true},
		{"underflow", 0, 1, 0, false},
		{"near underflow", 1, 1, 0, true},
		{"equal values", 128, 128, 0, true},
		{"max minus max", math.MaxUint8, math.MaxUint8, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := USub8(tt.a, tt.b)
			if ok != tt.wantBool {
				t.Errorf("USub8(%v, %v) overflow check = %v, want %v",
					tt.a, tt.b, ok, tt.wantBool)
			}
			if got != tt.want {
				t.Errorf("USub8(%v, %v) = %v, want %v",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestUnsignedMul(t *testing.T) {
	tests := []struct {
		name     string
		a        uint8
		b        uint8
		want     uint8
		wantBool bool
	}{
		{"zero multiplication", 0, 5, 0, true},
		{"simple multiplication", 2, 3, 6, true},
		{"max value", math.MaxUint8, 1, math.MaxUint8, true},
		{"overflow", math.MaxUint8, 2, 254, false}, // 255 * 2 = 510 -> 254 (mod 256)
		{"near overflow", math.MaxUint8 / 2, 2, math.MaxUint8 - 1, true},
		{"small values 1", 10, 10, 100, true},
		{"small values 2", 15, 16, 240, true},                 // 15 * 16 = 240
		{"overflow large", 16, 16, 0, false},                  // 256 % 256 = 0 (mod 256)
		{"max * max", math.MaxUint8, math.MaxUint8, 1, false}, // 255 * 255 = 65025 -> 1 (mod 256)
		{"half max", 128, 2, 0, false},                        // 128 * 2 = 256 -> 0 (mod 256)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := UMul8(tt.a, tt.b)
			if ok != tt.wantBool {
				t.Errorf("UMul8(%v, %v) overflow check = %v, want %v",
					tt.a, tt.b, ok, tt.wantBool)
			}
			if got != tt.want {
				t.Errorf("UMul8(%v, %v) = %v, want %v",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestUnsignedDiv(t *testing.T) {
	tests := []struct {
		name     string
		a        uint8
		b        uint8
		wantQuot uint8
		wantRem  uint8
		wantBool bool
	}{
		{"simple division", 10, 2, 5, 0, true},
		{"division with remainder", 10, 3, 3, 1, true},
		{"divide by zero", 10, 0, 0, 0, false},
		{"max value", math.MaxUint8, 2, math.MaxUint8 / 2, 1, true},
		{"divide by one", 127, 1, 127, 0, true},
		{"zero dividend", 0, 5, 0, 0, true},
		{"equal values", 128, 128, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quot, rem, ok := UQuotient8(tt.a, tt.b)
			if ok != tt.wantBool {
				t.Errorf("UQuotient8(%v, %v) status = %v, want %v",
					tt.a, tt.b, ok, tt.wantBool)
			}
			if ok {
				if quot != tt.wantQuot {
					t.Errorf("UQuotient8(%v, %v) quotient = %v, want %v",
						tt.a, tt.b, quot, tt.wantQuot)
				}
				if rem != tt.wantRem {
					t.Errorf("UQuotient8(%v, %v) remainder = %v, want %v",
						tt.a, tt.b, rem, tt.wantRem)
				}
			}
		})
	}
}

func TestUnsignedPanic(t *testing.T) {
	tests := []struct {
		name      string
		fn        func()
		wantPanic bool
	}{
		{
			name:      "addition overflow panic",
			fn:        func() { UAdd8p(math.MaxUint8, 1) },
			wantPanic: true,
		},
		{
			name:      "subtraction underflow panic",
			fn:        func() { USub8p(0, 1) },
			wantPanic: true,
		},
		{
			name:      "multiplication overflow panic",
			fn:        func() { UMul8p(math.MaxUint8, 2) },
			wantPanic: true,
		},
		{
			name:      "division by zero panic",
			fn:        func() { UDiv8p(10, 0) },
			wantPanic: true,
		},
		{
			name:      "valid addition no panic",
			fn:        func() { UAdd8p(1, 1) },
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("panic = %v, wantPanic = %v", r != nil, tt.wantPanic)
				}
			}()
			tt.fn()
		})
	}
}

func TestHigherBitOperations(t *testing.T) {
	t.Run("16-bit operations", func(t *testing.T) {
		result, ok := UAdd16(math.MaxUint16-1, 1)
		if !ok || result != math.MaxUint16 {
			t.Errorf("UAdd16(MaxUint16-1, 1) = %v, %v; want %v, true",
				result, ok, math.MaxUint16)
		}
	})

	t.Run("32-bit operations", func(t *testing.T) {
		result, ok := UAdd32(math.MaxUint32-1, 1)
		if !ok || result != math.MaxUint32 {
			t.Errorf("UAdd32(MaxUint32-1, 1) = %v, %v; want %v, true",
				result, ok, math.MaxUint32)
		}
	})

	t.Run("64-bit operations", func(t *testing.T) {
		result, ok := UAdd64(math.MaxUint64-1, 1)
		if !ok || result != math.MaxUint64 {
			t.Errorf("UAdd64(MaxUint64-1, 1) = %v, %v; want %s, true",
				result, ok, "18446744073709551615")
		}
	})
}

func abs8(t *testing.T, x int8) int8 {
	t.Helper()
	if x >= 0 {
		return x
	}
	return -x
}
