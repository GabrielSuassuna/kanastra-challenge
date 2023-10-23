package entity

import (
	"time"
)

type Billet struct {
	ID           string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string    `bson:"name" json:"name"`
	GovernmentId string    `bson:"government_id" json:"government_id"`
	Email        string    `bson:"email" json:"email"`
	DebtAmount   float64   `bson:"debt_amount" json:"debt_amount"`
	DebtDueDate  string    `bson:"debt_due_date" json:"debt_due_date"`
	DebtID       string    `bson:"debt_id" json:"debt_id"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

func NewBillet(name, governmentId, email, debtDueDate, DebtId string, debtAmount float64) *Billet {
	createdAt := time.Now()
	updatedAt := time.Now()

	return &Billet{
		Name:         name,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		GovernmentId: governmentId,
		Email:        email,
		DebtAmount:   debtAmount,
		DebtDueDate:  debtDueDate,
		DebtID:       DebtId,
	}
}
