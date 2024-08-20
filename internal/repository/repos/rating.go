package repos

import (
	"context"
	"github.com/timurzdev/social-rating-bot/internal/infrastructure/postgres"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/timurzdev/social-rating-bot/internal/repository/model"
)

const (
	ratingTable = "ratings"
)

type RatingRepository struct {
	db *postgres.DB
}

func NewRatingRepository(db *postgres.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) insert(ctx context.Context, tx *sqlx.Tx, rating *model.Rating) error {
	clause := map[string]interface{}{
		"user_id": rating.UserID,
		"chat_id": rating.ChatID,
		"rating":  rating.Rating,
		"level":   rating.Level,
	}
	query, args, _ := sq.Insert(ratingTable).
		SetMap(clause).
		ToSql()

	_, err := tx.ExecContext(ctx, query, args...)
	return err

}

func (r *RatingRepository) Create(ctx context.Context, rating *model.Rating) error {
	return postgres.Transaction(ctx,
		func(tx *sqlx.Tx) error {
			return r.insert(ctx, tx, rating)
		},
		r.db,
	)
}
