package scoring

func Upper90Club(result Result, prediction Prediction) int {
	if prediction == Prediction(result) {
		return 1
	}
	return 0
}
