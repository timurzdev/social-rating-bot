package handlers

import (
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	permissions permissionsProvider
	rating      ratingProvider
	events      eventsProvider
}

func New(
	p permissionsProvider,
	r ratingProvider,
	e eventsProvider,
) *Handler {
	return &Handler{
		permissions: p,
		rating:      r,
		events:      e,
	}
}

func (h *Handler) ReactionMessage(c tele.Context) error {
	return nil
}
