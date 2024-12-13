package utils

import (
	"fmt"
	"time"
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
    fmt.Printf(ColorRed+"[%s] ERROR: "+format+ColorReset+"\n", append([]interface{}{time.Now().Format("15:04:05")}, a...)...)
}

func PrintSuccess(format string, a ...interface{}) {
    fmt.Printf(ColorGreen+"[%s] SUCCESS: "+format+ColorReset+"\n", append([]interface{}{time.Now().Format("15:04:05")}, a...)...)
}

func PrintUserWritten(format string, a ...interface{}) {
    fmt.Printf(ColorBlue+"[%s] USER WRITTEN: "+format+ColorReset+"\n", append([]interface{}{time.Now().Format("15:04:05")}, a...)...)
}

func PrintInfo(format string, a ...interface{}) {
    fmt.Printf(ColorYellow+"[%s] INFO: "+format+ColorReset+"\n", append([]interface{}{time.Now().Format("15:04:05")}, a...)...)
}

func PrintUserRetrieved(format string, a ...interface{}) {
    fmt.Printf(ColorPurple+"[%s] USER RETRIEVED: "+format+ColorReset+"\n", append([]interface{}{time.Now().Format("15:04:05")}, a...)...)
}

func PrintUserProcessed(format string, a ...interface{}) {
    fmt.Printf(ColorCyan+"[%s] USER PROCESSED: "+format+ColorReset+"\n", append([]interface{}{time.Now().Format("15:04:05")}, a...)...)
}