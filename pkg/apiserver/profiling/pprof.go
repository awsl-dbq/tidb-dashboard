// Copyright 2021 PingCAP, Inc. Licensed under Apache-2.0.

package profiling

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/pingcap/tidb-dashboard/pkg/apiserver/model"
)

type pprofOptions struct {
	duration           uint
	fileNameWithoutExt string

	target        *model.RequestTargetNode
	fetcher       *profileFetcher
	profilingType TaskProfilingType
}

func fetchPprof(op *pprofOptions) (string, TaskRawDataType, error) {
	fetcher := &fetcher{profileFetcher: op.fetcher, target: op.target}
	tmpPath, rawDataType, err := fetcher.FetchAndWriteToFile(op.duration, op.fileNameWithoutExt, op.profilingType)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch and write to temp file: %v", err)
	}

	return tmpPath, rawDataType, nil
}

type fetcher struct {
	target         *model.RequestTargetNode
	profileFetcher *profileFetcher
}

func (f *fetcher) FetchAndWriteToFile(duration uint, fileNameWithoutExt string, profilingType TaskProfilingType) (string, TaskRawDataType, error) {
	var profilingRawDataType TaskRawDataType
	var fileExtenstion string
	secs := strconv.Itoa(int(duration))
	var url string
	switch profilingType {
	case ProfilingTypeCPU:
		url = "/debug/pprof/profile?seconds=" + secs
		profilingRawDataType = RawDataTypeProtobuf
		fileExtenstion = "*.proto"
	case ProfilingTypeHeap:
		url = "/debug/pprof/heap"
		profilingRawDataType = RawDataTypeProtobuf
		fileExtenstion = "*.proto"
	case ProfilingTypeGoroutine:
		url = "/debug/pprof/goroutine?debug=1"
		profilingRawDataType = RawDataTypeText
		fileExtenstion = "*.txt"
	case ProfilingTypeMutex:
		url = "/debug/pprof/mutex?debug=1"
		profilingRawDataType = RawDataTypeText
		fileExtenstion = "*.txt"
	}

	tmpfile, err := ioutil.TempFile("", fileNameWithoutExt+fileExtenstion)
	if err != nil {
		return "", "", fmt.Errorf("failed to create tmpfile to write profile: %v", err)
	}

	defer func() {
		_ = tmpfile.Close()
	}()

	resp, err := (*f.profileFetcher).fetch(&fetchOptions{ip: f.target.IP, port: f.target.Port, path: url})
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch profile with %v format: %v", fileExtenstion, err)
	}

	_, err = tmpfile.Write(resp)
	if err != nil {
		return "", "", fmt.Errorf("failed to write profile: %v", err)
	}

	return tmpfile.Name(), profilingRawDataType, nil
}
