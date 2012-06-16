package sudoku

import (
    "testing"
    "math"
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

func IsolateSingletons(board [][]int) [][]int {
    // Isolate singletons...
    // EFFING Magic Numbers!!!
    singletons := make([]int, 9)
    for i := range singletons {
        singletons[i] = -1
    }
    for i := range board {
        for _, v := range board[i] {
            if singletons[v] == -1 {
                singletons[v] = i
            } else {
                singletons[v] = -2
            }
        }
    }

    for i := range singletons {
        if singletons[i] > -1 {
            board[singletons[i]] = make([]int, 1)
            board[singletons[i]][0] = i
        }
    }

    return board
}

func NormalizeBoard(board [][]int) [][]int {
    max := len(board)
    outBoard := make([][]int, max)
    for i := range board {
        if len(board[i]) == 0 {
            outBoard[i] = make([]int, max)
            for j := 0; j < max; j++ {
                outBoard[i][j] = j + 1
            }
        } else {
            iLen := len(board[i])
            outBoard[i] = make([]int, iLen)
            for j := 0; j < iLen; j++ {
                outBoard[i][j] = board[i][j]
            }
        }
    }
    return outBoard
}

func ConstrainSet(board [][]int) [][]int {
    board = NormalizeBoard(board)
    board = IsolateSingletons(board)
    // Isolate any missing values

    notFound := MapMissingValues(board)

    missingValue := make([][]int, len(board))
    for i, cell := range board {
        missingValue[i] = make([]int, 0)
        if len(cell) == 0 {
            missingValue[i] = notFound
        } else if len(cell) != 1 {
            missingValue[i] = union(cell, notFound)
        } else {
            missingValue[i] = cell
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
    if len(ConstrainSet(emptyBoard)) != 0 {
        t.Fail()
    }
}

func TestOneSizedBoard(t *testing.T) {
    oneSized := [][]int {[]int{1}}
    result := ConstrainSet(oneSized)

    if len(result) != 1 {
        t.Errorf("Length should be 1, but was %v\n", result)
    }

    if !IsExactly(result[0], []int{1}) {
        t.Errorf("Single value in result should be 1, but was %v\n", result[0])
    }
}

func TestFindsMissingNumbersInList(t *testing.T) {
    input := [][]int {[]int{1}, []int{}}
    result := ConstrainSet(input)

    if !IsExactly(result[0], []int{1}) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[0])
    }

    if !IsExactly(result[1], []int{2}) {
        t.Errorf("Unknown value was not found. Expected 2, but was %v\n", result[1])
    }
}

func TestFindsMissingNumbersInList2(t *testing.T) {
    input := [][]int {[]int{}, []int{1}}
    result := ConstrainSet(input)

    if !IsExactly(result[1], []int{1}) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[1])
    }

    if !IsExactly(result[0], []int{2}) {
        t.Errorf("Unknown value was not found. Expected 2, but was %v\n", result[0])
    }
}

func TestFindsMissingNumberInSize3List(t *testing.T) {
    input := [][]int {[]int{1}, []int{2}, []int{}}
    result := ConstrainSet(input)
    if !IsExactly(result[0], []int{1}) || !IsExactly(result[1], []int{2}) {
        t.Errorf("Known values should not be changed.\n")
    }

    if !IsExactly(result[2], []int{3}) {
        t.Errorf("Unkonwn value was not found. Expected 3, but was %v\n", result[2])
    }

}

func TestFindsMultipleMissingNumbersInLongerList(t *testing.T) {
    input := [][]int {[]int{1}, []int{}, []int{}}
    result := ConstrainSet(input)

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
    result := ConstrainSet(input)

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
    result := ConstrainSet(input)
    if !IsExactly(result[1], []int{2,3}) {
        t.Errorf("Added some new numbers when it should not have. Expected []int{2,3}, but was %v\n", result[1])
    }
}

func TestIsolatesANumberWhichOnlyAppearsOnce(t *testing.T) {
    input := [][]int{[]int{1,2}, []int{1,2}, []int{1,2,3}}
    result := ConstrainSet(input)
    if !IsExactly(result[2], []int{3}) {
        t.Errorf("A number which appears exactly once should be the only possible number for that cell. Expected [3], but got %v\n", result[2])
    }
}

