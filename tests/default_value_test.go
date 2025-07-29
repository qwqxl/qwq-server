package qwqtest

import (
	"fmt"
	"qwqserver/pkg/defaultvalue"
	"testing"
	"time"
)

type Example struct {
	IDs       [][]int         `qwq-default:"[[3],[2]]"`
	Names     []string        `qwq-default:"2,xiaolin,xiaomi"`
	Options   []*Option       `qwq-default:"2"`
	Durations []time.Duration `qwq-default:"2"`
}

type Option struct {
	Value int `qwq-default:"5"`
}

func TestDefaultValue(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		ex := Example{}
		if err := defaultvalue.SetDefaults(&ex); err != nil {
			panic(err)
		}

		for i := range ex.IDs {
			fmt.Println("ids: ", ex.IDs[i]) // [0 0]
		}
		fmt.Println("names: ", ex.Names) // ["", ""]
		for _, opt := range ex.Options {
			fmt.Println("opt: ", opt.Value) // 5
		}
		fmt.Println("d: ", ex.Durations) // [0s, 0s]
	})
}
