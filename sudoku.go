package sudoku

import (
    "fmt"
    "reflect"
    "math"
)

type Board [][][]int

// Satisfy the Equalable interface so we can use matchers in the test.
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

// Insert into a set (row/column/square) any of the values which aren't known to be in that set.
func mapMissingValues(board [][]int) []int {
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

// Given two sets of ints, calculate the union
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

// For any value which appears in exactly one cell in a set, remove all other values from that cell
func IsolateSingletons(board [][]int) [][]int {
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

// Populate all cells with the '0' value with a full range of possible values
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

// Map b[r][c] to b[c][r]
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

// Map each sub-square of the board to a row in the output.
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

// Generate the coordinates map function to support the squaresOf() function.
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


// Given a row/col/square, propogate constraints on it.
// This is the pluggable part, though plugging doesn't seem to be necessary.
func ConstrainSet(board [][]int) [][]int {
    board = NormalizeBoard(board)
    board = IsolateSingletons(board)
    // Isolate any missing values

    notFound := mapMissingValues(board)

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

func (input Board) Solve() ([][][]int) {
    input.Step(ConstrainSet)
    for i := 0; i < 1000; i = i + 1 {
        input = input.Step(ConstrainSet)
    }
    return input
}
