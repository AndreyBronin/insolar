//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package lazyhead

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func lessTestFn(v1 interface{}, v2 interface{}) bool {
	return v1.(int) < v2.(int)
}

func TestNewHeadedLazySortedList(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, -1)
	require.Equal(t, 0, len(hl.data.data))

	require.Equal(t, 2, cap(hl.data.data))

	hl = NewHeadedLazySortedList(1, lessTestFn, 0)
	require.Equal(t, 2, cap(hl.data.data))

	hl = NewHeadedLazySortedList(3, lessTestFn, 1)
	require.Equal(t, 3, cap(hl.data.data))

	hl = NewHeadedLazySortedList(1, lessTestFn, 3)
	require.Equal(t, 3, cap(hl.data.data))
}

func TestInnerLen(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)

	require.Equal(t, 1, inl.Len())
}

func TestInnerLess(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	inl.Add(2)
	require.Equal(t, lessTestFn(0, 1), inl.Less(0, 1))
}

func TestInnerSwap(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	inl.Add(2)
	inl.Swap(0, 1)
	require.Equal(t, 2, inl.Get(0))

	require.Equal(t, 1, inl.Get(1))
}

func TestInnerAdd(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	require.Equal(t, 1, inl.Get(0))

	require.Equal(t, 1, inl.len)

	inl = innerHeadedLazySortedList{data: make([]interface{}, 2), less: lessTestFn}
	inl.Add(1)
	require.Equal(t, 1, inl.Get(0))

	require.Equal(t, 1, inl.len)

	inl.Add(3)
	require.Equal(t, 1, inl.Get(0))

	require.Equal(t, 3, inl.Get(1))

	require.Equal(t, 2, inl.len)
}

func TestInnerGet(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(2)
	inl.Add(3)
	require.Equal(t, 2, inl.Get(0))

	require.Equal(t, 3, inl.Get(1))

	require.Panics(t, func() { inl.Get(-1) })

	require.Panics(t, func() { inl.Get(2) })
}

func TestLen(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	require.Equal(t, 0, hl.Len())

	hl.Add(2)
	require.Equal(t, 1, hl.Len())

	hl.Add(3)
	require.Equal(t, 2, hl.Len())
}

func TestGet(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(2)
	hl.Add(3)
	require.Equal(t, 2, hl.Get(0))

	require.Equal(t, 3, hl.Get(1))

	require.Panics(t, func() { hl.Get(-1) })

	require.Panics(t, func() { hl.Get(2) })
}

func TestAdd(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(1)
	require.Equal(t, Sorted, hl.sorted)

	hl.Add(2)
	require.Equal(t, UnsortedTail, hl.sorted)

	hl.Add(0)
	require.Equal(t, UnsortedTail, hl.sorted)

	hl = NewHeadedLazySortedList(2, lessTestFn, 2)
	hl.Add(1)
	hl.sorted = UnsortedAll
	hl.Add(2)
	require.Equal(t, Sorted, hl.sorted)

	hl.Add(1)
	require.Equal(t, Sorted, hl.sorted)

	hl = NewHeadedLazySortedList(4, lessTestFn, 3)
	hl.Add(1)
	hl.Add(2)
	require.Equal(t, UnsortedAll, hl.sorted)

	hl.Add(1)
	require.Equal(t, UnsortedAll, hl.sorted)
}

func TestSortAll(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(2)
	hl.SortAll()
	require.Equal(t, Sorted, hl.sorted)

	hl.Add(3)
	require.Equal(t, UnsortedTail, hl.sorted)

	hl.SortAll()
	require.Equal(t, Sorted, hl.sorted)
}

func TestInnerCutOffHeadByLen(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(2)
	inl.Add(3)
	to := inl.cutOffHeadByLen(0, nil)
	require.Equal(t, 0, len(to))

	require.Equal(t, 2, len(inl.data))

	to = inl.cutOffHeadByLen(1, nil)
	require.Equal(t, 1, len(to))

	require.Equal(t, 1, inl.len)

	require.Equal(t, nil, inl.data[1])

	inl.Add(4)
	to2 := make([]interface{}, 1)
	to = inl.cutOffHeadByLen(2, to2)
	require.Equal(t, 3, len(to))

	require.Equal(t, 0, inl.len)

	require.Equal(t, nil, inl.data[1])
}

func TestGetReversedHead(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(2)
	require.Panics(t, func() { hl.GetReversedHead(-1) })

	require.Panics(t, func() { hl.GetReversedHead(2) })

	k := hl.GetReversedHead(1)
	require.Equal(t, 2, k.(int))
}

func TestHasFullHead(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	require.Panics(t, func() { hl.HasFullHead(-1) })

	require.Panics(t, func() { hl.HasFullHead(2) })

	require.True(t, hl.HasFullHead(1))

	require.False(t, hl.HasFullHead(0))
}

