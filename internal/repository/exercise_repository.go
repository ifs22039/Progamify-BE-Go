package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type ExerciseRepository interface {
	// FindById returns the exercise along with its questions (and answers).
	// If `theta` is non‑zero we will choose up to five items whose difficulty
	// (beta) is closest to the supplied ability using utils.SelectNextItem.
	FindById(id uint, theta float64, userID uint) (*model.Exercise, error)
	AddTakeExercise(userID uint, request model.SubmitExerciseRequest) (*model.TakeExercise, error)
}

type exerciseRepository struct {
	db       *gorm.DB
	userRepo UserRepository
}

func NewExerciseRepository(db *gorm.DB, userRepo UserRepository) ExerciseRepository {
	return &exerciseRepository{db, userRepo}
}

func (e *exerciseRepository) AddTakeExercise(
	userID uint,
	request model.SubmitExerciseRequest,
) (*model.TakeExercise, error) {

	fmt.Println("✅ AddTakeExercise DIPANGGIL | userID:", userID)

	// 1. Ambil Exercise beserta Questions dan Answers
	var exercise model.Exercise
	if err := e.db.Preload("Questions.Answers").
		First(&exercise, request.ExerciseID).Error; err != nil {
		return nil, err
	}

	// 2. Ambil Lesson
	var lesson model.Lesson
	if err := e.db.First(&lesson, exercise.LessonID).Error; err != nil {
		return nil, err
	}

	// 3. Ambil data user (untuk theta)
	user, err := e.userRepo.FindById(userID)
	if err != nil {
		return nil, err
	}

	thetaBefore := user.Theta
	betaBefore := exercise.Beta

	// 4. Hitung nomor attempt
	var attemptCount int64
	e.db.Model(&model.TakeExercise{}).
		Where("user_id = ? AND exercise_id = ?", userID, request.ExerciseID).
		Count(&attemptCount)
	attemptNumber := int(attemptCount) + 1

	// =========================
	// 5. GRADING
	// =========================

	answersJSON := request.Answers
	// if client didn't send any answers, nothing will be graded later
	if len(answersJSON) == 0 {
		log.Printf("⚠️ AddTakeExercise: empty answers payload, request=%+v", request)
	}
	takeExerciseAnswer := make(map[int]interface{})

	totalCorrect := 0.0
	totalExp := 0
	totalPoint := 0
	rewardExp := 0
	rewardPoint := 0

	// gunakan theta berganti secara progresif setiap soal sesuai permintaan (IRT setelah setiap check)
	thetaCurrent := thetaBefore

	// Map untuk tracking hasil jawaban benar/salah per soal (dipakai untuk update beta)
	questionResults := make(map[uint]bool)

	// counters used by the new IRT utilities; we keep separate integers for
	// theta updates (library UpdateTheta expects correct/total counts).
	correctCountProgress := 0
	answeredCount := 0

	// Threshold bisa diatur via config nanti (misal di env atau di tabel exercise)
	const essayGradingThreshold = 50.0

	for key, value := range answersJSON {
		// grading loop per soal
		detail, ok := value.(map[string]interface{})
		if !ok {
			log.Printf("Invalid answer detail format for key %v", key)
			continue
		}

		var question model.ExQuestion

		// Cari question berdasarkan question_id (cara utama)
		if qid, ok := detail["question_id"].(float64); ok {
			if err := e.db.Preload("Answers").First(&question, uint(qid)).Error; err != nil {
				log.Printf("Question not found: %v", err)
				continue
			}
		} else {
			// fallback (jarang dipakai)
			return nil, fmt.Errorf("missing or invalid question_id for answer key %v", key)
		}

		exp := 0
		point := 0
		rewardExp += question.Exp
		rewardPoint += question.Point

		switch question.Type {
		case "multiple_choice":
			var correct model.ExAnswer
			var correctIndex int

			for i, a := range question.Answers {
				if a.IsCorrect {
					correct = a
					correctIndex = i
					break
				}
			}

			wasCorrect := false
			if ansid, ok := detail["answer_id"].(float64); ok {
				if int(correct.ID) == int(ansid) {
					exp = question.Exp
					point = question.Point
					totalCorrect++
					wasCorrect = true
				}
			}
			questionResults[question.ID] = wasCorrect

			totalExp += exp
			totalPoint += point

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":          question.ID,
				"type":                 question.Type,
				"exp_gained":           exp,
				"point_gained":         point,
				"correct_answer_index": correctIndex,
			}

		case "matching":
			// build correct pairs from data stored in question.Content first,
			// then fall back to question.Answers if necessary.  The mobile app
			// populates the pair data in Content, so failing to read it would
			// always mark matching items incorrect (as seen in the logs).
			type pair struct{ keyword, explanation string }
			var correctPairs []pair
			// try parsing JSON array from Content field
			if strings.TrimSpace(question.Content) != "" {
				var arr []map[string]interface{}
				if err := json.Unmarshal([]byte(question.Content), &arr); err == nil {
					for _, item := range arr {
						kw, _ := item["keyword"].(string)
						exp, _ := item["explanation"].(string)
						correctPairs = append(correctPairs, pair{
							keyword:     strings.TrimSpace(kw),
							explanation: strings.TrimSpace(exp),
						})
					}
				}
			}
			// fallback to Answers table if Content had nothing useful
			if len(correctPairs) == 0 {
				for _, ans := range question.Answers {
					content := strings.TrimSpace(ans.Content)
					var parsed map[string]string
					if err := json.Unmarshal([]byte(content), &parsed); err == nil {
						correctPairs = append(correctPairs, pair{
							keyword:     strings.TrimSpace(parsed["keyword"]),
							explanation: strings.TrimSpace(parsed["explanation"]),
						})
					} else {
						correctPairs = append(correctPairs, pair{keyword: content})
					}
				}
			}
			if len(correctPairs) == 0 {
				log.Printf("⚠️ Matching question %d has no stored answer pairs", question.ID)
			}

			submittedPairs, ok := detail["answers"].([]interface{})
			if !ok || len(submittedPairs) == 0 {
				continue
			}

			correctCount := 0
			totalPairs := len(submittedPairs)

			var userMatches []map[string]interface{}

			for _, pairAny := range submittedPairs {
				pair, ok := pairAny.(map[string]interface{})
				if !ok {
					continue
				}

				submittedKeyword := strings.TrimSpace(pair["keyword"].(string))
				submittedExplanation := strings.TrimSpace(pair["explanation"].(string))

				userMatches = append(userMatches, map[string]interface{}{
					"explanation": submittedExplanation,
					"keyword":     submittedKeyword,
				})

				// look for an exact match in the correctPairs slice
				for _, cp := range correctPairs {
					if cp.keyword == submittedKeyword &&
						(cp.explanation == "" || cp.explanation == submittedExplanation) {
						correctCount++
						break
					}
				}
			}

			// compute ratio of correct pairs, use it for partial credit
			ratio := 0.0
			if totalPairs > 0 {
				ratio = float64(correctCount) / float64(totalPairs)
			}

			// award fractional totalCorrect and proportional exp/point
			totalCorrect += ratio
			exp = int(ratio * float64(question.Exp))
			point = int(ratio * float64(question.Point))

			isFullyCorrect := correctCount == totalPairs && totalPairs > 0
			questionResults[question.ID] = isFullyCorrect

			totalExp += exp
			totalPoint += point
			// expose correct keywords for debugging / response
			var correctKeywords []string
			for _, cp := range correctPairs {
				if cp.keyword != "" {
					correctKeywords = append(correctKeywords, cp.keyword)
				}
			}
			// override later when building response
			// (assignment below)

			// store local variables for later output
			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":      question.ID,
				"type":             question.Type,
				"exp_gained":       exp,
				"point_gained":     point,
				"correct_keywords": correctKeywords,
				"user_matches":     userMatches,
				"correct_count":    correctCount,
				"total_pairs":      totalPairs,
				"ratio":            ratio,
				"is_fully_correct": isFullyCorrect,
			}

case "essay":
		userAnswerText, ok := detail["answer_text"].(string)
		if !ok {
			userAnswerText = ""
		}
		userAnswerText = strings.TrimSpace(userAnswerText)

		if userAnswerText == "" {
			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":  question.ID,
				"type":         question.Type,
				"exp_gained":   0,
				"point_gained": 0,
				"user_answer":  "",
				"status":       "empty",
			}
			continue
		}

		// find a reference answer, if any
		var correctText string
		for _, ans := range question.Answers {
			if ans.IsCorrect {
				correctText = ans.Content
				break
			}
		}

		if correctText == "" {
			// no reference stored in database; skip external grading and
			// mark the response as ungraded.  This usually means the question
			// was created without an answer (perhaps an essay intended for
			// manual review).  The mobile client did send text, but we cannot
			// compare it automatically.
			log.Printf("⚠️ No correct answer found for essay question %d, skipping auto‑grading", question.ID)
		} else {
			// Grade essay using the external service and compare against threshold
			gradeResult, err := utils.EssayGrading(correctText, userAnswerText)
			if err != nil {
				log.Printf("Essay grading failed for qID %d: %v", question.ID, err)
				gradeResult = utils.EssayGradeResult{} // zero values
			}

			score := gradeResult.FinalScore * 100 // percentage used for IRT/points
			isCorrectEssay := score >= essayGradingThreshold

			if isCorrectEssay {
				exp = question.Exp
				point = question.Point
				totalCorrect += 1
			}
			questionResults[question.ID] = isCorrectEssay

			totalExp += exp
			totalPoint += point

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":      question.ID,
				"type":             question.Type,
				"exp_gained":       exp,
				"point_gained":     point,
				"user_answer":      userAnswerText,
				"correct_answer":   correctText,
				"final_score":       gradeResult.FinalScore,
				"similarity_score":  gradeResult.SimilarityScore,
				"essay_score":      score,
				"threshold":        essayGradingThreshold,
				"is_correct":       isCorrectEssay,
				"feedback":         question.Feedback,
			}
		}

		// when correctText == "" we still want to record the user's answer
		takeExerciseAnswer[key] = map[string]interface{}{
			"question_id":  question.ID,
			"type":         question.Type,
			"exp_gained":   0,
			"point_gained": 0,
			"user_answer":  userAnswerText,
			"status":       "no_reference",
		}

	case "short_answer":
		userAnswerText, ok := detail["answer_text"].(string)
		if !ok {
			userAnswerText = ""
		}
		userAnswerText = strings.TrimSpace(userAnswerText)

		if userAnswerText == "" {
			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":  question.ID,
				"type":         question.Type,
				"exp_gained":   0,
				"point_gained": 0,
				"user_answer":  "",
				"status":       "empty",
			}
			continue
		}

		// query database for correct answer(s)
		var correctAnswers []model.ExAnswer
		if err := e.db.Where("ex_question_id = ? AND is_correct = ?", question.ID, true).
			Find(&correctAnswers).Error; err != nil {
			log.Printf("Failed to fetch correct answers for short_answer question %d: %v", question.ID, err)
			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":  question.ID,
				"type":         question.Type,
				"exp_gained":   0,
				"point_gained": 0,
				"user_answer":  userAnswerText,
				"status":       "db_error",
			}
			continue
		}

		if len(correctAnswers) == 0 {
			log.Printf("No correct answer found in DB for short_answer question %d", question.ID)
			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":  question.ID,
				"type":         question.Type,
				"exp_gained":   0,
				"point_gained": 0,
				"user_answer":  userAnswerText,
				"status":       "no_correct_answer",
			}
			continue
		}

		// check if user answer matches any correct answer (exact match, case-insensitive)
		isCorrect := false
		var matchedAnswer string
		for _, correctAns := range correctAnswers {
			if strings.EqualFold(userAnswerText, strings.TrimSpace(correctAns.Content)) {
				isCorrect = true
				matchedAnswer = correctAns.Content
				break
			}
		}

		if isCorrect {
			exp = question.Exp
			point = question.Point
			totalCorrect += 1
		}
		questionResults[question.ID] = isCorrect
		totalExp += exp
		totalPoint += point

		takeExerciseAnswer[key] = map[string]interface{}{
			"question_id":    question.ID,
			"type":           question.Type,
			"exp_gained":     exp,
			"point_gained":   point,
			"user_answer":    userAnswerText,
			"correct_answer": matchedAnswer,
			"is_correct":     isCorrect,
			}

		case "multiple_answer":
		// simple set equality check; partial credit could be added later
		var userIDs []int
		if arr, ok := detail["answers"].([]interface{}); ok {
			for _, v := range arr {
				if m, ok := v.(map[string]interface{}); ok {
					if aid, ok := m["answer_id"].(float64); ok {
						userIDs = append(userIDs, int(aid))
					}
				}
			}
		}
		var correctIDs []int
		for _, a := range question.Answers {
			if a.IsCorrect {
				correctIDs = append(correctIDs, int(a.ID))
			}
		}
		sort.Ints(userIDs)
		sort.Ints(correctIDs)
		wasCorrect := reflect.DeepEqual(userIDs, correctIDs)
		if wasCorrect {
			exp = question.Exp
			point = question.Point
			totalCorrect += 1
		}
		questionResults[question.ID] = wasCorrect
		totalExp += exp
		totalPoint += point
		takeExerciseAnswer[key] = map[string]interface{}{
			"question_id":      question.ID,
			"type":             question.Type,
			"exp_gained":       exp,
			"point_gained":     point,
			"user_answers":     userIDs,
			"correct_answers":  correctIDs,
			"is_correct":        wasCorrect,
		}

	case "true_false":
		// user may send answer_text = "true"/"false"
		correct := ""
		if len(question.Answers) > 0 {
			correct = strings.ToLower(question.Answers[0].Content)
		}
		userAns, _ := detail["answer_text"].(string)
		userAns = strings.ToLower(strings.TrimSpace(userAns))
		wasCorrect := userAns == correct
		if wasCorrect {
			exp = question.Exp
			point = question.Point
			totalCorrect += 1
		}
		questionResults[question.ID] = wasCorrect
		totalExp += exp
		totalPoint += point
		takeExerciseAnswer[key] = map[string]interface{}{
			"question_id": question.ID,
			"type":        question.Type,
			"exp_gained":  exp,
			"point_gained": point,
			"user_answer":  userAns,
			"correct_answer": correct,
			"is_correct":   wasCorrect,
		}

	default:
		log.Printf("Unsupported question type '%s' for question %d", question.Type, question.ID)
	}

		if m, ok := takeExerciseAnswer[key].(map[string]interface{}); ok {
			m["question"] = question
			takeExerciseAnswer[key] = m
		}

		// === IRT per-soal (theta diperbarui berdasarkan hitungan jawaban) ===
		// maintain running totals so that theta is recalculated after every item
		// (makes the adaptation visible within a single exercise session).
		// NOTE: totalCorrect variable is float64 and used for scoring; we keep
		// separate integer counters for the theta update.
		
		// these variables are declared above the loop
		
		// add to progress counters
		answeredCount++
		if questionResults[question.ID] {
			correctCountProgress++
		}

		// probability used for logging/diagnostics only
		p := utils.RaschProbability(thetaCurrent, question.Beta)
		// compute new theta based on accumulated counts
		thetaCurrent = utils.UpdateTheta(correctCountProgress, answeredCount)
		log.Printf("🔁 IRT soal %d | beta=%.4f | correct=%v | theta -> %.4f | p=%.4f",
			question.ID, question.Beta, questionResults[question.ID], thetaCurrent, p)
	}

	// because sessions always consist of 5 questions, normalize score to 5
	denom := 5.0
	if float64(len(exercise.Questions)) < denom {
		denom = float64(len(exercise.Questions))
	}
	score := totalCorrect / denom

	// jumlah soal yang sebenarnya diberikan (biasanya 5)
	totalQuestion := len(exercise.Questions)

	// thetaAfter was updated inside the loop; use that value
	thetaAfter := thetaCurrent
	betaAfter := betaBefore // exercise-level beta tidak diubah; beta per soal diupdate di bawah

	// =========================
	// 6b. Update Beta per Soal
	// =========================
	// Setelah grading selesai, update difficulty (beta) setiap soal yang dijawab.
	// Kita menghitung statistik historis (jumlah benar/total) lalu menambahkan
	// jawaban saat ini; fungsi utils.UpdateBeta menerima hitungan tersebut.
	for qID, wasCorrect := range questionResults {
		var q model.ExQuestion
		if err := e.db.First(&q, qID).Error; err != nil {
			log.Printf("⚠️ UpdateBeta: soal %d tidak ditemukan: %v", qID, err)
			continue
		}
		correctCnt, totalCnt, _ := e.computeQuestionStats(qID)
		// include current response
		totalCnt++
		if wasCorrect {
			correctCnt++
		}
		newBeta := utils.UpdateBeta(correctCnt, totalCnt)
		if err := e.db.Model(&model.ExQuestion{}).Where("id = ?", qID).Update("beta", newBeta).Error; err != nil {
			log.Printf("⚠️ UpdateBeta: gagal update beta soal %d: %v", qID, err)
		} else {
			log.Printf("📊 Beta soal %d: %.4f → %.4f (prevStats=%d/%d) (correct=%v)",
				qID, q.Beta, newBeta, correctCnt, totalCnt, wasCorrect)
		}
	}

	// =========================
	// 7. Update user theta
	// =========================

	if err := e.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("theta", thetaAfter).Error; err != nil {
		return nil, err
	}

	// =========================
	// 8. Simpan TakeExercise
	// =========================

	answerDetail, err := json.Marshal(takeExerciseAnswer)
	if err != nil {
		return nil, err
	}

	newTake := model.TakeExercise{
		ExerciseID:    exercise.ID,
		LessonID:      lesson.ID,
		UserID:        userID,
		TopicID:       lesson.TopicID,
		AttemptNumber: attemptNumber,
		Answers:       answerDetail,

		Score:         score,
		TotalCorrect:  int(totalCorrect),
		TotalQuestion: int(totalQuestion),
		TotalExp:      totalExp,
		TotalPoint:    totalPoint,
		RewardExp:     rewardExp,
		RewardPoint:   rewardPoint,

		ThetaBefore: thetaBefore,
		ThetaAfter:  thetaAfter,
		BetaBefore:  betaBefore,
		BetaAfter:   betaAfter,
	}

	if err := e.db.Create(&newTake).Error; err != nil {
		return nil, err
	}

	// =========================
	// 9. Gamification
	// =========================

	_ = e.userRepo.AddPoint(userID, totalPoint)
	_ = e.userRepo.AddExp(userID, totalExp)
	_ = e.userRepo.CheckLevel(userID)

	fmt.Printf("✅ IRT OK | theta: %.4f → %.4f | score: %.2f%%\n", thetaBefore, thetaAfter, score*100)

	return &newTake, nil
}

