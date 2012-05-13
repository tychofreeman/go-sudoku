package sudoku

import (
    "testing"
)

func MapMissingValues(board [][]int) []int {
    found := make([]int, len(board))
    for _, v := range board {
        if len(v) == 1 {
            found[v[0]-1] = v[0]
        }
    }

    notFound := make([]int, 0)
    for i, v := range found {
        if v == 0 {
            notFound = append(notFound, i + 1)
        }
    }
    return notFound
}

func union(prevCalcd, remaining []int) []int {
    if len(prevCalcd) > 0 {
        rtn := make([]int, 0)
        for _, pv := range prevCalcd {
            for _, rv := range remaining {
                if pv == rv {
                    rtn = append(rtn, pv)
                }
            }
        }
        return rtn
    }
    return remaining
}

func Solution(board [][]int) [][]int {
    notFound := MapMissingValues(board)

    missingValue := make([][]int, len(board))
    for i := range board {
        missingValue[i] = make([]int, 0)
        if len(board[i]) != 1 {
            missingValue[i] = union(board[i], notFound)
        } else {
            missingValue[i] = board[i]
        }
    }
    return missingValue
}

func IsExactly(candidate []int, value []int) bool {
    if len(candidate) == len(value) {
        for i := range candidate {
            if candidate[i] != value[i] {
                return false
            }
        }
    } else {
        return false
    }
    return true
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
    emptyBoard := [][]int {}
    if len(Solution(emptyBoard)) != 0 {
        t.Fail()
    }
}

func TestOneSizedBoard(t *testing.T) {
    oneSized := [][]int {[]int{1}}
    result := Solution(oneSized)

    if len(result) != 1 {
        t.Errorf("Length should be 1, but was %v\n", result)
    }

    if !IsExactly(result[0], []int{1}) {
        t.Errorf("Single value in result should be 1, but was %v\n", result[0])
    }
}

func TestFindsMissingNumbersInList(t *testing.T) {
    input := [][]int {[]int{1}, []int{}}
    result := Solution(input)

    if !IsExactly(result[0], []int{1}) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[0])
    }

    if !IsExactly(result[1], []int{2}) {
        t.Errorf("Unknown value was not found. Expected 2, but was %v\n", result[1])
    }
}

func TestFindsMissingNumbersInList2(t *testing.T) {
    input := [][]int {[]int{}, []int{1}}
    result := Solution(input)

    if !IsExactly(result[1], []int{1}) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[1])
    }

    if !IsExactly(result[0], []int{2}) {
        t.Errorf("Unknown value was not found. Expected 2, but was %v\n", result[0])
    }
}

func TestFindsMissingNumberInSize3List(t *testing.T) {
    input := [][]int {[]int{1}, []int{2}, []int{}}
    result := Solution(input)
    if !IsExactly(result[0], []int{1}) || !IsExactly(result[1], []int{2}) {
        t.Errorf("Known values should not be changed.\n")
    }

    if !IsExactly(result[2], []int{3}) {
        t.Errorf("Unkonwn value was not found. Expected 3, but was %v\n", result[2])
    }

}

func TestFindsMultipleMissingNumbersInLongerList(t *testing.T) {
    input := [][]int {[]int{1}, []int{}, []int{}}
    result := Solution(input)

    if !IsExactly(result[0], []int{1}) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[0])
    }

    if !HasAllOf(result[1], []int {2,3}) {
        t.Errorf("Undecidable value should have all possible values. Expected {2,3}, but was %v\n", result[1])
    }
    if !HasAllOf(result[1], []int {2,3}) {
        t.Errorf("Undecidable value should have all possible values. Expected {2,3}, but was %v\n", result[2])
    }
}

func TestReducesMissingNumbersIfNumIsPresent(t *testing.T) {
    input := [][]int {[]int{1}, []int{2,1,3}, []int{2,1,3}}
    result := Solution(input)

    if !IsExactly(result[0], []int{1}) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[0])
    }

    if !HasAllOf(result[1], []int {2,3}) {
        t.Errorf("Undecidable value should have all possible values. Expected {2,3}, but was %v\n", result[1])
    }
    if !HasAllOf(result[1], []int {2,3}) {
        t.Errorf("Undecidable value should have all possible values. Expected {2,3}, but was %v\n", result[2])
    }
}

func TestDoesNotIntroduceNewNumbers(t *testing.T) {
    input := [][]int {[]int{1}, []int{2,3}, []int{2,3,4}, []int {2,3,4}}
    result := Solution(input)
    if !IsExactly(result[1], []int{2,3}) {
        t.Errorf("Added some new numbers when it should not have. Expected []int{2,3}, but was %v\n", result[1])
    }
}
