package sudoku

import (
    "testing"
)

func Solution(board []int) []int {
    return board
}

func TestZeroSizedBoard(t *testing.T) {
    emptyBoard := []int {}
    if len(Solution(emptyBoard)) != 0 {
        t.Fail()
    }
}

func TestOneSizedBoard(t *testing.T) {
    oneSized := []int {1}
    result := Solution(oneSized)

    if len(result) != 1 {
        t.Errorf("Length should be 1, but was %v\n", result)
    }

    if result[0] != 1 {
        t.Errorf("Single value in result should be 1, but was %v\n", result[0])
    }
}

func TestFindsMissingNumbersInList(t *testing.T) {
    input := []int {1, 0}
    result := Solution(input)

    if result[0] != 1 {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[0])
    }

    if result[1] != 2 {
/bin/bash: :w: command not found
    }
}
