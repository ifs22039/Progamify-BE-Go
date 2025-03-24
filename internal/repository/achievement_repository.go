package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"fmt"
	"log"
	"sort"

	// "time"

	// "golang.org/x/text/date"
	"gorm.io/gorm"
)
type TotalExercisesLessonsCompleted struct {
	TotalLessons       int64  `json:"total_lessons"`
	TotalTakeLessons   int64  `json:"total_take_lessons"`
	TotalExercises     int64  `json:"total_exercises"`
	TotalTakeExercises int64  `json:"total_take_exercises"`
	IsFinished				 bool		`json:"is_finished"`
}

type TotalTakesPerDay struct {
	TakeDate      			string `json:"take_date"`
	TotalTakeLessons   	int64  `json:"total_take_lessons"`
	TotalTakeExercises 	int64  `json:"total_take_exercises"`
	TotalTakeQuests 		int64  `json:"total_take_quests"`
}

type AchievementRepository interface {
	FindAchievement(id uint) (*model.Achievement, error)
	Create(achievement *model.Achievement) error
	AddHaveAchievement(userId uint, achievementId uint) (*model.HaveAchievement, error)
	AssignAchievementIfEligible(userId uint) ([]model.HaveAchievement, error)
	GetAchievements(id uint) ([]model.UserAchievement, error)
}

type achievementRepository struct {
	db *gorm.DB
	userRepo UserRepository
}

func NewAchievementRepository(db *gorm.DB, userRepo UserRepository) AchievementRepository {
	return &achievementRepository{db, userRepo}
}

