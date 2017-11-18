package main

import (
	"strconv"
	"os"
	"fmt"
	"github.com/kvantetore/rothTouchline"
	"flag"
	"text/tabwriter"
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

func listValues(url string) error {
	sensorCount, err := roth.GetSensorCount(url)
	if err != nil {
		return err
	}

	sensors, err := roth.GetSensors(url, sensorCount);
	if err != nil {
		return err
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.AlignRight | tabwriter.Debug)
	fmt.Fprintf(writer, "Id\tSensor\tValve State\tTarget Temp\tCurrent Temp\tMode\tProgram\t\n")
	for _, sensor := range sensors {
		var program string
		switch sensor.Program {
			case roth.ProgramConstant:
				program = "Constant"
			case roth.Program1:
				program = "Program1"
			case roth.Program2:
				program = "Program2"
			case roth.Program3:
				program = "Program3"
		}
		var mode string
		switch sensor.Mode {
			case roth.ModeDay:
				mode = "Day"
			case roth.ModeNight:
				mode = "Night"
			case roth.ModeHoliday:
				mode = "Holiday"
		}
		fmt.Fprintf(writer, "%v\t%v\t%v\t%.2f\t%.2f\t%v\t%v\t\n", sensor.Id, sensor.Name, sensor.GetValveState(), sensor.RoomTemperature, sensor.TargetTemperature, mode, program)
	}
	writer.Flush()
	return nil
}

func setTemperature(url string, sensorIds []int, temperature float32) error {
	for sensorID := range sensorIds {
		err := roth.SetTargetTemperature(url, sensorID, temperature)
		if err != nil {
			return err
		}
	}
	return nil
}

func setMode(url string, sensorIds []int, mode int) error {
	for sensorID := range sensorIds {
		err := roth.SetMode(url, sensorID, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

func setProgram(url string, sensorIds []int, program int) error {
	for sensorID := range sensorIds {
		err := roth.SetProgram(url, sensorID, program)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var sensorIds intarray
	flag.Var(&sensorIds, "sensor", "List of sensor ids")
	action := flag.String("action", "list", "One of list, temp, mode or program")
	value := flag.String("value", "", "Value to use with action (temperature, mode or program)")
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
			listValues(*url)

		case "temp":
			temperature, err := strconv.ParseFloat(*value, 32)
			if err != nil {
				fmt.Printf("Invalid value '%v', expected float\n", value)
			}
			setTemperature(*url, sensorIds, float32(temperature))

		case "mode":
			var mode int
			switch *value {
				case "day": 
					mode = roth.ModeDay
				case "night":
					mode = roth.ModeNight
				case "holiday": 
					mode = roth.ModeHoliday
				default:
					fmt.Printf("Invalid value '%v' should be one of day|night|holiday\n", *value)
					os.Exit(-1)
			}
			setMode(*url, sensorIds, mode)

		case "program":
			var program int
			switch *value {
				case "constant": 
					program = roth.ProgramConstant
				case "program1":
					program = roth.Program1
				case "program2":
					program = roth.Program2
				case "program3":
					program = roth.Program3
				default:
					fmt.Printf("Invalid value '%v' should be one of constant|program1|program2|program3\n", *value)
					os.Exit(-1)
			}
			setProgram(*url, sensorIds, program)
			
		default:
			fmt.Printf("Invalid action %v\n", *action)
			os.Exit(-1)
	}

}