func TestIsolatesANumberWhichOnlyAppearsOnceAndDoesNotFallForStupidTricks(t *testing.T) {
    input := [][]int{[]int{1,2}, []int{1,2}, []int{1,2,3}, []int{}}
    result := ConstrainSet(input)
    if IsExactly(result[2], []int{3}) {
        t.Errorf("An empty cell should be replaced with all possible missing values.")
    }
}

func columnsOf(board [][][]int) [][][]int {
    output := make([][][]int, len(board[0]))
    for i := range output {
        output[i] = make([][]int, len(board))
    }
    for i := range board {
        for j := range board[i] {
            output[j][i] = board[i][j]
        }
    }
    return output
}

func squaresOf(board [][][]int) [][][]int {
    output := make([][][]int, len(board))
    for i := range output {
        output[i] = make([][]int, len(board))
    }
    coords := coordsMapForBoardOfLength(len(board))
    for i := range board {
        for j := range board[i] {
            outI, outJ := coords(i, j)
            output[outI][outJ] = board[i][j]
        }
    }
    return output
}

func coordsMapForBoardOfLength(length int) (func(int,int) (int,int)) {
    fLen := float64(length)
    squareSize := int(math.Floor(math.Sqrt(fLen)))
    if squareSize * squareSize != length {
        return func(i,j int) (int,int) {
            return 0,0
        }
    }
    return func(i, j int) (int, int) {
        return ((i/squareSize)*squareSize + j/squareSize), (j%squareSize + (i%squareSize)*squareSize)
    }
}


func TestDegenerateCoords3By3MapTo1By1Squares(t *testing.T) {
    data := [][]int{
        []int{0,0,0,0},
        []int{0,1,0,0},
        []int{1,0,0,0},
        []int{2,2,0,0},
    }

    f := coordsMapForBoardOfLength(3)
    for _, datum := range data {
        i, j := f(datum[0], datum[1])
        if i != datum[2] || j != datum[3] {
            t.Errorf("For input %v,%v, expected %v,%v but got %v,%v\n", datum[0],datum[1], datum[2], datum[3], i, j)
        }
    }
}

func TestCoords9By9MapTo3By3Squares(t *testing.T) {
    data := [][]int{
        []int{0,0,0,0},
        []int{1,0,0,3},
        []int{1,1,0,4},
        []int{2,2,0,8},
        []int{0,6,2,0},
        []int{0,7,2,1},
        []int{0,8,2,2},
        []int{1,6,2,3},
        []int{2,6,2,6},
        []int{3,0,3,0},
        []int{3,2,3,2},
        []int{4,0,3,3},
        []int{3,6,5,0},
        []int{4,6,5,3},
        []int{5,8,5,8},
    }

    f := coordsMapForBoardOfLength(9)
    for _, datum := range data {
        i, j := f(datum[0], datum[1])
        if i != datum[2] || j != datum[3] {
            t.Errorf("For input %v,%v, expected %v,%v but got %v,%v\n", datum[0],datum[1], datum[2], datum[3], i, j)
        }
    }
}

func _TestSauaresOf81CellBoardAre3X3Nondrants(t *testing.T) {
    input := [][][]int{
        [][]int{[]int{1},[]int{1},[]int{1},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2}},
        [][]int{[]int{1},[]int{1},[]int{1},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2}},
        [][]int{[]int{1},[]int{1},[]int{1},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2}},
        [][]int{[]int{4},[]int{4},[]int{4},[]int{5},[]int{5},[]int{5},[]int{6},[]int{6},[]int{6}},
        [][]int{[]int{4},[]int{4},[]int{4},[]int{5},[]int{5},[]int{5},[]int{6},[]int{6},[]int{6}},
        [][]int{[]int{4},[]int{4},[]int{4},[]int{5},[]int{5},[]int{5},[]int{6},[]int{6},[]int{6}},
        [][]int{[]int{7},[]int{7},[]int{7},[]int{8},[]int{8},[]int{8},[]int{9},[]int{9},[]int{9}},
        [][]int{[]int{7},[]int{7},[]int{7},[]int{8},[]int{8},[]int{8},[]int{9},[]int{9},[]int{9}},
        [][]int{[]int{7},[]int{7},[]int{7},[]int{8},[]int{8},[]int{8},[]int{9},[]int{9},[]int{9}},
    }

    expected := [][][]int{
        [][]int{[]int{1},[]int{1},[]int{1},[]int{1},[]int{1},[]int{1},[]int{1},[]int{1},[]int{1},},
        [][]int{[]int{2},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2},[]int{2},},
        [][]int{[]int{3},[]int{3},[]int{3},[]int{3},[]int{3},[]int{3},[]int{3},[]int{3},[]int{3},},
        [][]int{[]int{4},[]int{4},[]int{4},[]int{4},[]int{4},[]int{4},[]int{4},[]int{4},[]int{4},},
        [][]int{[]int{5},[]int{5},[]int{5},[]int{5},[]int{5},[]int{5},[]int{5},[]int{5},[]int{5},},
        [][]int{[]int{6},[]int{6},[]int{6},[]int{6},[]int{6},[]int{6},[]int{6},[]int{6},[]int{6},},
        [][]int{[]int{7},[]int{7},[]int{7},[]int{7},[]int{7},[]int{7},[]int{7},[]int{7},[]int{7},},
        [][]int{[]int{8},[]int{8},[]int{8},[]int{8},[]int{8},[]int{8},[]int{8},[]int{8},[]int{8},},
        [][]int{[]int{9},[]int{9},[]int{9},[]int{9},[]int{9},[]int{9},[]int{9},[]int{9},[]int{9},},
    }

    validateSameCells(t, expected, squaresOf(input))
}

