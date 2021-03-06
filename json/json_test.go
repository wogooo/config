package json

import (
	"fmt"
	"github.com/gookit/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Example() {
	config.WithOptions(config.ParseEnv)

	// add Decoder and Encoder
	config.AddDriver(Driver)

	err := config.LoadFiles("testdata/json_base.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())

	err = config.LoadFiles("testdata/json_other.json")
	// config.LoadFiles("testdata/json_base.json", "testdata/json_other.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())
	fmt.Print("get config example:\n")

	name, ok := config.String("name")
	fmt.Printf("get string\n - ok: %v, val: %v\n", ok, name)

	arr1, ok := config.Strings("arr1")
	fmt.Printf("get array\n - ok: %v, val: %#v\n", ok, arr1)

	val0, ok := config.String("arr1.0")
	fmt.Printf("get sub-value by path 'arr.index'\n - ok: %v, val: %#v\n", ok, val0)

	map1, ok := config.StringMap("map1")
	fmt.Printf("get map\n - ok: %v, val: %#v\n", ok, map1)

	val0, ok = config.String("map1.key")
	fmt.Printf("get sub-value by path 'map.key'\n - ok: %v, val: %#v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.DefString("envKey1", ""))

	// set value
	config.Set("name", "new name")
	name, ok = config.String("name")
	fmt.Printf("set string\n - ok: %v, val: %v\n", ok, name)

	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, config.JSON)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("json", Driver.Name())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.JSON))
	c.AddDriver(Driver)

	st.True(c.HasDecoder(config.JSON))
	st.True(c.HasEncoder(config.JSON))

	m := struct {
		N string
	}{}
	err := Decoder([]byte(`{
// comments
"n":"v"}
`), &m)
	st.Nil(err)
	st.Equal("v", m.N)
}
