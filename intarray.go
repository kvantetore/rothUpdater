package main

import (
	"fmt"
	"strconv"
)

type intarray []int

func (ids *intarray) String() string {
	return fmt.Sprintf("%v", *ids)
}

func (ids *intarray) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*ids = append(*ids, i)
	return nil
}