func TestColumnsOf(t *testing.T) {
    input := [][][]int{
            [][]int{[]int{1},[]int{2},[]int{3}},
            [][]int{[]int{1},[]int{2},[]int{3}},
            [][]int{[]int{1},[]int{2},[]int{3}},
            }

    expected := [][][]int{
            [][]int{[]int{1},[]int{1},[]int{1}},
            [][]int{[]int{2},[]int{2},[]int{2}},
            [][]int{[]int{3},[]int{3},[]int{3}},
            }
    validateSameCells(t, expected, columnsOf(input))
}

func Step(board [][][]int, filter func([][]int)) {
    for i := range board {
        filter(board[i])
    }
    cols := columnsOf(board)
    for _, col := range cols {
        filter(col)
    }
    squares := squaresOf(board)
    for _, square := range squares {
        filter(square)
    }
}

func TestCallsFunction(t *testing.T) {
    wasCalled := false
    input := [][][]int{[][]int{[]int{1}}}
    Step(input, func(board [][]int) { wasCalled = true })
    if !wasCalled {
        t.Errorf("Expected function to be called by Step(), but was not.")
    }
}

func TestCallsFunctionOnRows(t *testing.T) {
    rows := [][][]int{}
    input := [][][]int{[][]int{[]int{1},[]int{2},[]int{3}}, [][]int{[]int{1},[]int{2},[]int{3}}, [][]int{[]int{1},[]int{2},[]int{3}},
                       [][]int{}, [][]int{}, [][]int{},
                       [][]int{}, [][]int{}, [][]int{}}

    Step(input, func(board [][]int) {
        rows = append(rows, board)
    })

    validateSameCells(t, input, rows)
}

func containsData(container [][][]int, containee [][]int) bool {
    Row: for i := range container {
        for j := range container[i] {
            for k := range container[i][j] {
                if len(containee) <= j || len(containee[j]) <= k || containee[j][k] != container[i][j][k] {
                    continue Row

                }
            }
        }
        return true
    }
    return false
}

func validateSameCells(t *testing.T, expected [][][]int, rows [][][]int) {

    for _, r := range expected {
        if !containsData(rows, r) {
            t.Errorf("Cannot find row data %v in %v", r, rows)
        }
    }
}

func TestCallsFunctionOnCols(t *testing.T) {
    cols := [][][]int{}
    input := [][][]int{
                [][]int{[]int{1},[]int{2},[]int{6}},
                [][]int{[]int{4},[]int{5},[]int{8}},
                [][]int{[]int{3},[]int{9},[]int{7}},
            }

    expected := [][][]int{
        [][]int{[]int{1},[]int{4},[]int{3}},
        [][]int{[]int{2},[]int{5},[]int{9}},
        [][]int{[]int{6},[]int{8},[]int{7}},
    }

    Step(input, func(board [][]int) {
        cols = append(cols, board)
    })

    validateSameCells(t, expected, cols)
}

