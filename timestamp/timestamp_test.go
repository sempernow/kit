package timestamp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sempernow/kit/convert"
	"github.com/sempernow/kit/timestamp"
)

func TestTimeUtils(t *testing.T) {
	t.Skip()
	fmt.Println(timestamp.EpochSecToTimeUTC(1594562637))       // 2020-07-12 14:03:57 +0000 UTC
	fmt.Println(timestamp.EpochSecToTimeLocal(1594562637))     // 2020-07-12 10:03:57 -0400 EDT
	fmt.Println(timestamp.EpochMsecToTimeUTC(1594562637123))   // 2020-07-12 14:03:57 +0000 UTC
	fmt.Println(timestamp.EpochMsecToTimeLocal(1594562637123)) // 2020-07-12 10:03:57 -0400 EDT
	fmt.Println(timestamp.EpochSecToMsec(1594562637))          // 1594562637000

	year, month, day := time.Now().Date()
	fmt.Printf("year: %s, month: %s, day: %s\n",
		convert.IntToString(year), month.String(), convert.IntToString(day),
	)
	tt := time.Now()
	yr := tt.Year()  // type int
	mo := tt.Month() // type time.Month
	d := tt.Day()    // type int
	fmt.Printf("year: %s, month: %s, day: %s\n",
		convert.IntToString(yr), mo.String(), convert.IntToString(d),
	)
	x := time.Now()
	l := timestamp.TimeStringLocal(x)
	z := timestamp.TimeStringZulu(x)
	fmt.Println(l)                                                //2020-07-22T10:21:51-04:00
	fmt.Println(z)                                                //2020-07-22T14:21:51Z
	fmt.Println(timestamp.TimeToEpochSec(x), "::", x)             //1595428805
	fmt.Println(timestamp.TimeToEpochSec(x.UTC()), "::", x.UTC()) //1595428805
}
