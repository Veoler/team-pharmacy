package services

import (
	"errors"
	"strings"
	"time"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrPromocodeNotFound = errors.New("промокод не найден")
	ErrPromocodeExpired  = errors.New("срок действия промокода истек или еще не начался")
	ErrPromocodeInactive = errors.New("промокод неактивен")
	ErrPromocodeLimitExceeded = errors.New("превышен лимит использований промокода")
)

type PromocodeService interface {
	CreatePromocode(req *models.PromocodeCreateRequest) (*models.Promocode, error)
	GetAllPromocodes() ([]models.Promocode, error)
	UpdatePromocode(id uint, req models.PromocodeUpdateRequest) (*models.Promocode, error)
	DeletePromocode(id uint) error
	ValidatePromocode(req models.PromocodeValidateRequest) (*models.PromocodeValidateResponse, error)
	ApplyPromocode(code string, totalPrice int) (discountTotal int, finalPrice int, err error)
}


type promocodeService struct {
	promocode repository.PromocodeRepository
}

func NewPromocodeService(promocode repository.PromocodeRepository) PromocodeService {
	return &promocodeService{promocode: promocode}
}

func (s *promocodeService) CreatePromocode(req *models.PromocodeCreateRequest) (*models.Promocode, error) {
	if err := s.validatePromocodeCreate(req); err != nil {
		return nil, err
	}

	promocode := &models.Promocode{
		Code:			req.Code,
    	Description:	req.Description,
    	DiscountType:	req.DiscountType,
    	DiscountValue:	req.DiscountValue,
    	ValidFrom:		req.ValidFrom,
    	ValidTo:		req.ValidTo,
    	MaxUses:		req.MaxUses,
    	MaxUsesPerUser:	req.MaxUsesPerUser,
    	IsActive:		req.IsActive,
	}

	if err := s.promocode.Create(promocode); err != nil {
		return nil, err
	}

	return promocode, nil
}

func (s *promocodeService) GetAllPromocodes() ([]models.Promocode, error) {

	return s.promocode.GetAll()
}

func (s *promocodeService) UpdatePromocode(id uint, req models.PromocodeUpdateRequest) (*models.Promocode, error) {
	promocode, err := s.promocode.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPromocodeNotFound
		}
		return nil, err
	}

	if err := s.applyPromocodeUpdate(promocode, req); err != nil {
		return nil, err
	}

	if err := s.promocode.Update(promocode); err != nil {
		return nil, err
	}

	return promocode, nil
}

func (s *promocodeService) DeletePromocode(id uint) error {
	if _, err := s.promocode.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPromocodeNotFound
		}
		return err
	}

	if err := s.promocode.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *promocodeService) ValidatePromocode(req models.PromocodeValidateRequest) (*models.PromocodeValidateResponse, error) {
	promo, err := s.promocode.GetByCode(req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPromocodeNotFound
		}
		return nil, err
	}

	if !promo.IsActive {
		return nil, ErrPromocodeInactive
	}
 
	if err := checkPromocodeValidity(promo); err != nil {
		return nil, err
	}
 
	discountAmount := calcDiscount(promo, req.OrderAmount)

	// Скидка не может быть больше суммы заказа
	if discountAmount > req.OrderAmount {
		discountAmount = req.OrderAmount
	}

	return &models.PromocodeValidateResponse{
		Code:         promo.Code,
		Discount:     discountAmount,
		FinalAmount:  req.OrderAmount - discountAmount,
		IsApplicable: true,
	}, nil
}


func (s *promocodeService) validatePromocodeCreate(req *models.PromocodeCreateRequest) error {
	if req.Code == "" {
		return errors.New("поле code не должно быть пустым")
	}
	
	if !isValidDiscountType(req.DiscountType) {
		return errors.New("поле discount_type должно иметь значение fixed или percent")
	}
	
	if req.DiscountValue <= 0 {
		return errors.New("поле discount_value должно быть больше 0")
	}

	if req.ValidFrom.IsZero() {
		return errors.New("поле valid_from не должно быть пустым")
	}

	if req.ValidTo.IsZero() {
		return errors.New("поле valid_to не должно быть пустым")
	}

	if req.ValidTo.Before(req.ValidFrom) {
        return errors.New("дата окончания (valid_to) не может быть раньше даты начала (valid_from)")
    }
	
	return nil
}

func (s *promocodeService) applyPromocodeUpdate(promocode *models.Promocode, req models.PromocodeUpdateRequest) error {
	if req.Code != nil {
		trimmed := strings.TrimSpace(*req.Code)
		if trimmed == "" {
			return errors.New("поле code не должно быть пустым")
		}
		promocode.Code = trimmed
	}

	if req.Description != nil {
		promocode.Description = *req.Description
	}

	if req.DiscountType != nil {
		if !isValidDiscountType(*req.DiscountType) {
			return errors.New("поле discount_type должно быть fixed или percent")
		}
		promocode.DiscountType = *req.DiscountType
	}

	if req.DiscountValue != nil {
		if *req.DiscountValue <= 0 {
			return errors.New("поле discount_value должно быть больше нуля")
		}
		promocode.DiscountValue = *req.DiscountValue
	}

	if req.ValidFrom != nil {
		if req.ValidFrom.IsZero() {
			return errors.New("поле valid_from не должно быть пустым")
		}
		promocode.ValidFrom = *req.ValidFrom
	}

	if req.ValidTo != nil {
		if req.ValidTo.IsZero() {
			return errors.New("поле valide_to не должно быть пустым")
		} 
		if req.ValidTo.Before(*req.ValidFrom) {
			return errors.New("дата окончания(valide_to) не может быть раньше даты начала(valide_from)")
		}
		promocode.ValidTo = *req.ValidTo
	}

	if req.MaxUses != nil {
		promocode.MaxUses = req.MaxUses
	}

	if req.MaxUsesPerUser != nil {
		promocode.MaxUsesPerUser = req.MaxUsesPerUser
	}

	if req.IsActive != nil {
		promocode.IsActive = *req.IsActive
	}

	return nil
}

func isValidDiscountType(disType models.DisType) bool {
	switch disType {
	case models.DisTypeFixed, models.DisTypePercent:
		return true
	default:
		return false
	}
}

func (s *promocodeService) ApplyPromocode(code string, totalPrice int) (int, int, error) {
	promo, err := s.promocode.GetByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, ErrPromocodeNotFound
		}
		return 0, 0, err
	}
 
	if !promo.IsActive {
		return 0, 0, ErrPromocodeInactive
	}
 
	if err := checkPromocodeValidity(promo); err != nil {
		return 0, 0, err
	}
 
	discountTotal := int(calcDiscount(promo, float64(totalPrice)))
 
	// Скидка не может сделать итоговую сумму отрицательной.
	if discountTotal > totalPrice {
		discountTotal = totalPrice
	}
 
	finalPrice := totalPrice - discountTotal
	return discountTotal, finalPrice, nil
}

// проверяетпериод действия промокода.
func checkPromocodeValidity(promo *models.Promocode) error {
	now := time.Now()
	if now.Before(promo.ValidFrom) || now.After(promo.ValidTo) {
		return ErrPromocodeExpired
	}
	return nil
}
 
// сумма скидки в зависимости от типа промокода.
func calcDiscount(promo *models.Promocode, orderAmount float64) float64 {
	switch promo.DiscountType {
	case models.DisTypePercent:
		return (orderAmount * promo.DiscountValue) / 100
	default:
		return promo.DiscountValue
	}
}