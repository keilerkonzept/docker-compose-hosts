package main

import (
	"fmt"
	"os"
)

func expandEnv(s string) (string, error) {
	var unset []string
	mapping := func(k string) string {
		v, ok := os.LookupEnv(k)
		if !ok {
			unset = append(unset, k)
			return "$" + k
		}
		return v
	}
	s = expand(s, '$', mapping)
	var err error
	if len(unset) > 0 {
		err = fmt.Errorf("unset environment variables used: %q", unset)
	}
	return s, err
}
