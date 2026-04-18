package scoring

func Upper90Club(result Result, prediction Prediction) int {
	if prediction == Prediction(result) {
		return 2
	}
	if outcome(result.Home, result.Away) == outcome(prediction.Home, prediction.Away) {
		return 1
	}
	return 0
}
