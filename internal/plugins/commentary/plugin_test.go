package commentary

import (
	"github.com/qvistgaard/openrms/internal/plugins/commentary/engines/playht"
	"testing"
	"time"
)

func Test(t *testing.T) {
	c := &Config{}
	c.Plugin.Commentary.Enabled = true
	c.Plugin.Commentary.Engine = "playht"
	c.Plugin.Commentary.PlayHT = &playht.PlayHTConfig{
		Voice:  "Oliver (Advertising)",
		ApiKey: "b6726df13b3648c5868863b5d0ec4d90",
		UserId: "0ZiGPSDn9TN6GsezvpAv9MCDYY93",
		Speed:  1.1,
		Cache:  "cache",
	}

	plugin, err := New(c)

	if err != nil {
		panic(err)
	}

	plugin.Announce("With a sound that could drown out Big Ben, the cars take off, each driver as eager as a kid in a sweet shop.")
	plugin.Announce("Off they zoom, like a group of pensioners on motorized scooters - only faster, and with fewer complaints about their backs.")
	plugin.OptionalAnnouncement("test")
	//	plugin.Announce("Test2")

	time.Sleep(30 * time.Second)
}