func (r *achievementRepository) FindAchievement(id uint) (*model.Achievement, error) {
	var achievement model.Achievement
	err := r.db.First(&achievement, id).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

func (r *achievementRepository) Create(achievement *model.Achievement) error {
	return r.db.Create(achievement).Error
}

func (r *achievementRepository) AddHaveAchievement(userId uint, achievementId uint) (*model.HaveAchievement, error) {
	var user model.User

	err := r.db.First(&user, userId).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	var achievement model.Achievement
	err = r.db.First(&achievement, achievementId).Error
	if err != nil {
		return nil, fmt.Errorf("achievement not found: %w", err)
	} else {
		log.Printf("✅ Achievement found: ID=%d, Title=%s", achievement.ID, achievement.Title)
	}

	haveAchievement := &model.HaveAchievement{
		UserID:  userId,
		AchievementID: achievementId,
		Achievement: achievement,
	}

	if err := r.db.Create(haveAchievement).Error; err != nil {
		return nil, fmt.Errorf("failed to assign achievement: %w", err)
	}

	return haveAchievement, nil
}


func (r *achievementRepository) GetAchievements(userId uint) ([]model.UserAchievement, error) {
	var haveAchievements []model.HaveAchievement
	err := r.db.Preload("Achievement").Where("user_id = ?", userId).Find(&haveAchievements).Error
	if err != nil {
		return nil, err
	}

	// Mapping achievement_id ke jumlahnya
	achievementCount := make(map[uint]int)
	userAchievementsMap := make(map[uint]model.UserAchievement)

	for _, hb := range haveAchievements {
		achievementCount[hb.AchievementID]++

		// Jika belum ada dalam map, tambahkan
		if _, exists := userAchievementsMap[hb.AchievementID]; !exists {
			userAchievementsMap[hb.AchievementID] = model.UserAchievement{
				ID:          hb.Achievement.ID,
				Title:       hb.Achievement.Title,
				Description: hb.Achievement.Description,
				Picture:     hb.Achievement.Picture,
				Count:       0, // Akan diisi di bawah
			}
		}
	}

	// Update count di setiap achievement yang ditemukan
	var userAchievements []model.UserAchievement
	for achievementID, achievement := range userAchievementsMap {
		achievement.Count = achievementCount[achievementID]
		userAchievements = append(userAchievements, achievement)
	}

	sort.Slice(userAchievements, func(i, j int) bool {
		return userAchievements[i].ID < userAchievements[j].ID
	})

	return userAchievements, nil
}



func (r *achievementRepository) AssignAchievementIfEligible(userId uint) ([]model.HaveAchievement, error) {

	var haveBadges []model.HaveBadge
	err := r.db.
    Table("have_badges").
    Select("*").
    Where("user_id = ?", userId).
    Group("badge_id").
    Order("MIN(created_at) ASC").
    Find(&haveBadges).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch take badges: %w", err)
	}

	log.Printf("Successfully fetched %d badges for user_id: %d", len(haveBadges), userId)

	// Menampilkan detail setiap badge yang ditemukan
	for i, badge := range haveBadges {
		log.Printf("Badge %d: ID=%d, BadgeId=%d, UserId=%d, CreatedAt=%v", i+1, badge.ID, badge.BadgeID, badge.UserID, badge.CreatedAt)
	}
	var awardedAchievements []model.HaveAchievement

	// Helper function untuk menambahkan achievement jika belum dimiliki
	addAchievementIfNotOwned := func(achievementID uint) error {
		var countAchievement int64
		r.db.Model(&model.HaveAchievement{}).Where("user_id = ? AND achievement_id = ?", userId, achievementID).Count(&countAchievement)

		if countAchievement == 0 { // ✅ Cek apakah user sudah punya achievement ini
			log.Printf("User %d belum memiliki achievement %d, menambahkan achievement...", userId, achievementID)
			achievement, err := r.AddHaveAchievement(userId, achievementID)
			if err != nil {
				return err
			}
			awardedAchievements = append(awardedAchievements, *achievement)
		}
		return nil
	}

	// ✅ 1. Badge Collector: diberikan HANYA SEKALI setelah 3 Badges pertama
	if len(haveBadges) >= 3 {
		_ = addAchievementIfNotOwned(1)
	}


	// ✅ 2. Ambitious Learner : diberikan jika sudah belajar selama 30 hari
	var totalTakesPerDay []TotalTakesPerDay;

	queryCase2 := `
    SELECT 
    take_date,
    COUNT(DISTINCT take_lessons.id) AS total_lessons,
    COUNT(DISTINCT take_exercises.id) AS total_exercises,
    COUNT(DISTINCT take_quests.id) AS total_quests
		FROM (
				-- Ambil tanggal unik dari setiap tabel
				SELECT DISTINCT DATE(created_at) AS take_date FROM take_lessons WHERE user_id = ?
				UNION
				SELECT DISTINCT DATE(created_at) FROM take_exercises WHERE user_id = ?
				UNION
				SELECT DISTINCT DATE(created_at) FROM take_quests WHERE user_id = ?
		) AS unique_dates
		LEFT JOIN take_lessons ON DATE(take_lessons.created_at) = unique_dates.take_date AND take_lessons.user_id = 1
		LEFT JOIN take_exercises ON DATE(take_exercises.created_at) = unique_dates.take_date AND take_exercises.user_id = 1
		LEFT JOIN take_quests ON DATE(take_quests.created_at) = unique_dates.take_date AND take_quests.user_id = 1
		GROUP BY take_date
		ORDER BY take_date ASC;
    `

	err = r.db.Raw(queryCase2, userId, userId, userId).Scan(&totalTakesPerDay).Error

	if err != nil {
		return nil, err
	}

	if len(totalTakesPerDay) >= 30 {
    _ = addAchievementIfNotOwned(4)
	}


	// ✅ 3. Ultimate Badge Collector: diberikan jika sudah mendapat ke-6 badges
	if len(haveBadges) >= 6 {
		_ = addAchievementIfNotOwned(3)
	}

	// ✅ 4. Final Boss: diberikan jika sudah selesai semua section dalam courses
	var totalExercisesLessonsCompleted TotalExercisesLessonsCompleted;

	queryCase4 := `
    SELECT 
        topics.id, 
        topics.name, 
        COUNT(DISTINCT lessons.id) AS total_lessons,
        COUNT(DISTINCT take_lessons.id) AS total_take_lessons,
        COUNT(DISTINCT exercises.id) AS total_exercises,
        COUNT(DISTINCT take_exercises.id) AS total_take_exercises
    FROM topics
		LEFT JOIN lessons ON topics.id = lessons.topic_id
		LEFT JOIN take_lessons ON topics.id = take_lessons.topic_id AND take_lessons.user_id = ?
		LEFT JOIN exercises ON topics.id = exercises.topic_id
		LEFT JOIN take_exercises ON topics.id = take_exercises.topic_id AND take_exercises.user_id = ?
    GROUP BY topics.id, topics.name
    `

	err = r.db.Raw(queryCase4, userId, userId).Scan(&totalExercisesLessonsCompleted).Error

	if err != nil {
		return nil, err
	}

	if(totalExercisesLessonsCompleted.IsFinished){
		_ = addAchievementIfNotOwned(4)
	}

	// Jika tidak ada achievement yang diberikan, return nil
	if len(awardedAchievements) == 0 {
		return nil, nil
	}

	return awardedAchievements, nil
}
