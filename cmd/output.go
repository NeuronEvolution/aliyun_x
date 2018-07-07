package main

import (
	"bytes"
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"io/ioutil"
	"os"
	"time"
)

func output(result *cloud.ResourceManagement, duration time.Duration) {
	fmt.Printf("time=%f\n", duration.Seconds())

	outputFile := fmt.Sprintf("_output/submit_%s", time.Now().Format("20060102_150405"))
	err := ioutil.WriteFile(outputFile+".csv", result.DeployCommandHistory.OutputCSV(), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	summaryBuf := bytes.NewBufferString("")
	result.DebugStatus(summaryBuf)
	summaryBuf.WriteString(fmt.Sprintf("time=%f\n", duration.Seconds()))
	err = ioutil.WriteFile(fmt.Sprintf(outputFile+"_summary.csv"),
		summaryBuf.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
}
