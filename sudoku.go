package sudoku

import (
    "fmt"
    "math"
)


// Insert into a set (row/column/square) any of the values which aren't known to be in that set.
func findMissingValues(set Set) Cell {
    found := make(Cell, len(set))
    for _, v := range set {
        if len(v) == 1 {
            found[v[0]-1] = v[0]
        }
    }

    notFound := make(Cell, 0)
    for i, v := range found {
        if v == 0 {
            notFound = append(notFound, i + 1)
        }
    }
    return notFound
}

// This gets rewritten when Cell is rewritten as bit-values on an int
func oneOffsetComplementOf(orig []int, max int) []int {
    compl := []int{}
    for i := 1; i <= max; i++ {
        found := false
        for _, v := range orig {
            if i == v {
                found = true
            }
        }
        if !found {
            compl = append(compl, i)
        }
    }
    return compl
}

func zeroOffsetComplementOf(orig []int, max int) []int {
    compl := []int{}
    for i := 0; i < max; i++ {
        for _, v := range orig {
            if i != v {
                compl = append(compl, i)
            }
        }
    }
    return compl
}

func constrainForSet(s Set, indexesInComplement []int, constrained int) Set {
    for _, index := range indexesInComplement {
        s[index] = s[index].remove(constrained)
    }
    return s
}

func findMissingFor(input Set, onlyIn []int) []int {
    max := len(input)
    found := C()
    for _, v := range onlyIn {
        found = found.union(input[v])
    }
    missing := oneOffsetComplementOf(found, max)
    return missing
}

func ConstrainLinearAndSquare(input []Set, intersection []int) []Set {
    indexesInComplement := zeroOffsetComplementOf(intersection, len(input))
    constraineds := findMissingFor(input[1], intersection)
    for _, constrained := range constraineds {
        input[0] = constrainForSet(input[0], indexesInComplement, constrained)
    }
    return input
}

// For any value which appears in exactly one cell in a set, remove all other values from that cell
func IsolateSingletons(board Set) Set {
    singletons := make(Cell, len(board) + 1)
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
            board[singletons[i]] = make(Cell, 1)
            board[singletons[i]][0] = i
        }
    }

    return board
}

// Fill out all possible values in a cell
func Normalize(max int, cell *Cell) {
    for j := 0; j < max; j++ {
        (*cell)[j] = j + 1
    }
}

func Create(max int) Cell {
    cell := make(Cell, max)
    Normalize(max, &cell)
    return cell
}

func Copy(boardCell Cell) Cell {
    iLen := len(boardCell)
    cell := make(Cell, iLen)
    for j := 0; j < iLen; j++ {
        cell[j] = boardCell[j]
    }
    return cell
}

// Populate all cells with the '0' value with a full range of possible values
func NormalizeBoard(board Set) Set {
    max := len(board)
    outBoard := make(Set, max)
    for i := range board {
        if len(board[i]) == 0 {
            outBoard[i] = Create(max)
        } else {
            outBoard[i] = Copy(board[i])
        }
    }
    return outBoard
}

// Map b[r][c] to b[c][r]
func columnsOf(board Board) Board {
    output := make(Board, len(board[0]))
    for i := range output {
        output[i] = make(Set, len(board))
    }
    for i := range board {
        for j := range board[i] {
            output[j][i] = board[i][j]
        }
    }
    return output
}

// Map each sub-square of the board to a row in the output.
func squaresOf(board Board) Board {
    output := make(Board, len(board))
    for i := range output {
        output[i] = make(Set, len(board))
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


func (board Board) Step(filter func(Set) Set) (Board) {
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
// This is the pluggable part
func ConstrainSet(set Set) Set {
    set = NormalizeBoard(set)
    set = IsolateSingletons(set)

    // Isolate any missing values

    notFound := findMissingValues(set)

    missingValue := make(Set, len(set))
    for i, cell := range set {
        missingValue[i] = Cell{}
        if len(cell) == 0 {
            missingValue[i] = notFound
        } else if len(cell) != 1 {
            missingValue[i] = cell.intersection(notFound)
        } else {
            missingValue[i] = cell
        }
    }

    return missingValue
}

func (input Board) DebugString() string {
    out := ""
    maxWidths := make(Cell, len(input))
    for _, row := range input {
        for col, cell := range row {
            cellWidth := len(fmt.Sprintf("%v", cell))
            if maxWidths[col] < cellWidth {
                maxWidths[col] = cellWidth
            }
        }
    }

    for _, row := range input {
        out += "\t"
        for col, cell := range row {
            str := fmt.Sprintf("%v", cell)
            out += fmt.Sprintf("%*s|", maxWidths[col], str)
        }
        out += "\n"
    }
    return out
}

func (input Board) GoString() string {
    out := ""
    for row, cols := range input {
        for col, cell := range cols {
            sep := " "
            if col % 3 == 2 {
                sep = "|"
            }
            if len(cell) == 1 {
                out += fmt.Sprintf("%d%s", cell[0], sep)
            } else {
                out += fmt.Sprintf(" %s", sep)
            }
        }
        if row % 3 == 2 {
            out += fmt.Sprintf("\n------------------\n")
        } else {
            out += "\n"
        }
    }
    return out
}

func (input Board) IsSolved() bool {
    for _, row := range input {
        for _, cell := range row {
            if !cell.IsSolved() {
                return false
            }
        }
    }
    return true
}

func (input Board) Solve() (Board) {
    input.Step(ConstrainSet)
    for !input.IsSolved() {
        input = input.Step(ConstrainSet)
        fmt.Printf("Board: \n%#v\n", input)
    }
    return input
}
