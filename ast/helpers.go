package ast

import "fmt"

func stringResult(result interface{}, err error) (string, error) {
	return result.(string), err
}

func noErrStringResult(result interface{}, err error) string {
	if err != nil {
		err = fmt.Errorf("noErrStringResult: %w", err)
		panic(err)
	}
	return result.(string)
}
