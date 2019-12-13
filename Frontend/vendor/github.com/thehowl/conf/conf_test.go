package conf

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const sampleForParsing = `Key1=Value
Test=Xd
What=1294; This be a comment.
; This be another comment.


Multi=Line\
                String!
Key\=With\=Equal=Indeed!
ShouldNotOutput=
=ShouldNotOutput
ShouldNotOutput

Can=Contain\;Semicolons!
EqualSign=Can also = not = be = escaped!

ThisShouldComeUp=Yup`

func TestParse(t *testing.T) {
	kvs := Parse([]byte(sampleForParsing))
	t.Log("Result:")
	for _, v := range kvs {
		t.Logf("  %s - %s\n", v.Field, v.Value)
	}
}

func BenchmarkParseWithSample(b *testing.B) {
	sp := []byte(sampleForParsing)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(sp)
	}
}

const (
	escapeOriginal = `;aqsdwe103====aa\a\\aa!a;a

xd|`
	escapeExpected = `\;aqsdwe103\=\=\=\=aa\\a\\\\aa!a\;a\
\
xd|`
)

func TestEscape(t *testing.T) {
	escaped := Escape(escapeOriginal)
	if escaped != escapeExpected {
		t.Fatalf("Expected Escape to return '%s', got '%s' instead.", escapeExpected, escaped)
	}
}

func BenchmarkEscape(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Escape(escapeOriginal)
	}
}

const sampleForUnmarshaling = `
; Welcome to hell.
; As we have to try every single possible value a struct can have, all of this is required.
; Fuck.
TestInt=1510591
TestInt8=-100
TestInt16=-12844
TestInt32=-3211
TestInt64=410491

TestUint=284129419
TestUint8=255
TestUint16=65535
TestUint32=151529;Can't remember the max value. Won't bother googling.
TestUint64=24148999

TestFloat32=24.1123
TestFloat64=110.134141223

TestBool=1

TestString=How much wood would a woodchuck chuck if a woodchuck could chuck wood?
TestEmpty=
`

type sampleTest struct {
	TestInt   int
	TestInt8  int8
	TestInt16 int16
	TestInt32 int32
	TestInt64 int64

	TestUint   uint
	TestUint8  uint8
	TestUint16 uint16
	TestUint32 uint32
	TestUint64 uint64

	TestFloat32 float32
	TestFloat64 float64

	TestBool bool

	TestString string
	TestEmpty  string
}

func (si sampleTest) check(t *testing.T) {
	if true &&
		si.TestInt == 1510591 &&
		si.TestInt8 == -100 &&
		si.TestInt16 == -12844 &&
		si.TestInt32 == -3211 &&
		si.TestInt64 == 410491 &&

		si.TestUint == 284129419 &&
		si.TestUint8 == 255 &&
		si.TestUint16 == 65535 &&
		si.TestUint32 == 151529 &&
		si.TestUint64 == 24148999 &&

		si.TestFloat32 == 24.1123 &&
		si.TestFloat64 == 110.134141223 &&

		si.TestBool == true &&

		si.TestString == "How much wood would a woodchuck chuck if a woodchuck could chuck wood?" &&
		si.TestEmpty == "" {
		t.Log("It surprisingly worked.")
	} else {
		t.Log("Nope.")
		t.Fatalf("%#v", si)
	}
}

