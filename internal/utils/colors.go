package utils

import (
	"fmt"
)

const (
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorBlue   = "\033[34m"
    ColorPurple = "\033[35m"
    ColorCyan   = "\033[36m"
    ColorReset  = "\033[0m"
)

func PrintError(format string, a ...interface{}) {
    fmt.Printf(ColorRed+"ERROR: "+format+ColorReset+"\n", a...)
}

func PrintSuccess(format string, a ...interface{}) {
    fmt.Printf(ColorGreen+"SUCCESS: "+format+ColorReset+"\n", a...)
}

func PrintUserWritten(format string, a ...interface{}) {
    fmt.Printf(ColorBlue+"USER WRITTEN: "+format+ColorReset+"\n", a...)
}

func PrintInfo(format string, a ...interface{}) {
    fmt.Printf(ColorYellow+"INFO: "+format+ColorReset+"\n", a...)
}

func PrintUserRetrieved(format string, a ...interface{}) {
    fmt.Printf(ColorPurple+"USER RETRIEVED: "+format+ColorReset+"\n", a...)
}

func PrintUserProcessed(format string, a ...interface{}) {
    fmt.Printf(ColorCyan+"USER PROCESSED: "+format+ColorReset+"\n", a...)
}