package sudoku

import (
    "sort"
)

type Cell []int

func C(vals ...int) Cell {
    return Cell(vals)
}

// Given two sets of ints, calculate the intersection
func (prevCalcd Cell) intersection(remaining Cell) Cell {
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

func (c Cell) union(other Cell) Cell {
    for _, o := range other {
        c = append(c, o)
    }
    sort.Ints(c)
    return c
}

func (c Cell) remove(toRemove int) Cell {
    c2 := C()
    for _, v := range c {
        if v != toRemove {
            c2 = append(c2, v)
        }
    }
    return c2
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

func (cell Cell) IsSolved() bool {
    return len(cell) == 1
}