func TestLoad(t *testing.T) {
	// Prepare file
	err := ioutil.WriteFile("test.conf", []byte(sampleForUnmarshaling), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.conf")

	// Prepare struct
	si := sampleTest{}

	// Load
	err = Load(&si, "test.conf")
	if err != nil {
		t.Fatal(err)
	}

	// Check for all values to be valid.
	si.check(t)
}

func TestLoadShouldErrNotAStruct(t *testing.T) {
	if err := Load(2, "test.conf"); err != ErrNotAStruct {
		t.Fatalf("Should have panicked with ErrNotAStruct, got error '%s' instead.", err)
	}
}

func TestLoadShouldErrNoFile(t *testing.T) {
	si := sampleTest{}
	if err := Load(&si, "test.conf"); err != ErrNoFile {
		t.Fatalf("Should have panicked with ErrNoFile, got error '%s' instead.", err)
	}
}

func TestMustLoadShouldPanic(t *testing.T) {
	defer func() {
		c := recover()
		if c == nil {
			t.Fatal("MustLoad with wrong values didn't panic!")
		}
	}()
	MustLoad(2, "test.conf")
}

func TestFailBool(t *testing.T) {
	// Prepare file
	err := ioutil.WriteFile("test.conf", []byte(`TestBool=Lolno.`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.conf")

	// Prepare struct
	si := sampleTest{}

	// Load
	err = Load(&si, "test.conf")
	if err == nil {
		t.Fatal("Should have returned an error, didn't")
	}
}

func TestFailInt(t *testing.T) {
	// Prepare file
	err := ioutil.WriteFile("test.conf", []byte(`TestInt=10934104912049120491204912031301293102`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.conf")

	// Prepare struct
	si := sampleTest{}

	// Load
	err = Load(&si, "test.conf")
	if err == nil {
		t.Fatal("Should have returned an error, didn't")
	}
}

func TestFailFloat(t *testing.T) {
	// Prepare file
	err := ioutil.WriteFile("test.conf", []byte(`TestFloat32=2931.23111aaddfffeep`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.conf")

	// Prepare struct
	si := sampleTest{}

	// Load
	err = Load(&si, "test.conf")
	if err == nil {
		t.Fatal("Should have returned an error, didn't")
	}
}

func TestFailUint(t *testing.T) {
	// Prepare file
	err := ioutil.WriteFile("test.conf", []byte(`TestUint8=asd`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.conf")

	// Prepare struct
	si := sampleTest{}

	// Load
	err = Load(&si, "test.conf")
	if err == nil {
		t.Fatal("Should have returned an error, didn't")
	}
}

func TestLoadRaw(t *testing.T) {
	si := sampleTest{}

	LoadRaw(&si, []byte(sampleForUnmarshaling))
	si.check(t)
}

func TestMustLoadRawShouldPanic(t *testing.T) {
	defer func() {
		c := recover()
		if c == nil {
			t.Fatal("MustLoad with wrong values didn't panic!")
		}
	}()
	MustLoadRaw(2, []byte("a"))
}

func BenchmarkLoadRaw(b *testing.B) {
	si := &sampleTest{}
	data := []byte(sampleForUnmarshaling)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadRaw(si, data)
	}
}

type exportTypeTest struct {
	Name string
	Age  int `description:"The age of the patient.\nDescriptions can be multi-line."`
}

const expectedExportType = `Name=Jack
; The age of the patient.
; Descriptions can be multi-line.
Age=19
`

func TestMustExportRaw(t *testing.T) {
	defer func() {
		c := recover()
		if c != nil {
			t.Fatal(c)
		}
	}()
	e := &exportTypeTest{
		Name: "Jack",
		Age:  19,
	}
	data := MustExportRaw(e)
	if string(data) != expectedExportType {
		t.Fatalf("Expected '%s', got '%s'", expectedExportType, string(data))
	}
}

func TestMustExportShouldPanic(t *testing.T) {
	defer func() {
		c := recover()
		if c == nil {
			t.Fatal(c)
		}
	}()
	MustExport(2, "a.conf")
}

func TestMustExportRawShouldPanic(t *testing.T) {
	defer func() {
		c := recover()
		if c == nil {
			t.Fatal(c)
		}
	}()
	MustExportRaw(2)
}

type exportWithInvalidType struct {
	InvalidType []byte
	Name        string
}

func TestExport(t *testing.T) {
	err := Export(exportWithInvalidType{
		InvalidType: []byte("well"),
		Name:        "xd",
	}, "test.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("test.conf")
}

func TestFullChain(t *testing.T) {
	initial := exportTypeTest{
		Name: "Pisellone",
		Age:  133337,
	}
	c := MustExportRaw(initial)
	secondary := exportTypeTest{}
	MustLoadRaw(&secondary, c)
	if !reflect.DeepEqual(initial, secondary) {
		t.Fatalf("Initial struct %#v is not the same as the derivate %#v.", initial, secondary)
	}
}

func TestCRLF(t *testing.T) {
	const w = "Key1=Nice\r\nKey2=Meme"
	vals := Parse([]byte(w))
	for _, i := range vals {
		switch i.Field {
		case "Key1":
			if i.Value != "Nice" {
				t.Fatalf("Expected '%s', got '%s'", "Nice", i.Value)
			}
		case "Key2":
			if i.Value != "Meme" {
				t.Fatalf("Expected '%s', got '%s'", "Meme", i.Value)
			}
		default:
			t.Fatalf("Unexpected key '%s'", i.Field)
		}
	}
}
