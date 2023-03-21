package math

type Vector struct {
	AxisX float64
	AxisY float64
	AxisZ float64
}

// Add return new Vector with Added value
func (v Vector) Add(vector Vector) Vector {
	return Vector{
		AxisX: v.AxisX + vector.AxisX,
		AxisY: v.AxisY + vector.AxisY,
		AxisZ: v.AxisZ + vector.AxisZ,
	}
}

// PartOf return new part of a Vector
func (v Vector) PartOf(portion float64) Vector {
	total := v.AxisX + v.AxisY + v.AxisZ
	portionFromTotal := portion / total
	if portionFromTotal > 1 {
		portionFromTotal = 1
	}

	return Vector{
		AxisX: v.AxisX * portionFromTotal,
		AxisY: v.AxisY * portionFromTotal,
		AxisZ: v.AxisZ * portionFromTotal,
	}
}
