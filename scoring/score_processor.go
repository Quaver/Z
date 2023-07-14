package scoring

import (
	"example.com/Quaver/Z/common"
	"math"
)

type ScoreProcessor struct {
	DifficultyRating  float64
	Modifiers         int64
	Accuracy          float64
	PerformanceRating float64
	Combo             int
	MaxCombo          int
	Judgements        map[common.Judgements]int
}

func NewScoreProcessor(difficultyRating float64, modifiers int64) *ScoreProcessor {
	return &ScoreProcessor{
		DifficultyRating: difficultyRating,
		Modifiers:        modifiers,
		Judgements:       map[common.Judgements]int{},
	}
}

// AddJudgements Adds new judgements to the score
func (sp *ScoreProcessor) AddJudgements(judgements []common.Judgements) {
	for _, j := range judgements {
		sp.addJudgement(j)
	}

	sp.calculateAccuracy()
	sp.calculatePerformanceRating()
}

// addJudgement Adds a singular judgement to the score
func (sp *ScoreProcessor) addJudgement(judgement common.Judgements) {
	sp.Judgements[judgement]++

	if judgement != common.JudgementMiss {
		sp.Combo++

		if sp.Combo > sp.MaxCombo {
			sp.MaxCombo = sp.Combo
		}
	} else {
		sp.Combo = 0
	}
}

// Calculates the accuracy of the current score
func (sp *ScoreProcessor) calculateAccuracy() {
	var acc float64 = 0

	// Since its being used in multiple places, Just to keep it shorter
	marvWeight := getJudgementAccuracyWeight(common.JudgementMarv)

	acc += float64(sp.Judgements[common.JudgementMarv]) * marvWeight
	acc += float64(sp.Judgements[common.JudgementPerf]) * getJudgementAccuracyWeight(common.JudgementPerf)
	acc += float64(sp.Judgements[common.JudgementGreat]) * getJudgementAccuracyWeight(common.JudgementGreat)
	acc += float64(sp.Judgements[common.JudgementGood]) * getJudgementAccuracyWeight(common.JudgementGood)
	acc += float64(sp.Judgements[common.JudgementOkay]) * getJudgementAccuracyWeight(common.JudgementOkay)
	acc += float64(sp.Judgements[common.JudgementMiss]) * getJudgementAccuracyWeight(common.JudgementMiss)

	totalCount := float64(sp.Judgements[common.JudgementMarv] + sp.Judgements[common.JudgementPerf] + sp.Judgements[common.JudgementGreat] +
		sp.Judgements[common.JudgementGood] + sp.Judgements[common.JudgementOkay] + sp.Judgements[common.JudgementMiss])

	sp.Accuracy = math.Max(acc/(totalCount*marvWeight), 0) * marvWeight
}

// Calculates the performance rating of the current score
func (sp *ScoreProcessor) calculatePerformanceRating() {
	sp.PerformanceRating = sp.DifficultyRating * math.Pow(sp.Accuracy/98, 6)
}

// GetJudgementAccuracyWeight Returns the accuracy weighting for a given judgement
func getJudgementAccuracyWeight(j common.Judgements) float64 {
	switch j {
	case common.JudgementMarv:
		return 100
	case common.JudgementPerf:
		return 98.25
	case common.JudgementGreat:
		return 65
	case common.JudgementGood:
		return 25
	case common.JudgementOkay:
		return -100
	case common.JudgementMiss:
		return -50
	default:
		return 0
	}
}
