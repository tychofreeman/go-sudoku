package sudoku

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
