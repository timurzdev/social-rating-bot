package repos

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/timurzdev/social-rating-bot/internal/infrastructure/sqlite"
	"github.com/timurzdev/social-rating-bot/internal/repository/model"
)

type RatingRepository struct {
	db *sqlite.DB
}

func NewRatingRepository(db *sqlite.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) insert(ctx context.Context, tx *sqlx.Tx, rating *model.Rating) error {
	clause := map[string]interface{}{
		"user_id": rating.UserID,
		"chat_id": rating.ChatID,
		"rating":  rating.Rating,
		"level":   rating.Level,
	}
	query, args, _ := sq.Insert("ratings").
		SetMap(clause).
		ToSql()

	_, err := tx.ExecContext(ctx, query, args...)
	return err

}

func (r *RatingRepository) Create(ctx context.Context, rating *model.Rating) error {
	return sqlite.Transaction(context.Background(), func(tx *sqlx.Tx) error {
		return r.insert(ctx, tx, rating)
	}, r.db)
}
