package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type BadgeRepository interface {
	FindBadge(id uint) (*model.Badge, error)
	Create(badge *model.Badge) error
	AddHaveBadge(userId uint, badgeId uint) (*model.HaveBadge, error)
	AssignBadgeIfEligible(userId uint) ([]model.HaveBadge, error)
	GetBadges(id uint) ([]model.UserBadge, error)
}

type badgeRepository struct {
	db *gorm.DB
	userRepo UserRepository
}

func NewBadgeRepository(db *gorm.DB, userRepo UserRepository) BadgeRepository {
	return &badgeRepository{db, userRepo}
}

func (r *badgeRepository) FindBadge(id uint) (*model.Badge, error) {
	var badge model.Badge
	err := r.db.First(&badge, id).Error
	if err != nil {
		return nil, err
	}
	return &badge, nil
}

func (r *badgeRepository) Create(badge *model.Badge) error {
	return r.db.Create(badge).Error
}

func (r *badgeRepository) AddHaveBadge(userId uint, badgeId uint) (*model.HaveBadge, error) {
	var user model.User

	err := r.db.First(&user, userId).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	var badge model.Badge
	err = r.db.First(&badge, badgeId).Error
	if err != nil {
		return nil, fmt.Errorf("badge not found: %w", err)
	} else {
		log.Printf("✅ Badge found: ID=%d, Title=%s", badge.ID, badge.Title)
	}

	haveBadge := &model.HaveBadge{
		UserID:  userId,
		BadgeID: badgeId,
		Badge: badge,
	}

	if err := r.db.Create(haveBadge).Error; err != nil {
		return nil, fmt.Errorf("failed to assign badge: %w", err)
	}

	return haveBadge, nil
}


func (r *badgeRepository) GetBadges(userId uint) ([]model.UserBadge, error) {
	var haveBadges []model.HaveBadge
	err := r.db.Preload("Badge").Where("user_id = ?", userId).Find(&haveBadges).Error
	if err != nil {
		return nil, err
	}

	// Mapping badge_id ke jumlahnya
	badgeCount := make(map[uint]int)
	userBadgesMap := make(map[uint]model.UserBadge)

	for _, hb := range haveBadges {
		badgeCount[hb.BadgeID]++

		// Jika belum ada dalam map, tambahkan
		if _, exists := userBadgesMap[hb.BadgeID]; !exists {
			userBadgesMap[hb.BadgeID] = model.UserBadge{
				ID:          hb.Badge.ID,
				Title:       hb.Badge.Title,
				Description: hb.Badge.Description,
				Picture:     hb.Badge.Picture,
				Count:       0, // Akan diisi di bawah
			}
		}
	}

	// Update count di setiap badge yang ditemukan
	var userBadges []model.UserBadge
	for badgeID, badge := range userBadgesMap {
		badge.Count = badgeCount[badgeID]
		userBadges = append(userBadges, badge)
	}

	return userBadges, nil
}



func (r *badgeRepository) AssignBadgeIfEligible(userId uint) ([]model.HaveBadge, error) {
	var takeQuests []model.TakeQuest
	err := r.db.Where("user_id = ?", userId).Order("created_at asc").Find(&takeQuests).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch take quests: %w", err)
	}

	if len(takeQuests) == 0 {
		return nil, nil
	}

	var awardedBadges []model.HaveBadge

	// Helper function untuk menambahkan badge jika belum dimiliki
	addBadgeIfNotOwned := func(badgeID uint) error {
		var count int64
		r.db.Model(&model.HaveBadge{}).Where("user_id = ? AND badge_id = ?", userId, badgeID).Count(&count)

		if count == 0 { // ✅ Cek apakah user sudah punya badge ini
			badge, err := r.AddHaveBadge(userId, badgeID)
			if err != nil {
				return err
			}
			awardedBadges = append(awardedBadges, *badge)
		}
		return nil
	}

	// ✅ Quest Beginner: diberikan HANYA SEKALI setelah 10 TakeQuest pertama
	if len(takeQuests) >= 10 {
		_ = addBadgeIfNotOwned(1)
	}

	// ✅ Warrior: 7 salah berturut-turut sebelum 1 benar
	if len(takeQuests) >= 8 {
		consecutiveWrong := 0
		for i := len(takeQuests) - 2; i >= 0; i-- { // Cek 7 terakhir sebelum yang benar
			if takeQuests[i].IsCorrect {
				break // Harus benar di quest ke-8
			}
			consecutiveWrong++
		}
		if consecutiveWrong >= 7 && takeQuests[len(takeQuests)-1].IsCorrect {
			_ = addBadgeIfNotOwned(2)
		}
	}

	// ✅ Perfect Streak (3, 5, 10): hanya jika streak benar baru dimulai setelah salah
	if len(takeQuests) >= 4 {
		correctStreak := 0
		for i := len(takeQuests) - 1; i >= 0; i-- {
			if takeQuests[i].IsCorrect {
				correctStreak++
			} else {
				break // Harus ada satu false sebelum streak benar
			}
		}

		if correctStreak == 3 {
			_ = addBadgeIfNotOwned(3)
		}
		if correctStreak == 5 {
			_ = addBadgeIfNotOwned(4)
		}
		if correctStreak == 10 {
			_ = addBadgeIfNotOwned(5)
		}
	}

	// ✅ Unstoppable Challenger: minimal 1 quest per hari selama 7 hari berturut-turut, lalu di-reset
	streakDays := 1
	lastAwardedDay := time.Time{} // Untuk mereset streak setelah badge diberikan
	for i := 1; i < len(takeQuests); i++ {
		diffHours := takeQuests[i].CreatedAt.Sub(takeQuests[i-1].CreatedAt).Hours()
		if diffHours <= 24 {
			streakDays++
		} else {
			streakDays = 1 // Reset jika ada jeda lebih dari 1 hari
		}

		if streakDays >= 7 {
			if lastAwardedDay.IsZero() || takeQuests[i].CreatedAt.Sub(lastAwardedDay).Hours() > 24 {
				_ = addBadgeIfNotOwned(6)
				lastAwardedDay = takeQuests[i].CreatedAt // Reset streak setelah badge diberikan
			}
			streakDays = 0 // Mulai hitung ulang
		}
	}

	// Jika tidak ada badge yang diberikan, return nil
	if len(awardedBadges) == 0 {
		return nil, nil
	}

	return awardedBadges, nil
}
