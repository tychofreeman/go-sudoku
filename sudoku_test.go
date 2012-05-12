package sudoku

import (
    "testing"
)

func GetMissingValue(board []int) int {
    missing := make([]int, len(board))
    missingValue := 0
    for _, v := range board {
        if v > 0 {
            missing[v-1] = v
        }
    }
    for i := range board {
        if missing[i] == 0 {
            missingValue = i+1
        }
    }
    return missingValue
}

func Solution(board []int) [][]int {
    missingValue := GetMissingValue(board)

    result := make([][]int, len(board))
    for i, v := range board {
        if v == 0 {
            v = missingValue
        }
        result[i] = []int {v}
    }
    return result
}

func IsExactly(candidate []int, value int) bool {
    return len(candidate) == 1 && candidate[0] == value
}

func HasAllOf(candidate []int, values []int) bool {
    for _, value := range values {
        found := false
        for _, candidateValue := range candidate {
            if candidateValue == value {
                found = true
            }
        }
        if !found {
            return false
        }
    }
    return true
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

    if !IsExactly(result[0], 1) {
        t.Errorf("Single value in result should be 1, but was %v\n", result[0])
    }
}

func TestFindsMissingNumbersInList(t *testing.T) {
    input := []int {1, 0}
    result := Solution(input)

    if !IsExactly(result[0], 1) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[0])
    }

    if !IsExactly(result[1], 2) {
        t.Errorf("Unknown value was not found. Expected 2, but was %v\n", result[1])
    }
}

func TestFindsMissingNumbersInList2(t *testing.T) {
    input := []int {0, 1}
    result := Solution(input)

    if !IsExactly(result[1], 1) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[1])
    }

    if !IsExactly(result[0], 2) {
        t.Errorf("Unkonwn value was not found. Expected 2, but was %v\n", result[1])
    }
}

func TestFindsMissingNumberInSize3List(t *testing.T) {
    input := []int {1, 2, 0}
    result := Solution(input)
    if !IsExactly(result[0], 1) || !IsExactly(result[1], 2) {
        t.Errorf("Known values should not be changed.\n")
    }

    if !IsExactly(result[2], 3) {
        t.Errorf("Unkonwn value was not found. Expected 3, but was %v\n", result[2])
    }

}

/*
func TestFindsMultipleMissingNumbersInLongerList(t *testing.T) {
    input := []int {1, 0, 0}
    result := Solution(input)

    if !IsExactly(result[0], 1) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[0])
    }

    if !HasAllOf(result[1], []int {2,3}) {
        t.Errorf("Undecidable value should have all possible values. Expected {2,3}, but was %v\n", result[1])
    }
    if !HasAllOf(result[1], []int {2,3}) {
        t.Errorf("Undecidable value should have all possible values. Expected {2,3}, but was %v\n", result[2])
    }
}
*/