// computeQuestionStats loads previous attempts and returns the number of
// correct answers / total answers for the specified question.  The routine
// inspects the Answers JSON inside TakeExercise rows and treats any
// exp_gained>0 as a correct response (this mirrors grading logic used above).
func (e *exerciseRepository) computeQuestionStats(qID uint) (correct, total int, err error) {
	var takes []model.TakeExercise
	if err = e.db.Find(&takes, "answers LIKE ?", fmt.Sprintf("%%\"question_id\":%d%%", qID)).Error; err != nil {
		return
	}
	for _, t := range takes {
		var ansMap map[string]interface{}
		if err2 := json.Unmarshal(t.Answers, &ansMap); err2 != nil {
			continue
		}
		for _, v := range ansMap {
			if detail, ok := v.(map[string]interface{}); ok {
				if idF, ok := detail["question_id"].(float64); ok && uint(idF) == qID {
					total++
					if expF, ok := detail["exp_gained"].(float64); ok && expF > 0 {
						correct++
					}
				}
			}
		}
	}
	return
}

// userCorrectQuestions returns a set of question IDs that the specified user
// has already answered correctly for the given exercise.  This is used to avoid
// presenting the same item again during adaptive selection.
func (e *exerciseRepository) userCorrectQuestions(userID, exerciseID uint) (map[int]struct{}, error) {
	res := make(map[int]struct{})
	if userID == 0 {
		return res, nil
	}

	var takes []model.TakeExercise
	if err := e.db.Where("user_id = ? AND exercise_id = ?", userID, exerciseID).Find(&takes).Error; err != nil {
		return nil, err
	}

	for _, t := range takes {
		var ansMap map[string]interface{}
		if err := json.Unmarshal(t.Answers, &ansMap); err != nil {
			continue
		}
		for _, v := range ansMap {
			if detail, ok := v.(map[string]interface{}); ok {
				if idF, ok := detail["question_id"].(float64); ok {
					if expF, ok := detail["exp_gained"].(float64); ok && expF > 0 {
						res[int(idF)] = struct{}{}
					}
				}
			}
		}
	}
	return res, nil
}

