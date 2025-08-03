package output

import (
	"fmt"
)

func PrettyPrint(i interface{}, indent ...string) {
	fmt.Println(Prettify(i, indent...))
}
