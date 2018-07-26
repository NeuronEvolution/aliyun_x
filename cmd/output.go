package main

import (
	"bytes"
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"io/ioutil"
	"os"
	"time"
)

func output(result *cloud.ResourceManagement, duration time.Duration) (err error) {
	fmt.Printf("time=%f\n", duration.Seconds())

	a5506, err := ioutil.ReadFile("./_output/submit_20180710_031628.csv")
	if err != nil {
		return err
	}

	outputFile := fmt.Sprintf("_output/submit_%s", time.Now().Format("20060102_150405"))

	buf := bytes.NewBuffer(a5506)
	buf.WriteString("#\n")
	buf.Write(result.DeployCommandHistory.OutputCSV())

	err = ioutil.WriteFile(outputFile+".csv", buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	summaryBuf := bytes.NewBufferString("")
	result.DebugStatus(summaryBuf)
	summaryBuf.WriteString(fmt.Sprintf("time=%f秒\n", duration.Seconds()))
	err = ioutil.WriteFile(fmt.Sprintf(outputFile+"_summary.csv"),
		summaryBuf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
