package dto

type CreateAccountingEntity struct {
 Code string
 Name string
}

type CreateAccountingBook struct {
 EntityID string
 Code string
 Name string
 BaseCurrency string
}
