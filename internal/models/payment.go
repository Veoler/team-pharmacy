package models

import (
	"time"

	"gorm.io/gorm"
)


type PayStatus string

const (
	PayStatusPending	PayStatus = "pending"
	PayStatusSuccess	PayStatus = "success"
	PayStatusFailed		PayStatus = "failed"
)

type PayMethod string

const (
	PayMethodCard			PayMethod = "card"
	PayMethodCash			PayMethod = "cash"
	PayMethodOnlineWallet	PayMethod = "online_wallet"
)

type Payment struct {
	gorm.Model
	Order	*Order		`json:"-"`	
	OrderID	uint		`json:"order_id" gorm:"not null;index"`
	Amount	int			`json:"amount" gorm:"not null"`					// сумма платежа:
	Status	PayStatus	`json:"status" gorm:"type:varchar(16);not null"`// pending, success, failed;
	Method	PayMethod 	`json:"method" gorm:"type:varchar(16);not null"`// card, cash, online_wallet и т.п.;
	PaidAt	*time.Time	`json:"paid_at"`								// время успешного платежа.
}

type PaymentCreateRequest struct {
	OrderID		uint		`json:"order_id"`
	Amount		int			`json:"amount"`
	Method		PayMethod		`json:"method"`
	Status		PayStatus		`json:"status"`
	PaidAt		*time.Time	`json:"paid_at,omitempty"`
}
