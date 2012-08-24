package sudoku

import (
    "testing"
    "math"
    "fmt"
    "matchers"
    "reflect"
)

type Board [][][]int

func (b Board) Equals(other interface{}) (bool, string) {
    switch o := other.(type) {
        case [][][]int:
            boardLength := len(b)
            boardWidth := 0
            if boardLength > 0 {
                boardWidth = len((b)[0])
            }

            otherLength := len((o))
            otherWidth := 0
            if otherLength > 0 {
                otherWidth = len((o)[0])
            }
            if boardLength != otherLength || boardWidth != otherWidth {
                return false, fmt.Sprintf("mismatch between %v-by-%v and %v-by-%v boards", boardLength, boardWidth, otherLength, otherWidth)
            }
            equals := true
            msg := fmt.Sprintf("\n")
            for i := range b {
                if boardWidth != len((b)[i]) || otherWidth != len((o)[i]) {
                    return false, fmt.Sprintf("the board (or other board) is not of equal widths!")
                }
                for j := range b[i] {
                    sameLen := len(b[i][j]) == len(o[i][j])
                    if len(b[i][j]) == 0 {
                        equals = equals && sameLen
                        msg += "| X "
                    } else if len(o[i][j]) == 0 {
                        equals = equals && sameLen
                        msg += "| # "
                    } else if b[i][j][0] != o[i][j][0] {
                        equals = false
                        msg += fmt.Sprintf("|%v %v", b[i][j][0], o[i][j][0])
                    } else {
                        msg += fmt.Sprintf("|   ")
                    }
                }
                msg += fmt.Sprintf("\n")
            }

            return equals, msg
    }
    return false, fmt.Sprintf("a Board cannot equal a %v", reflect.TypeOf(other))
}

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
    // EFFING Magic Numbers!!!
    singletons := make([]int, len(board) + 1)
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

