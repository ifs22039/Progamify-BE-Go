package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"fmt"
	"log"
	"sort"
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

	sort.Slice(userBadges, func(i, j int) bool {
		return userBadges[i].ID < userBadges[j].ID
	})

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
		var countQuestBeginner int64
		r.db.Model(&model.HaveBadge{}).Where("user_id = ? AND badge_id = ?", userId, 1).Count(&countQuestBeginner)

		if countQuestBeginner == 0 { // ✅ Cek apakah user sudah punya badge ini
			log.Printf("User %d belum memiliki badge %d, menambahkan badge...", userId, badgeID)
			badge, err := r.AddHaveBadge(userId, badgeID)
			if err != nil {
				return err
			}
			awardedBadges = append(awardedBadges, *badge)
		}
		if badgeID != 1 {
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
	if len(takeQuests) >= 3 {
    // Ambil HaveBadge terakhir berdasarkan ID 2,3,4
    var lastHaveBadge2, lastHaveBadge3, lastHaveBadge4 model.HaveBadge

    r.db.Where("user_id = ? AND badge_id = ?", userId, 3).
        Order("created_at DESC").First(&lastHaveBadge2)

    r.db.Where("user_id = ? AND badge_id = ?", userId, 4).
        Order("created_at DESC").First(&lastHaveBadge3)

    r.db.Where("user_id = ? AND badge_id = ?", userId, 5).
        Order("created_at DESC").First(&lastHaveBadge4)

    // Cek takeQuest ke-3 terakhir
    if len(takeQuests) >= 3 {
        lastTakeQuest3 := takeQuests[len(takeQuests)-3]
        if lastTakeQuest3.IsCorrect && takeQuests[len(takeQuests)-2].IsCorrect && takeQuests[len(takeQuests)-1].IsCorrect {
            if lastHaveBadge2.ID == 0 || lastHaveBadge2.CreatedAt.Before(lastTakeQuest3.CreatedAt) {
                _ = addBadgeIfNotOwned(3)
            }
        }
    }

    // Cek takeQuest ke-5 terakhir
    if len(takeQuests) >= 5 {
        lastTakeQuest5 := takeQuests[len(takeQuests)-5]
        if lastTakeQuest5.IsCorrect && takeQuests[len(takeQuests)-4].IsCorrect &&
            takeQuests[len(takeQuests)-3].IsCorrect && takeQuests[len(takeQuests)-2].IsCorrect && takeQuests[len(takeQuests)-1].IsCorrect {
            if lastHaveBadge3.ID == 0 || lastHaveBadge3.CreatedAt.Before(lastTakeQuest5.CreatedAt) {
                _ = addBadgeIfNotOwned(4)
            }
        }
    }

    // Cek takeQuest ke-10 terakhir
    if len(takeQuests) >= 10 {
        lastTakeQuest10 := takeQuests[len(takeQuests)-10]
        if lastTakeQuest10.IsCorrect && takeQuests[len(takeQuests)-9].IsCorrect &&
            takeQuests[len(takeQuests)-8].IsCorrect && takeQuests[len(takeQuests)-7].IsCorrect &&
            takeQuests[len(takeQuests)-6].IsCorrect && takeQuests[len(takeQuests)-5].IsCorrect &&
            takeQuests[len(takeQuests)-4].IsCorrect && takeQuests[len(takeQuests)-3].IsCorrect &&
            takeQuests[len(takeQuests)-2].IsCorrect && takeQuests[len(takeQuests)-1].IsCorrect {
            if lastHaveBadge4.ID == 0 || lastHaveBadge4.CreatedAt.Before(lastTakeQuest10.CreatedAt) {
                _ = addBadgeIfNotOwned(5)
            }
        }
    }
	}


	// ✅ Unstoppable Challenger: 1 quest per hari selama 7 hari berturut-turut
	streak := 0
	var lastDate time.Time

	var lastHaveBadge model.HaveBadge
	r.db.Where("user_id = ? AND badge_id = ?", userId, 6).
		Order("created_at DESC").
		First(&lastHaveBadge)

	for _, tq := range takeQuests {
		date := tq.CreatedAt.Truncate(24 * time.Hour)
		if lastDate.IsZero() || date.Sub(lastDate) == 24*time.Hour {
			streak++
		} else if date.Sub(lastDate) > 24*time.Hour {
			streak = 1 // Reset streak jika ada hari yang terlewat
		}
		lastDate = date

		if streak == 7 {
			if lastHaveBadge.ID == 0 || lastHaveBadge.CreatedAt.Before(tq.CreatedAt) {
				if err := addBadgeIfNotOwned(6); err != nil {
					log.Printf("[ERROR] Failed to assign badge: %v", err)
					return nil, err
				}
				log.Printf("[INFO] Badge assigned successfully to user_id: %d", userId)
			}
			streak = 0
		}
	}


	// Jika tidak ada badge yang diberikan, return nil
	if len(awardedBadges) == 0 {
		return nil, nil
	}

	return awardedBadges, nil
}
