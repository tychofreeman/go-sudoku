package sudoku

import (
    "testing"
)

func Solution(board []int) []int {
    return nil
}

func TestZeroSizedBoard(t *testing.T) {
    emptyBoard := []int {}
    if len(Solution(emptyBoard)) != 0 {
        t.Fail()
    }
}
