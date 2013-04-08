package sudoku

import (
    matchers "github.com/tychofreeman/go-matchers"
    "testing"
    "fmt"
)

func containsData(container Board, containee Set) bool {
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

func validateSameCells(t *testing.T, expected Board, rows Board) {

    for _, r := range expected {
        if !containsData(rows, r) {
            t.Errorf("Cannot find row data %v in %v", r, rows)
        }
    }
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

func TestNormalizeEmptyBoard(t *testing.T) {
    input := Set{C(), C()}
    expected := Set{C(1,2), C(1,2)}

    actual := NormalizeBoard(input)

    for i := range actual {
        for j := range actual[i] {
            if actual[i][j] != expected[i][j] {
                t.Errorf("At %v,%v actual %v but expected %v", i, j, actual[i][j], expected[i][j])
            }
        }
    }
}

func TestZeroSizedBoard(t *testing.T) {
    emptyBoard := Set {}
    if len(ConstrainSet(emptyBoard)) != 0 {
        t.Fail()
    }
}

func TestOneSizedBoard(t *testing.T) {
    oneSized := Set {C(1)}
    result := ConstrainSet(oneSized)

    if len(result) != 1 {
        t.Errorf("Length should be 1, but was %v\n", result)
    }

    if !IsExactly(result[0], C(1)) {
        t.Errorf("Single value in result should be 1, but was %v\n", result[0])
    }
}

func TestFindsMissingNumbersInList(t *testing.T) {
    input := Set {C(1), C()}
    result := ConstrainSet(input)

    if !IsExactly(result[0], C(1)) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[0])
    }

    if !IsExactly(result[1], C(2)) {
        t.Errorf("Unknown value was not found. Expected 2, but was %v\n", result[1])
    }
}

func TestFindsMissingNumbersInList2(t *testing.T) {
    input := Set {C(), C(1)}
    result := ConstrainSet(input)

    if !IsExactly(result[1], C(1)) {
        t.Errorf("Known value should not be changed. Expected 1, but was %v\n", result[1])
    }

    if !IsExactly(result[0], C(2)) {
        t.Errorf("Unknown value was not found. Expected 2, but was %v\n", result[0])
    }
}

func TestFindsMissingNumberInSize3List(t *testing.T) {
    input := Set {C(1), C(2), C()}
    result := ConstrainSet(input)
    if !IsExactly(result[0], C(1)) || !IsExactly(result[1], C(2)) {
        t.Errorf("Known values should not be changed.\n")
    }

    if !IsExactly(result[2], C(3)) {
        t.Errorf("Unkonwn value was not found. Expected 3, but was %v\n", result[2])
    }

}

