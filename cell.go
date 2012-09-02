package sudoku

import (
    "fmt"
)

type Cell []int

// Given two sets of ints, calculate the union
func (prevCalcd Cell) union(remaining Cell) Cell {
    if len(prevCalcd) > 0 {
        rtn := make(Cell, 0)
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

func (c Cell) isEmpty() bool {
    return len(c) == 0
}

func (c Cell) Equals(o Cell) bool {
    if len(c) != len(o) {
        return false
    }
    for _, v := range c {
        found := false
        for _, v2 := range o {
            if v2 == v {
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }
    return true
}

func (cell Cell) GetWidth() int {
    return len(fmt.Sprintf("%v", cell))
}

func (cell Cell) IsSolved() bool {
    return len(cell) == 1
}
