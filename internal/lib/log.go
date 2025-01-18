package lib

import "fmt"

func LogMessage(errScope, errTemplate string, payload ...interface{}) {
	fmt.Printf("[%s]: %s\n", errScope, fmt.Sprintf(errTemplate, payload...))
}

func LogError(errScope, errTemplate string, err error, payload ...interface{}) {
	fmt.Println(
		fmt.Errorf("[%s]: %s; see: %w", errScope, fmt.Sprintf(errTemplate, payload...), err),
	)
}
