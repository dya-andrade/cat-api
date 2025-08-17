package domain

import "time"

/*
validate:"...": regras de validação do campo.
required: obrigatório.
min, max: tamanho mínimo/máximo.
gte, lte: valor mínimo/máximo.
omitempty: omite o campo no JSON se estiver vazio/nulo.
*tipo: ponteiro, permite valor nulo.
*/

type Cat struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=64"`
	AgeYears  int       `json:"age_years" validate:"gte=0,lte=40"`
	Breed     *string   `json:"breed,omitempty"`
	CoatColor *string   `json:"coat_color,omitempty"`
	WeightKG  *float64  `json:"weight_kg,omitempty" validate:"omitempty,gte=0,lte=50"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Para criação/atualização parciais

type CatCreate struct {
	Name      string   `json:"name" validate:"required,min=2,max=64"`
	AgeYears  int      `json:"age_years" validate:"gte=0,lte=40"`
	Breed     *string  `json:"breed"`
	CoatColor *string  `json:"coat_color"`
	WeightKG  *float64 `json:"weight_kg" validate:"omitempty,gte=0,lte=50"`
}

type CatUpdate struct {
	Name      *string  `json:"name" validate:"required,min=2,max=64"`
	AgeYears  *int     `json:"age_years" validate:"gte=0,lte=40"`
	Breed     *string  `json:"breed"`
	CoatColor *string  `json:"coat_color"`
	WeightKG  *float64 `json:"weight_kg" validate:"omitempty,gte=0,lte=50"`
}