func TestFindsMultipleMissingNumbersInLongerList(t *testing.T) {
    input := Set {C(1), C(), C()}
    result := ConstrainSet(input)

    if !IsExactly(result[0], C(1)) {
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
    input := Set {C(1), C(2,1,3), C(2,1,3)}
    result := ConstrainSet(input)

    if !IsExactly(result[0], C(1)) {
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
    input := Set {C(1), C(2,3), C(2,3,4), C(2,3,4)}
    result := ConstrainSet(input)
    if !IsExactly(result[1], C(2,3)) {
        t.Errorf("Added some new numbers when it should not have. Expected C(2,3), but was %v\n", result[1])
    }
}

func TestIsolatesANumberWhichOnlyAppearsOnce(t *testing.T) {
    input := Set{C(1,2), C(1,2), C(1,2,3)}
    result := ConstrainSet(input)
    if !IsExactly(result[2], C(3)) {
        t.Errorf("A number which appears exactly once should be the only possible number for that cell. Expected [3], but got %v\n", result[2])
    }
}

func TestIsolatesANumberWhichOnlyAppearsOnceAndDoesNotFallForStupidTricks(t *testing.T) {
    input := Set{C(1,2), C(1,2), C(1,2,3), C()}
    result := ConstrainSet(input)
    if IsExactly(result[2], C(3)) {
        t.Errorf("An empty cell should be replaced with all possible missing values.")
    }
}

func TestDegenerateCoords3By3MapTo1By1Squares(t *testing.T) {
    data := Set{
        C(0,0,0,0),
        C(0,1,0,0),
        C(1,0,0,0),
        C(2,2,0,0),
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
    data := Set{
        C(0,0,0,0),
        C(1,0,0,3),
        C(1,1,0,4),
        C(2,2,0,8),
        C(0,6,2,0),
        C(0,7,2,1),
        C(0,8,2,2),
        C(1,6,2,3),
        C(2,6,2,6),
        C(3,0,3,0),
        C(3,2,3,2),
        C(4,0,3,3),
        C(3,6,5,0),
        C(4,6,5,3),
        C(5,8,5,8),
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
    input := Board{
        Set{C(1),C(1),C(1),C(2),C(2),C(2),C(2),C(2),C(2)},
        Set{C(1),C(1),C(1),C(2),C(2),C(2),C(2),C(2),C(2)},
        Set{C(1),C(1),C(1),C(2),C(2),C(2),C(2),C(2),C(2)},
        Set{C(4),C(4),C(4),C(5),C(5),C(5),C(6),C(6),C(6)},
        Set{C(4),C(4),C(4),C(5),C(5),C(5),C(6),C(6),C(6)},
        Set{C(4),C(4),C(4),C(5),C(5),C(5),C(6),C(6),C(6)},
        Set{C(7),C(7),C(7),C(8),C(8),C(8),C(9),C(9),C(9)},
        Set{C(7),C(7),C(7),C(8),C(8),C(8),C(9),C(9),C(9)},
        Set{C(7),C(7),C(7),C(8),C(8),C(8),C(9),C(9),C(9)},
    }

    expected := Board{
        Set{C(1),C(1),C(1),C(1),C(1),C(1),C(1),C(1),C(1),},
        Set{C(2),C(2),C(2),C(2),C(2),C(2),C(2),C(2),C(2),},
        Set{C(3),C(3),C(3),C(3),C(3),C(3),C(3),C(3),C(3),},
        Set{C(4),C(4),C(4),C(4),C(4),C(4),C(4),C(4),C(4),},
        Set{C(5),C(5),C(5),C(5),C(5),C(5),C(5),C(5),C(5),},
        Set{C(6),C(6),C(6),C(6),C(6),C(6),C(6),C(6),C(6),},
        Set{C(7),C(7),C(7),C(7),C(7),C(7),C(7),C(7),C(7),},
        Set{C(8),C(8),C(8),C(8),C(8),C(8),C(8),C(8),C(8),},
        Set{C(9),C(9),C(9),C(9),C(9),C(9),C(9),C(9),C(9),},
    }

    validateSameCells(t, expected, squaresOf(input))
}

func TestColumnsOf(t *testing.T) {
    input := Board{
            Set{C(1),C(2),C(3)},
            Set{C(1),C(2),C(3)},
            Set{C(1),C(2),C(3)},
            }

    expected := Board{
            Set{C(1),C(1),C(1)},
            Set{C(2),C(2),C(2)},
            Set{C(3),C(3),C(3)},
            }
    validateSameCells(t, expected, columnsOf(input))
}

func TestCallsFunction(t *testing.T) {
    wasCalled := false
    input := Board{Set{C(1)}}
    input.Step(func(board Set) Set { wasCalled = true; return Set{};  })
    if !wasCalled {
        t.Errorf("Expected function to be called by Step(), but was not.")
    }
}

func TestCallsFunctionOnRows(t *testing.T) {
    rows := Board{}
    input := Board{Set{C(1),C(2),C(3)}, Set{C(1),C(2),C(3)}, Set{C(1),C(2),C(3)},
                       Set{}, Set{}, Set{},
                       Set{}, Set{}, Set{}}

    input.Step(func(board Set) Set {
        rows = append(rows, board)
        return Set{};
    })

    validateSameCells(t, input, rows)
}

func TestCallsFunctionOnCols(t *testing.T) {
    cols := Board{}
    input := Board{
                Set{C(1),C(2),C(6)},
                Set{C(4),C(5),C(8)},
                Set{C(3),C(9),C(7)},
            }

    expected := Board{
        Set{C(1),C(4),C(3)},
        Set{C(2),C(5),C(9)},
        Set{C(6),C(8),C(7)},
    }

    input.Step(func(board Set) Set {
        cols = append(cols, board)
        return Set{}
    })

    validateSameCells(t, expected, cols)
}

func TestCallsFunctionOnSquares(t *testing.T) {
    squares := Board{}
    input := Board{
        Set{C(1),C(2),C(3),C(4),C(5),C(6),C(7),C(8),C(9)},
        Set{C(9),C(1),C(2),C(3),C(4),C(5),C(6),C(7),C(8)},
        Set{C(8),C(9),C(1),C(2),C(3),C(4),C(5),C(6),C(7)},
        Set{C(7),C(8),C(9),C(1),C(2),C(3),C(4),C(5),C(6)},
        Set{C(6),C(7),C(8),C(9),C(1),C(2),C(3),C(4),C(5)},
        Set{C(5),C(6),C(7),C(8),C(9),C(1),C(2),C(3),C(4)},
        Set{C(4),C(5),C(6),C(7),C(8),C(9),C(1),C(2),C(3)},
        Set{C(3),C(4),C(5),C(6),C(7),C(8),C(9),C(1),C(2)},
        Set{C(2),C(3),C(4),C(5),C(6),C(7),C(8),C(9),C(1)},
    }

    expected := Board{
        Set{C(1),C(2),C(3),C(9),C(1),C(2),C(8),C(9),C(1)},
        Set{C(4),C(5),C(6),C(3),C(4),C(5),C(2),C(3),C(4)},
        Set{C(7),C(8),C(9),C(6),C(7),C(8),C(5),C(6),C(7)},
        Set{C(7),C(8),C(9),C(6),C(7),C(8),C(5),C(6),C(7)},
        Set{C(1),C(2),C(3),C(9),C(1),C(2),C(8),C(9),C(1)},
        Set{C(4),C(5),C(6),C(3),C(4),C(5),C(2),C(3),C(4)},
        Set{C(4),C(5),C(6),C(3),C(4),C(5),C(2),C(3),C(4)},
        Set{C(7),C(8),C(9),C(6),C(7),C(8),C(5),C(6),C(7)},
        Set{C(1),C(2),C(3),C(9),C(1),C(2),C(8),C(9),C(1)},
    }

    input.Step(func(board Set) Set {
        squares = append(squares, board)
        return Set{}
    })

    validateSameCells(t, expected, squares)
}

var unsolved Board = Board{
    Set{C( ),C(1),C( ),C(6),C( ),C(7),C( ),C( ),C(4)},
    Set{C( ),C(4),C(2),C( ),C( ),C( ),C( ),C( ),C( )},
    Set{C(8),C(7),C( ),C(3),C( ),C( ),C(6),C( ),C( )},
    Set{C( ),C(8),C( ),C( ),C(7),C( ),C( ),C(2),C( )},
    Set{C( ),C( ),C( ),C(8),C(9),C(3),C( ),C( ),C( )},
    Set{C( ),C(3),C( ),C( ),C(6),C( ),C( ),C(1),C( )},
    Set{C( ),C( ),C(8),C( ),C( ),C(6),C( ),C(4),C(5)},
    Set{C( ),C( ),C( ),C( ),C( ),C( ),C(1),C(7),C( )},
    Set{C(4),C( ),C( ),C(9),C( ),C(8),C( ),C(6),C( )},
}
var solved Board = Board{
    Set{C(9),C(1),C(3),C(6),C(2),C(7),C(5),C(8),C(4)},
    Set{C(6),C(4),C(2),C(5),C(8),C(9),C(7),C(3),C(1)},
    Set{C(8),C(7),C(5),C(3),C(4),C(1),C(6),C(9),C(2)},
    Set{C(5),C(8),C(9),C(1),C(7),C(4),C(3),C(2),C(6)},
    Set{C(2),C(6),C(1),C(8),C(9),C(3),C(4),C(5),C(7)},
    Set{C(7),C(3),C(4),C(2),C(6),C(5),C(8),C(1),C(9)},
    Set{C(1),C(2),C(8),C(7),C(3),C(6),C(9),C(4),C(5)},
    Set{C(3),C(9),C(6),C(4),C(5),C(2),C(1),C(7),C(8)},
    Set{C(4),C(5),C(7),C(9),C(1),C(8),C(2),C(6),C(3)},
}

func TestStepThroughIt(t *testing.T) {
    unsolved.Step(ConstrainSet)
}

func TestSolvesThis(t *testing.T) {
    input := Board{
        Set{C( ),C(1),C( ),C(6),C( ),C(7),C( ),C( ),C(4)},
        Set{C( ),C(4),C(2),C( ),C( ),C( ),C( ),C( ),C( )},
        Set{C(8),C(7),C( ),C(3),C( ),C( ),C(6),C( ),C( )},
        Set{C( ),C(8),C( ),C( ),C(7),C( ),C( ),C(2),C( )},
        Set{C( ),C( ),C( ),C(8),C(9),C(3),C( ),C( ),C( )},
        Set{C( ),C(3),C( ),C( ),C(6),C( ),C( ),C(1),C( )},
        Set{C( ),C( ),C(8),C( ),C( ),C(6),C( ),C(4),C(5)},
        Set{C( ),C( ),C( ),C( ),C( ),C( ),C(1),C(7),C( )},
        Set{C(4),C( ),C( ),C(9),C( ),C(8),C( ),C(6),C( )},
    }

    expected := Board{
        Set{C(9),C(1),C(3),C(6),C(2),C(7),C(5),C(8),C(4)},
        Set{C(6),C(4),C(2),C(5),C(8),C(9),C(7),C(3),C(1)},
        Set{C(8),C(7),C(5),C(3),C(4),C(1),C(6),C(9),C(2)},
        Set{C(5),C(8),C(9),C(1),C(7),C(4),C(3),C(2),C(6)},
        Set{C(2),C(6),C(1),C(8),C(9),C(3),C(4),C(5),C(7)},
        Set{C(7),C(3),C(4),C(2),C(6),C(5),C(8),C(1),C(9)},
        Set{C(1),C(2),C(8),C(7),C(3),C(6),C(9),C(4),C(5)},
        Set{C(3),C(9),C(6),C(4),C(5),C(2),C(1),C(7),C(8)},
        Set{C(4),C(5),C(7),C(9),C(1),C(8),C(2),C(6),C(3)},
    }
    output := input.Solve()

    matchers.AssertThat(t, output, matchers.Equals(expected))
}

func DISABLED_TestSolvesExtremePuzzle(t *testing.T) {
    input := Board{
        Set{C( ),C( ),C(5),C(6),C( ),C( ),C( ),C( ),C(7)},
        Set{C( ),C(6),C( ),C( ),C(4),C( ),C( ),C(8),C( )},
        Set{C( ),C( ),C(9),C( ),C( ),C( ),C( ),C( ),C(1)},
        Set{C(7),C( ),C( ),C( ),C( ),C( ),C(1),C( ),C( )},
        Set{C( ),C(8),C( ),C( ),C(1),C( ),C( ),C(2),C( )},
        Set{C( ),C( ),C(2),C( ),C( ),C( ),C( ),C( ),C(4)},
        Set{C(5),C( ),C( ),C( ),C( ),C( ),C(3),C( ),C( )},
        Set{C( ),C(2),C( ),C( ),C(9),C( ),C( ),C(6),C( )},
        Set{C(4),C( ),C( ),C( ),C( ),C(7),C(5),C( ),C( )},
    }

    output := input.Solve()

    matchers.AssertThat(t, output.IsSolved(), matchers.IsTrue)
}

func TestRemoves1sFromRestOfSquareWhenASubsetMustContainThem(t *testing.T) {
    input := []Set{
        Set{C(1,2),C(1,3),C(1,4),C(1,4)},
        Set{C(1,4),C(1,4),C(2  ),C(3  )},
    }
    intersect := []int{2,3}
    expected := []Set{
        Set{C(2  ),C(3  ),C(1,4),C(1,4)},
        Set{C(1,4),C(1,4),C(2  ),C(3  )},
    }

    output := ConstrainLinearAndSquare(input, intersect)

    matchers.AssertThat(t, output, matchers.Equals(expected))
}

type CellSet Set
func (cs CellSet) Equals(other interface{}) (b bool, s string) {
    b = true
    s = ""
    
    switch o := other.(type) {
    case Set:
        for i := range cs {
            found := true
            if len(cs[i]) != len(o[i]) {
                found = false
            } else {
                for j := range cs[i] {
                    inOther := false
                    for _, otherValue := range o[i] {
                        if otherValue == cs[i][j] {
                            inOther = true
                            break
                        }
                    }
                    if !inOther {
                        found = false
                        break
                    }
                }
            }
            if !found {
                b = false
                s += fmt.Sprintf("Differ on line %d - %v vs %v", i, cs[i], o[i])
            }
        }
    default:
        b = false
        s = "Cannot compare Set against given type."
    }

    return
}


