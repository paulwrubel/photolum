package cmd

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Hey, it works!\n")
	envVar, isSet := os.LookupEnv("PRINT_THIS_VAR")
	if isSet {
		fmt.Printf("You can find your requested env var [PRINT_THIS_VAR] here: %s\n", envVar)
	} else {
		fmt.Printf("You blew it! Your requested env var [PRINT_THIS_VAR] has not been set!\n")
	}
}
