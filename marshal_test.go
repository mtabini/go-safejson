package safejson

import (
	"testing"
)

type tc struct {
	n string
	t interface{}
	e bool
	r interface{}
}

type testS struct {
	A int `safejson:"a" json:"aa"`
	B int `json:"b"`
	C int `safejson:"test"`
}

type nestS struct {
	A testS `safejson:"a1"`
	B testS
}

type testS2 struct {
	A int `safejson:"a"`
	b int `safejson:"b"`
}

func TestMarshal(t *testing.T) {
	tests := []tc{
		tc{"Nil", nil, false, "null"},
		tc{"Non-struct", 10, true, ""},
		tc{"Struct", testS{1, 2, 3}, false, `{"a":1,"test":3}`},
		tc{"Nested Struct", nestS{testS{1, 2, 3}, testS{1, 2, 3}}, false, `{"a1":{"a":1,"test":3}}`},
		tc{"Slice", []testS{testS{1, 2, 3}}, false, `[{"a":1,"test":3}]`},
		tc{"Unexported Field", testS2{1, 2}, false, `{"a":1}`},
	}

	for _, test := range tests {
		if out, err := Marshal(test.t); err != nil && !test.e {
			t.Errorf("Test `%s` reported error `%s`", test.n, err)
		} else if err == nil && test.e {
			t.Errorf("Test `%s` expected an error, but no error was returned", test.n)
		} else if string(out) != test.r {
			t.Errorf("Test `%s` expected %s, but returned %s", test.n, test.r, string(out))
		}
	}
}

func BenchmarkMarshal(t *testing.B) {
	tests := []tc{
		tc{"Nil", nil, false, "null"},
		tc{"Non-struct", 10, true, ""},
		tc{"Struct", testS{1, 2, 3}, false, `{"a":1,"test":3}`},
		tc{"Nested Struct", nestS{testS{1, 2, 3}, testS{1, 2, 3}}, false, `{"a1":{"a":1,"test":3}}`},
		tc{"Slice", []testS{testS{1, 2, 3}}, false, `[{"a":1,"test":3}]`},
		tc{"Unexported Field", testS2{1, 2}, false, `{"a":1}`},
	}

	for i := 0; i < t.N; i++ {
		for _, test := range tests {
			if out, err := Marshal(test.t); err != nil && !test.e {
				t.Errorf("Test `%s` reported error `%s`", test.n, err)
			} else if err == nil && test.e {
				t.Errorf("Test `%s` expected an error, but no error was returned", test.n)
			} else if string(out) != test.r {
				t.Errorf("Test `%s` expected %s, but returned %s", test.n, test.r, string(out))
			}
		}
	}
}

func ExampleMarshal() {
	s := struct {
		// A field that is marshaled and unmarshaled
		A int `json:"a" safejson:"a"`

		// A field that is unmarshaled using a different name
		B int `json:"b" safejson:"bout"`

		// A field that is ignored when unmarshaling, but marshaled
		C int `json:"-" safejson:"c"`

		// A field that is unmarshaled, but never output during marshaling

		D int `json:"d"`
	}{1, 2, 3, 4}

	Marshal(s) // Returns {"a":1,"bout":2,"c":3}
}
