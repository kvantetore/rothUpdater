package main

import (
	"os"
	"fmt"
	"github.com/kvantetore/rothTouchline"
	"flag"
)

const (
	rothManagementURL = "http://ROTH-01A6D5"
)

func makeRange(start int, count int) (ret []int) {
	for i := 0; i < count; i++ {
		ret = append(ret, i + start)
	}
	return ret;
}

func main() {
	var sensorIds intarray
	flag.Var(&sensorIds, "sensor", "List of sensor ids")
	action := flag.String("action", "list", "One of list, temp, mode or program")
	value := flag.Float64("value", -1, "Value to use with action (temperature, mode or program)")
	url := flag.String("url", rothManagementURL, "Base url to Roth Touchline management device")
	help := flag.Bool("help", false, "Prints help")
	flag.Parse()

	if (*help) {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(sensorIds) == 0 {
		sensorCount, err := roth.GetSensorCount(*url)
		sensorIds = makeRange(0, sensorCount)
		if (err != nil) {
			fmt.Printf("Error reading sensor count, %v\n", err)
			os.Exit(-1)
		}
	}

	if len(sensorIds) == 0 {
		fmt.Printf("No sensors\n")
		os.Exit(-1)
	}

	switch *action {
		case "list":
			fmt.Printf("Listing values...\n")
		case "temp":
			fmt.Printf("Setting temperature...%v\n", value)
		case "mode":
			fmt.Printf("Setting mode to %v\n", int(*value))
		default:
			fmt.Printf("Invalid action %v\n", *action)
			os.Exit(-1)
	}

}