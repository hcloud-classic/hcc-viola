package main

import (
	violaInit "hcc/viola/init"
)

func init() {
	err := violaInit.MainInit()
	if err != nil {
		panic(err)
	}
}

func main() {

	defer func() {
		violinEnd.MainEnd()
	}()

}
