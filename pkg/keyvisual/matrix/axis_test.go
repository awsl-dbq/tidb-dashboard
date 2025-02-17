// Copyright 2021 PingCAP, Inc. Licensed under Apache-2.0.

package matrix

import (
	. "github.com/pingcap/check"
)

var _ = Suite(&testAxisSuite{})

type testAxisSuite struct{}

func (s *testAxisSuite) TestChunkReduce(c *C) {
	testcases := []struct {
		keys      []string
		values    []uint64
		newKeys   []string
		newValues []uint64
	}{
		{
			[]string{"", "a", "b", "c", ""},
			[]uint64{1, 10, 100, 1000},
			[]string{"", "b", ""},
			[]uint64{11, 1100},
		},
		{
			[]string{"", "a", "b", "c", ""},
			[]uint64{1, 10, 100, 1000},
			[]string{"", ""},
			[]uint64{1111},
		},
	}

	for _, testcase := range testcases {
		originChunk := createChunk(testcase.keys, testcase.values)
		reduceChunk := originChunk.Reduce(testcase.newKeys)
		c.Assert(reduceChunk.Values, DeepEquals, testcase.newValues)
	}
}