func TestNormalizeEmptyBoard(t *testing.T) {
    input := [][]int{[]int{}, []int{}}
    expected := [][]int{[]int{1,2}, []int{1,2}}

    actual := NormalizeBoard(input)

    for i := range actual {
        for j := range actual[i] {
            if actual[i][j] != expected[i][j] {
                t.Errorf("At %v,%v actual %v but expected %v", i, j, actual[i][j], expected[i][j])
            }
        }
    }
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
            if len(output) <= outI || len(output[outI]) <= outJ {
                fmt.Printf("Failed at output[%v][%v]: %v, %v\n", outI, outJ, i, j)
            }
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

func (board Board) Step(filter func([][]int) [][]int) (Board) {
    for i := range board {
        board[i] = filter(board[i])
    }

    cols := columnsOf(board)
    for i, col := range cols {
        updatedCol := filter(col)
        for j := range updatedCol {
            board[j][i] = updatedCol[j]
        }
    }

    squares := squaresOf(board)
    for i, square := range squares {
        squares[i] = filter(square)
    }
    mapper := coordsMapForBoardOfLength(len(board))
    for i := range board {
        for j := range board[i] {
            squareI, squareJ := mapper(i, j)
            board[i][j] = squares[squareI][squareJ]
        }
    }
    return board
}

func TestCallsFunction(t *testing.T) {
    wasCalled := false
    input := Board{[][]int{[]int{1}}}
    input.Step(func(board [][]int) [][]int { wasCalled = true; return [][]int{};  })
    if !wasCalled {
        t.Errorf("Expected function to be called by Step(), but was not.")
    }
}

func TestCallsFunctionOnRows(t *testing.T) {
    rows := Board{}
    input := Board{[][]int{[]int{1},[]int{2},[]int{3}}, [][]int{[]int{1},[]int{2},[]int{3}}, [][]int{[]int{1},[]int{2},[]int{3}},
                       [][]int{}, [][]int{}, [][]int{},
                       [][]int{}, [][]int{}, [][]int{}}

    input.Step(func(board [][]int) [][]int {
        rows = append(rows, board)
        return [][]int{};
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
    cols := Board{}
    input := Board{
                [][]int{[]int{1},[]int{2},[]int{6}},
                [][]int{[]int{4},[]int{5},[]int{8}},
                [][]int{[]int{3},[]int{9},[]int{7}},
            }

    expected := Board{
        [][]int{[]int{1},[]int{4},[]int{3}},
        [][]int{[]int{2},[]int{5},[]int{9}},
        [][]int{[]int{6},[]int{8},[]int{7}},
    }

    input.Step(func(board [][]int) [][]int {
        cols = append(cols, board)
        return [][]int{}
    })

    validateSameCells(t, expected, cols)
}

func TestCallsFunctionOnSquares(t *testing.T) {
    squares := Board{}
    input := Board{
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

    expected := Board{
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

    input.Step(func(board [][]int) [][]int {
        squares = append(squares, board)
        return [][]int{}
    })

    validateSameCells(t, expected, squares)
}

func (input Board) Solve() ([][][]int) {
    input.Step(ConstrainSet)
    for i := 0; i < 1000; i = i + 1 {
        input = input.Step(ConstrainSet)
    }
    return input
}

var unsolved Board = Board{
    [][]int{[]int{ },[]int{1},[]int{ },[]int{6},[]int{ },[]int{7},[]int{ },[]int{ },[]int{4}},
    [][]int{[]int{ },[]int{4},[]int{2},[]int{ },[]int{ },[]int{ },[]int{ },[]int{ },[]int{ }},
    [][]int{[]int{8},[]int{7},[]int{ },[]int{3},[]int{ },[]int{ },[]int{6},[]int{ },[]int{ }},
    [][]int{[]int{ },[]int{8},[]int{ },[]int{ },[]int{7},[]int{ },[]int{ },[]int{2},[]int{ }},
    [][]int{[]int{ },[]int{ },[]int{ },[]int{8},[]int{9},[]int{3},[]int{ },[]int{ },[]int{ }},
    [][]int{[]int{ },[]int{3},[]int{ },[]int{ },[]int{6},[]int{ },[]int{ },[]int{1},[]int{ }},
    [][]int{[]int{ },[]int{ },[]int{8},[]int{ },[]int{ },[]int{6},[]int{ },[]int{4},[]int{5}},
    [][]int{[]int{ },[]int{ },[]int{ },[]int{ },[]int{ },[]int{ },[]int{1},[]int{7},[]int{ }},
    [][]int{[]int{4},[]int{ },[]int{ },[]int{9},[]int{ },[]int{8},[]int{ },[]int{6},[]int{ }},
}
var solved Board = Board{
    [][]int{[]int{9},[]int{1},[]int{3},[]int{6},[]int{2},[]int{7},[]int{5},[]int{8},[]int{4}},
    [][]int{[]int{6},[]int{4},[]int{2},[]int{5},[]int{8},[]int{9},[]int{7},[]int{3},[]int{1}},
    [][]int{[]int{8},[]int{7},[]int{5},[]int{3},[]int{4},[]int{1},[]int{6},[]int{9},[]int{2}},
    [][]int{[]int{5},[]int{8},[]int{9},[]int{1},[]int{7},[]int{4},[]int{3},[]int{2},[]int{6}},
    [][]int{[]int{2},[]int{6},[]int{1},[]int{8},[]int{9},[]int{3},[]int{4},[]int{5},[]int{7}},
    [][]int{[]int{7},[]int{3},[]int{4},[]int{2},[]int{6},[]int{5},[]int{8},[]int{1},[]int{9}},
    [][]int{[]int{1},[]int{2},[]int{8},[]int{7},[]int{3},[]int{6},[]int{9},[]int{4},[]int{5}},
    [][]int{[]int{3},[]int{9},[]int{6},[]int{4},[]int{5},[]int{2},[]int{1},[]int{7},[]int{8}},
    [][]int{[]int{4},[]int{5},[]int{7},[]int{9},[]int{1},[]int{8},[]int{2},[]int{6},[]int{3}},
}

func TestStepThroughIt(t *testing.T) {
    unsolved.Step(ConstrainSet)
}

func TestSolvesThis(t *testing.T) {
    input := Board{
        [][]int{[]int{ },[]int{1},[]int{ },[]int{6},[]int{ },[]int{7},[]int{ },[]int{ },[]int{4}},
        [][]int{[]int{ },[]int{4},[]int{2},[]int{ },[]int{ },[]int{ },[]int{ },[]int{ },[]int{ }},
        [][]int{[]int{8},[]int{7},[]int{ },[]int{3},[]int{ },[]int{ },[]int{6},[]int{ },[]int{ }},
        [][]int{[]int{ },[]int{8},[]int{ },[]int{ },[]int{7},[]int{ },[]int{ },[]int{2},[]int{ }},
        [][]int{[]int{ },[]int{ },[]int{ },[]int{8},[]int{9},[]int{3},[]int{ },[]int{ },[]int{ }},
        [][]int{[]int{ },[]int{3},[]int{ },[]int{ },[]int{6},[]int{ },[]int{ },[]int{1},[]int{ }},
        [][]int{[]int{ },[]int{ },[]int{8},[]int{ },[]int{ },[]int{6},[]int{ },[]int{4},[]int{5}},
        [][]int{[]int{ },[]int{ },[]int{ },[]int{ },[]int{ },[]int{ },[]int{1},[]int{7},[]int{ }},
        [][]int{[]int{4},[]int{ },[]int{ },[]int{9},[]int{ },[]int{8},[]int{ },[]int{6},[]int{ }},
    }

    expected := Board{
        [][]int{[]int{9},[]int{1},[]int{3},[]int{6},[]int{2},[]int{7},[]int{5},[]int{8},[]int{4}},
        [][]int{[]int{6},[]int{4},[]int{2},[]int{5},[]int{8},[]int{9},[]int{7},[]int{3},[]int{1}},
        [][]int{[]int{8},[]int{7},[]int{5},[]int{3},[]int{4},[]int{1},[]int{6},[]int{9},[]int{2}},
        [][]int{[]int{5},[]int{8},[]int{9},[]int{1},[]int{7},[]int{4},[]int{3},[]int{2},[]int{6}},
        [][]int{[]int{2},[]int{6},[]int{1},[]int{8},[]int{9},[]int{3},[]int{4},[]int{5},[]int{7}},
        [][]int{[]int{7},[]int{3},[]int{4},[]int{2},[]int{6},[]int{5},[]int{8},[]int{1},[]int{9}},
        [][]int{[]int{1},[]int{2},[]int{8},[]int{7},[]int{3},[]int{6},[]int{9},[]int{4},[]int{5}},
        [][]int{[]int{3},[]int{9},[]int{6},[]int{4},[]int{5},[]int{2},[]int{1},[]int{7},[]int{8}},
        [][]int{[]int{4},[]int{5},[]int{7},[]int{9},[]int{1},[]int{8},[]int{2},[]int{6},[]int{3}},
    }
    output := input.Solve()

    matchers.AssertThat(t, output, matchers.Equals(expected))
}
