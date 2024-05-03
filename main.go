package main

import (
	"github.com/aauren/evermarkable/cmd"
	"k8s.io/klog/v2"
)

func main() {
	if err := cmd.Execute(); err != nil {
		klog.Fatalf("error encountered during run: %v", err)
	}
}