func TestCheckHeadLen(t *testing.T) {
	hl := NewHeadedLazySortedList(3, lessTestFn, 1)
	hl.Add(2)
	require.Panics(t, func() { hl.checkAndGetAdjustedHeadLen(-1) })

	require.Panics(t, func() { hl.checkAndGetAdjustedHeadLen(4) })

	hl.Add(3)
	n := hl.checkAndGetAdjustedHeadLen(1)
	require.Equal(t, hl.headLen, n)

	n = hl.checkAndGetAdjustedHeadLen(0)
	require.Equal(t, hl.data.len, n)
}

func TestCutOffHeadInto(t *testing.T) {
	hl := NewHeadedLazySortedList(3, lessTestFn, 1)
	item := 2
	hl.Add(item)
	to := hl.CutOffHeadInto(1, nil)
	require.Equal(t, 1, len(to))

	require.Equal(t, []interface{}{item}, to)
}

func TestCutOffHeadByLenInto(t *testing.T) {
	hl := NewHeadedLazySortedList(3, lessTestFn, 1)
	item1 := 2
	hl.Add(item1)
	to := hl.CutOffHeadByLenInto(0, nil)
	require.Equal(t, []interface{}(nil), to)

	to2 := make([]interface{}, 1)
	to = hl.CutOffHeadByLenInto(0, to2)
	require.Equal(t, len(to2), len(to))

	item2 := 3
	hl.Add(item2)
	hl.sorted = UnsortedAll
	to = hl.CutOffHeadByLenInto(1, to2)
	require.Equal(t, Sorted, hl.sorted)

	require.Equal(t, 1, hl.data.len)

	require.Equal(t, item2, hl.Get(0))

	require.Equal(t, item1, to[1].(int))

	to = hl.CutOffHeadByLenInto(1, to2)
	require.Equal(t, Sorted, hl.sorted)

	require.Equal(t, item2, to[1].(int))

	hl.Add(3)
	hl.Add(2)
	hl.Add(3)
	hl.Add(5)
	hl.Add(7)
	require.Equal(t, UnsortedTail, hl.sorted)

	to = hl.CutOffHeadByLenInto(2, to2)
	require.Equal(t, UnsortedAll, hl.sorted)
}

func TestCutOffHead(t *testing.T) {
	hl := NewHeadedLazySortedList(2, lessTestFn, 1)
	require.Equal(t, []interface{}{}, hl.CutOffHead(0))

	item1 := 4
	hl.Add(item1)
	item2 := 2
	hl.Add(item2)
	require.Equal(t, []interface{}{item2, item1}, hl.CutOffHead(2))

	hl.Add(item1)
	hl.Add(item2)
	item3 := 3
	hl.Add(item3)
	require.Equal(t, []interface{}{item2, item3}, hl.CutOffHead(2))

	require.Panics(t, func() { hl.CutOffHead(3) })
}

func TestCutOffHeadByLen(t *testing.T) {
	hl := NewHeadedLazySortedList(2, lessTestFn, 1)
	require.Equal(t, []interface{}{}, hl.CutOffHeadByLen(0))

	item1 := 4
	hl.Add(item1)
	item2 := 2
	hl.Add(item2)
	require.Equal(t, []interface{}{item2, item1}, hl.CutOffHeadByLen(2))

	hl.Add(item1)
	hl.Add(item2)
	item3 := 3
	hl.Add(item3)
	require.Equal(t, []interface{}{item2, item3}, hl.CutOffHeadByLen(2))

	require.Panics(t, func() { hl.CutOffHeadByLen(3) })
}

func TestFlush(t *testing.T) {
	hl := NewHeadedLazySortedList(2, lessTestFn, 1)
	item1 := 4
	hl.Add(item1)
	item2 := 2
	hl.Add(item2)
	res := hl.Flush()
	require.Equal(t, []interface{}{item2, item1}, res)

	require.Equal(t, Sorted, hl.sorted)
}

func TestGetHeadLen(t *testing.T) {
	headLen := 2
	hl := NewHeadedLazySortedList(headLen, lessTestFn, 1)
	require.Equal(t, headLen, hl.GetHeadLen())
}

func TestGetAvailableHeadLen(t *testing.T) {
	headLen := 2
	hl := NewHeadedLazySortedList(headLen, lessTestFn, 1)
	item1 := 4
	hl.Add(item1)
	item2 := 2
	hl.Add(item2)
	require.Equal(t, headLen, hl.GetAvailableHeadLen(1))

	require.Panics(t, func() { hl.GetAvailableHeadLen(-1) })

	require.Panics(t, func() { hl.GetAvailableHeadLen(3) })

	headLen = 3
	hl = NewHeadedLazySortedList(headLen, lessTestFn, 1)
	hl.Add(item1)
	require.Equal(t, 1, hl.GetAvailableHeadLen(1))
}
