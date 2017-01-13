package flatjson

import (
	"reflect"
	"testing"
)

func TestScanObjects(t *testing.T) {
	tests := []struct {
		Name string

		Data string

		WantStart int
		WantEnd   int

		WantNumber []tnumber
		WantString []tstring
		WantBool   []tbool
		WantNull   []tnull

		WantErrError  string
		WantErrOffset int
	}{

		// happy path
		{
			Name:      "empty object",
			Data:      `{}`,
			WantStart: 0,
			WantEnd:   2,
		},
		{
			Name:      "simple string object",
			Data:      `{"hello":"world"}`,
			WantStart: 0,
			WantEnd:   17,
			WantString: []tstring{
				{name: `"hello"`, value: `"world"`},
			},
		},
		{
			Name:      "simple number object",
			Data:      `{"hello":-49.14159e-2}`,
			WantStart: 0,
			WantEnd:   22,
			WantNumber: []tnumber{
				{name: `"hello"`, value: -49.14159e-2},
			},
		},
		{
			Name:      "simple true bool object",
			Data:      `{"hello":true}`,
			WantStart: 0,
			WantEnd:   14,
			WantBool: []tbool{
				{name: `"hello"`, value: true},
			},
		},
		{
			Name:      "simple false bool object",
			Data:      `{"hello":false}`,
			WantStart: 0,
			WantEnd:   15,
			WantBool: []tbool{
				{name: `"hello"`, value: false},
			},
		},
		{
			Name:      "simple null object",
			Data:      `{"hello":null}`,
			WantStart: 0,
			WantEnd:   14,
			WantNull: []tnull{
				{name: `"hello"`},
			},
		},

		{
			Name:      "simple composite object",
			Data:      `{"a":"1","b":2,"c":true,"d":false,"e":null}`,
			WantStart: 0,
			WantEnd:   43,
			WantString: []tstring{
				{name: `"a"`, value: `"1"`},
			},
			WantNumber: []tnumber{
				{name: `"b"`, value: 2},
			},
			WantBool: []tbool{
				{name: `"c"`, value: true},
				{name: `"d"`, value: false},
			},
			WantNull: []tnull{
				{name: `"e"`},
			},
		},

		{
			Name: "composite object with whitespace",
			Data: `
            {
                "a" :   "1",
                "b" :   2,
                "c" :   true,
                "d" :   false,
                "e":    null
}`,
			WantStart: 13,
			WantEnd:   162,
			WantString: []tstring{
				{name: `"a"`, value: `"1"`},
			},
			WantNumber: []tnumber{
				{name: `"b"`, value: 2},
			},
			WantBool: []tbool{
				{name: `"c"`, value: true},
				{name: `"d"`, value: false},
			},
			WantNull: []tnull{
				{name: `"e"`},
			},
		},

		{
			Name: "composite object with weird whitespace",
			Data: `
            {
                "a" :   "1"
                ,"b" :   2,
                "c" :true ,
                "d" :   false
                ,
                "e":    null
}`,
			WantStart: 13,
			WantEnd:   177,
			WantString: []tstring{
				{name: `"a"`, value: `"1"`},
			},
			WantNumber: []tnumber{
				{name: `"b"`, value: 2},
			},
			WantBool: []tbool{
				{name: `"c"`, value: true},
				{name: `"d"`, value: false},
			},
			WantNull: []tnull{
				{name: `"e"`},
			},
		},

		// special cases
		{
			Name: "empty object with whitespace",
			Data: `
    {

    }

`,
			WantStart: 5,
			WantEnd:   13,
		},

		// errors
		{
			Name:         "empty string",
			Data:         ``,
			WantErrError: noOpeningBracketFound,
		},
		{
			Name:          "only opening brakcet",
			Data:          `{`,
			WantErrError:  endOfDataNoClosingBracket,
			WantErrOffset: 1,
		},
		{
			Name:          "single pair with no closing bracket (number)",
			Data:          `{"hello":0`,
			WantErrError:  endOfDataNoClosingBracket,
			WantErrOffset: 11,
		},
		{
			Name:          "single pair with no closing bracket (bool)",
			Data:          `{"hello":true`,
			WantErrError:  endOfDataNoClosingBracket,
			WantErrOffset: 14,
		},
		{
			Name:          "single pair with no closing bracket (bool and comma)",
			Data:          `{"hello":true,`,
			WantErrError:  endOfDataNoClosingBracket,
			WantErrOffset: 14,
		},
		{
			Name:          "single pair with no closing bracket (number) and space",
			Data:          `{"hello":0 `,
			WantErrError:  endOfDataNoClosingBracket,
			WantErrOffset: 12,
		},
		{
			Name:          "single pair with no closing bracket (bool) and space",
			Data:          `{"hello":true `,
			WantErrError:  endOfDataNoClosingBracket,
			WantErrOffset: 15,
		},
		{
			Name:          "single pair with no closing bracket (bool and comma)",
			Data:          `{"hello":true, `,
			WantErrError:  endOfDataNoNamePair,
			WantErrOffset: 15,
		},

		{
			Name:          "missing name in name/value pair",
			Data:          `{:true, `,
			WantErrError:  expectingNameBeforeValue + ", " + reachedEndScanningCharacters,
			WantErrOffset: 1,
		},
		{
			Name:          "missing semicolon in name/value pair",
			Data:          `{"hello" true, `,
			WantErrError:  noSemicolonFound,
			WantErrOffset: 9,
		},
		{
			Name:          "nothing follows the name",
			Data:          `{"hello" `,
			WantErrError:  endOfDataNoSemicolon,
			WantErrOffset: 9,
		},
		{
			Name:          "nothing follows the semicolon",
			Data:          `{"hello": `,
			WantErrError:  endOfDataNoValue,
			WantErrOffset: 10,
		},
		{
			Name:          "malformed number value (garbage)",
			Data:          `{"hello": 7162hhhh}`,
			WantErrError:  malformedNumber,
			WantErrOffset: 10,
		},
		{
			Name:          "malformed number value (incomplete)",
			Data:          `{"hello": 7162.}`,
			WantErrError:  beginNumberValueButError + ", " + scanningForFraction + ", " + needAtLeastOneDigit,
			WantErrOffset: 10,
		},
		{
			Name:          "malformed string value (incomplete)",
			Data:          `{"hello": "world}`,
			WantErrError:  beginStringValueButError + ", " + reachedEndScanningCharacters,
			WantErrOffset: 10,
		},
		{
			Name:          "random crap for value",
			Data:          `{"hello": lololool}`,
			WantErrError:  expectValueButNoKnownType,
			WantErrOffset: 10,
		},

		{
			Name:          "no closing bracket at end of object",
			Data:          `{"hello": "hello"   `,
			WantErrError:  endOfDataNoClosingBracket,
			WantErrOffset: 21,
		},
	}

	for _, tt := range tests {
		t.Logf("====> %s", tt.Name)

		data := []byte(tt.Data)
		var gotNumber []tnumber
		onNumber := func(v Number) {
			gotNumber = append(gotNumber, tnumber{
				name:  v.Name.String(data),
				value: v.Value,
			})
		}
		var gotString []tstring
		onString := func(v String) {
			gotString = append(gotString, tstring{
				name:  v.Name.String(data),
				value: v.Value.String(data),
			})
		}
		var gotBool []tbool
		onBool := func(v Bool) {
			gotBool = append(gotBool, tbool{
				name:  v.Name.String(data),
				value: v.Value,
			})
		}
		var gotNull []tnull
		onNull := func(v Null) {
			gotNull = append(gotNull, tnull{
				name: v.Name.String(data),
			})
		}

		start, end, err := ScanObject([]byte(data), onNumber, onString, onBool, onNull)

		gotErr, _ := err.(*SyntaxError)

		// if we expect errors
		if tt.WantErrError != "" && gotErr == nil {
			t.Errorf("want an error, got none")
		} else if tt.WantErrError != "" && gotErr != nil {
			wantOffset := tt.WantErrOffset
			if wantOffset != gotErr.Offset {
				t.Errorf("want err offset %d, was %d", wantOffset, gotErr.Offset)
			}
			if want, got := tt.WantErrError, gotErr.Error(); want != got {
				t.Errorf("want error: %q", want)
				t.Errorf(" got error: %q", got)
			}
		} else if gotErr != nil {
			t.Errorf("offset=%d", gotErr.Offset)
			t.Error(gotErr)
		} else {

			if want, got := tt.WantStart, start; want != got {
				t.Errorf("want start %+v", want)
				t.Errorf(" got start %+v", got)
			}

			if want, got := tt.WantEnd, end; want != got {
				t.Errorf("want end %+v", want)
				t.Errorf(" got end %+v", got)
			}

			if want, got := tt.WantNumber, gotNumber; !reflect.DeepEqual(want, got) {
				t.Errorf("want number %+v", want)
				t.Errorf(" got number %+v", got)
			}

			if want, got := tt.WantString, gotString; !reflect.DeepEqual(want, got) {
				t.Errorf("want string %+v", want)
				t.Errorf(" got string %+v", got)
			}

			if want, got := tt.WantBool, gotBool; !reflect.DeepEqual(want, got) {
				t.Errorf("want bool %+v", want)
				t.Errorf(" got bool %+v", got)
			}

			if want, got := tt.WantNull, gotNull; !reflect.DeepEqual(want, got) {
				t.Errorf("want null %+v", want)
				t.Errorf(" got null %+v", got)
			}
		}
	}
}

type tnumber struct {
	name  string
	value float64
}

type tstring struct {
	name  string
	value string
}

type tbool struct {
	name  string
	value bool
}

type tnull struct {
	name string
}
