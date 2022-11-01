package carrier

func (j JetIDs) removeJetID(jetID string) JetIDs {
	if index := j.indexOf(jetID); index != -1 {
		return j.removeAt(index)
	}
	return j
}

func (j JetIDs) indexOf(jetID string) int {
	for index, id := range j {
		if id == jetID {
			return index
		}
	}
	return -1
}

func (j JetIDs) removeAt(i int) JetIDs {
	j[i] = j[len(j)-1]
	return j[:len(j)-1]
}

func (j JetIDs) append(jetID string) (JetIDs, error) {
	if index := j.indexOf(jetID); index != -1 {
		return j, ErrJetIDExisted
	}
	return append(j, jetID), nil
}
