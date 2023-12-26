package fuel

import (
	"fmt"
	"testing"
)

func Test_average_reportUsage(t *testing.T) {
	t.Run("average calculation", func(t *testing.T) {
		a := average{}
		fmt.Printf("%+v\n", a)
		a = a.reportUsage(5)
		fmt.Printf("%+v\n", a)
		a = a.reportUsage(10)
		fmt.Printf("%+v\n", a)
		a = a.reportUsage(15)
		fmt.Printf("%+v\n", a)
		a = a.reportUsage(0)
		fmt.Printf("%+v\n", a)
		a = a.reportUsage(5)
		fmt.Printf("%+v\n", a)
	})
}
