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

package stats

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

type StatTable struct {
	columns  []Column
	rows     []*Row
	rowCount int
	summary  []uint32
}

func NewStatTable(maxValue uint8, columns int) StatTable {
	if columns > math.MaxUint16 {
		panic("too many columns")
	}
	r := StatTable{columns: make([]Column, columns), rows: make([]*Row, 0, columns)}
	for i := 0; i < columns; i++ {
		r.columns[i].colIndex = uint16(i)
		r.columns[i].summary = make([]uint16, maxValue+1)
	}
	r.summary = make([]uint32, maxValue+1)
	return r
}

func (t *StatTable) NewRow() *Row {
	nr := NewStatRow(t.MaxValue(), t.ColumnCount())
	return &nr
}

func (t *StatTable) AddRow(row *Row) int {
	row.ensureForTable(t)
	row.rowIndex = len(t.rows)
	t.rows = append(t.rows, row)
	t.rowCount++

	for i, v := range row.values {
		t.columns[i].summary[v]++
		t.summary[v]++
	}
	return row.rowIndex
}

func (t *StatTable) PutRow(rowIndex int, row *Row) {
	row.ensureForTable(t)
	switch {
	case rowIndex == len(t.rows):
		t.rows = append(t.rows, row)
	case rowIndex > len(t.rows):
		t.rows = append(t.rows, make([]*Row, rowIndex-len(t.rows)+1)...)
		t.rows[rowIndex] = row
	case t.rows[rowIndex] != nil:
		panic("row is in use")
	default:
		t.rows[rowIndex] = row
	}
	row.rowIndex = rowIndex
	t.rowCount++

	for i, v := range row.values {
		t.columns[i].summary[v]++
		t.summary[v]++
	}
}

func (t *StatTable) GetRow(rowIndex int) (row *Row, ok bool) {
	if rowIndex >= len(t.rows) {
		return nil, false
	}
	row = t.rows[rowIndex]
	return row, row != nil
}

func (t *StatTable) RemoveRow(rowIndex int) (ok bool) {
	if rowIndex >= len(t.rows) {
		return false
	}
	row := t.rows[rowIndex]
	if row == nil {
		return false
	}
	t.rows[rowIndex] = nil
	for i, v := range row.values {
		t.columns[i].summary[v]--
		t.summary[v]--
	}
	t.rowCount--
	return true
}

func (t *StatTable) RowCount() int {
	return t.rowCount
}

func (t *StatTable) ColumnCount() int {
	return len(t.columns)
}

func (t *StatTable) GetSummaryByValue(value uint8) uint32 {
	return t.summary[value]
}

func (t *StatTable) MaxValue() uint8 {
	return uint8(len(t.summary) - 1)
}

func (t *StatTable) GetColumn(colIndex int) *Column {
	return &t.columns[colIndex]
}

func (t *StatTable) String() string {
	return fmt.Sprintf("stats[v=%d, c=%d, r=%d/%d]", t.MaxValue()+1, t.ColumnCount(), t.RowCount(), len(t.rows))
}

func (t *StatTable) AsText(header string) string {
	return t.TableFmt(header, nil)
}

func (t *StatTable) TableFmt(header string, fmtFn RowValueFormatFunc) string {
	widths := make([]int, t.ColumnCount())
	builder := strings.Builder{}
	if fmtFn != nil {
		builder.WriteString("LEGEND [")
		for i := uint8(0); i <= t.MaxValue(); i++ {
			if i != 0 {
				builder.WriteRune(' ')
			}
			builder.WriteString(fmtFn(i))
		}
		builder.WriteRune(']')
		builder.WriteRune(' ')
	}
	builder.WriteString(header)
	builder.WriteString("\n###")
	for i, c := range t.columns {
		s := fmt.Sprintf("|%03d%+v", c.colIndex, c.summary)
		widths[i] = utf8.RuneCountInString(s)
		builder.WriteString(s)
	}
	builder.WriteString("|∑")
	stringSummary32Fmt(t.summary, &builder, fmtFn)
	builder.WriteByte('\n')
	for i, r := range t.rows {
		if r == nil {
			continue
		}
		builder.WriteString(fmt.Sprintf("%03d", i))
		for j, v := range r.values {
			if fmtFn == nil {
				builder.WriteString(fmt.Sprintf("|%*d", widths[j]-1, v))
			} else {
				builder.WriteString(fmt.Sprintf("|%*s", widths[j]-1, fmtFn(v)))
			}
		}
		builder.WriteString("|∑")

		if fmtFn == nil {
			builder.WriteString(fmt.Sprintf("%+v\n", r.GetSummary()))
		} else {
			stringSummary16Fmt(r.GetSummary(), &builder, fmtFn)
			builder.WriteByte('\n')
		}
	}
	builder.WriteByte('\n')
	return builder.String()
}