// FindById fetches the exercise along with its questions and answers.  When
// a non‑zero theta is supplied we perform an adaptive selection of up to five
// items whose beta values are closest to theta using utils.SelectNextItem.  If
// theta == 0 we fall back to a random subset (previous behaviour).
func (e *exerciseRepository) FindById(id uint, theta float64, userID uint) (*model.Exercise, error) {
	var exercise model.Exercise
	if err := e.db.
		Preload("Questions.Answers").
		First(&exercise, id).Error; err != nil {
		return nil, err
	}

	// remove questions the user has already answered correctly
	if userID != 0 {
		if answered, err := e.userCorrectQuestions(userID, id); err == nil && len(answered) > 0 {
			remaining := make([]model.ExQuestion, 0, len(exercise.Questions))
			for _, q := range exercise.Questions {
				if _, ok := answered[int(q.ID)]; ok {
					continue
				}
				remaining = append(remaining, q)
			}
			exercise.Questions = remaining
		}
	}

	// if there are five or fewer items left just return them; otherwise we'll
	// perform either adaptive or random selection depending on theta.
	if len(exercise.Questions) <= 5 {
		return &exercise, nil
	}

	avail := make(map[int]float64, len(exercise.Questions))
	for _, q := range exercise.Questions {
		avail[int(q.ID)] = q.Beta
	}

	var chosenIDs []uint
	if theta != 0 {
		for len(chosenIDs) < 5 {
			next := utils.SelectNextItem(theta, avail)
			if next == -1 {
				break
			}
			chosenIDs = append(chosenIDs, uint(next))
			delete(avail, next)
		}
	} else {
		keys := make([]int, 0, len(avail))
		for k := range avail {
			keys = append(keys, k)
		}
		rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
		for i := 0; i < 5 && i < len(keys); i++ {
			chosenIDs = append(chosenIDs, uint(keys[i]))
		}
	}

	var filtered []model.ExQuestion
	for _, q := range exercise.Questions {
		for _, id2 := range chosenIDs {
			if q.ID == id2 {
				filtered = append(filtered, q)
				break
			}
		}
	}
	exercise.Questions = filtered
	return &exercise, nil
}

