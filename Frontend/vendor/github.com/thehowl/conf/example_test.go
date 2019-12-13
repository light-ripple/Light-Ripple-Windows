package conf_test

import (
	"fmt"

	"github.com/thehowl/conf"
)

type myConf struct {
	Name string
	Age  int
}

const myConfString = `Name=Jack
Age=19`

func Example() {
	c := myConf{}
	conf.LoadRaw(&c, []byte(myConfString))
	fmt.Printf("%#v\n", c)
}