func TestCallsFunctionOnSquares(t *testing.T) {
    squares := [][][]int{}
    input := [][][]int{
        [][]int{[]int{1},[]int{2},[]int{3},[]int{4},[]int{5},[]int{6},[]int{7},[]int{8},[]int{9}},
        [][]int{[]int{9},[]int{1},[]int{2},[]int{3},[]int{4},[]int{5},[]int{6},[]int{7},[]int{8}},
        [][]int{[]int{8},[]int{9},[]int{1},[]int{2},[]int{3},[]int{4},[]int{5},[]int{6},[]int{7}},
        [][]int{[]int{7},[]int{8},[]int{9},[]int{1},[]int{2},[]int{3},[]int{4},[]int{5},[]int{6}},
        [][]int{[]int{6},[]int{7},[]int{8},[]int{9},[]int{1},[]int{2},[]int{3},[]int{4},[]int{5}},
        [][]int{[]int{5},[]int{6},[]int{7},[]int{8},[]int{9},[]int{1},[]int{2},[]int{3},[]int{4}},
        [][]int{[]int{4},[]int{5},[]int{6},[]int{7},[]int{8},[]int{9},[]int{1},[]int{2},[]int{3}},
        [][]int{[]int{3},[]int{4},[]int{5},[]int{6},[]int{7},[]int{8},[]int{9},[]int{1},[]int{2}},
        [][]int{[]int{2},[]int{3},[]int{4},[]int{5},[]int{6},[]int{7},[]int{8},[]int{9},[]int{1}},
    }

    expected := [][][]int{
        [][]int{[]int{1},[]int{2},[]int{3},[]int{9},[]int{1},[]int{2},[]int{8},[]int{9},[]int{1}},
        [][]int{[]int{4},[]int{5},[]int{6},[]int{3},[]int{4},[]int{5},[]int{2},[]int{3},[]int{4}},
        [][]int{[]int{7},[]int{8},[]int{9},[]int{6},[]int{7},[]int{8},[]int{5},[]int{6},[]int{7}},
        [][]int{[]int{7},[]int{8},[]int{9},[]int{6},[]int{7},[]int{8},[]int{5},[]int{6},[]int{7}},
        [][]int{[]int{1},[]int{2},[]int{3},[]int{9},[]int{1},[]int{2},[]int{8},[]int{9},[]int{1}},
        [][]int{[]int{4},[]int{5},[]int{6},[]int{3},[]int{4},[]int{5},[]int{2},[]int{3},[]int{4}},
        [][]int{[]int{4},[]int{5},[]int{6},[]int{3},[]int{4},[]int{5},[]int{2},[]int{3},[]int{4}},
        [][]int{[]int{7},[]int{8},[]int{9},[]int{6},[]int{7},[]int{8},[]int{5},[]int{6},[]int{7}},
        [][]int{[]int{1},[]int{2},[]int{3},[]int{9},[]int{1},[]int{2},[]int{8},[]int{9},[]int{1}},
    }

    Step(input, func(board [][]int) {
        squares = append(squares, board)
    })

    validateSameCells(t, expected, squares)
}

/*
func TestSolvesThis(t *testing.T) {
    input := [][][]int{
        [][]int{[]int{},[]int{1},[]int{},[]int{6},[]int{},[]int{7},[]int{},[]int{},[]int{4}},
        [][]int{[]int{},[]int{4},[]int{2},[]int{},[]int{},[]int{},[]int{},[]int{},[]int{}},
        [][]int{[]int{8},[]int{7},[]int{},[]int{3},[]int{},[]int{},[]int{6},[]int{},[]int{}},
        [][]int{[]int{},[]int{8},[]int{},[]int{},[]int{7},[]int{},[]int{},[]int{2},[]int{}},
        [][]int{[]int{},[]int{},[]int{},[]int{8},[]int{9},[]int{3},[]int{},[]int{},[]int{}},
        [][]int{[]int{},[]int{3},[]int{},[]int{},[]int{6},[]int{},[]int{},[]int{1},[]int{}},
        [][]int{[]int{},[]int{},[]int{8},[]int{},[]int{},[]int{6},[]int{},[]int{4},[]int{5}},
        [][]int{[]int{},[]int{},[]int{},[]int{},[]int{},[]int{},[]int{1},[]int{7},[]int{}},
        [][]int{[]int{4},[]int{},[]int{},[]int{9},[]int{},[]int{8},[]int{},[]int{6},[]int{}},
    }

    Step(input, func(board [][]int) {
        input = ConstrainSet(input)
    })
}
*/
