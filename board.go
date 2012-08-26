package sudoku

import (
    "fmt"
    "reflect"
)

type Board []Set

// Satisfy the Equalable interface so we can use matchers in the test.
func (b Board) Equals(other interface{}) (bool, string) {
    switch o := other.(type) {
        case Board:
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
                    if b[i][j].isEmpty() {
                        equals = equals && sameLen
                        msg += "| X "
                    } else if o[i][j].isEmpty() {
                        equals = equals && sameLen
                        msg += "| # "
                    } else if !b[i][j].Equals(o[i][j]) {
                        equals = false
                        msg += fmt.Sprintf("|%v %v", b[i][j], o[i][j])
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
