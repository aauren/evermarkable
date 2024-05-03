package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"k8s.io/klog/v2"
)

func ReadCode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter one-time code (go to https://my.remarkable.com/device/desktop/connect): ")
	code, _ := reader.ReadString('\n')

	code = strings.TrimSuffix(code, "\n")
	code = strings.TrimSuffix(code, "\r")

	if len(code) != 8 {
		klog.Error("Code has the wrong length, it should be 8")
		return ReadCode()
	}

	return code
}
