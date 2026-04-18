package scoring

type Result struct{ Home, Away int }
type Prediction struct{ Home, Away int }

func AcesRadio(result Result, prediction Prediction) int {
	if prediction == Prediction(result) {
		return 15
	}
	if prediction.Home == result.Away && prediction.Away == result.Home {
		return -15
	}
	if outcome(result.Home, result.Away) == outcome(prediction.Home, prediction.Away) {
		return 10
	}
	return 0
}

func outcome(home, away int) int {
	switch {
	case home > away:
		return 1
	case away > home:
		return -1
	default:
		return 0
	}
}
