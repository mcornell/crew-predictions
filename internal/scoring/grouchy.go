package scoring

func Grouchy(result Result, prediction Prediction, targetIsHome bool) int {
	if grouchyCategory(crewMargin(result, targetIsHome)) == grouchyCategory(crewMargin(predictionToResult(prediction), targetIsHome)) {
		return 1
	}
	return 0
}

func crewMargin(r Result, targetIsHome bool) int {
	if targetIsHome {
		return r.Home - r.Away
	}
	return r.Away - r.Home
}

func predictionToResult(p Prediction) Result {
	return Result{Home: p.Home, Away: p.Away}
}

func grouchyCategory(margin int) int {
	switch {
	case margin >= 2:
		return 2
	case margin == 1:
		return 1
	case margin == 0:
		return 0
	case margin == -1:
		return -1
	default:
		return -2
	}
}
