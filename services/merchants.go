package services

import (
	"github.com/Guesstrain/EthBankok/config"
	"github.com/Guesstrain/EthBankok/models"
)

type MerchantService interface {
	AddMerchant(merchant *models.Merchants) error
	GetMerchantByID(id uint) (models.Merchants, error)
	GetAllMerchants() ([]models.Merchants, error)
}

type MerchantServiceImpl struct{}

func (m *MerchantServiceImpl) AddMerchant(merchant *models.Merchants) error {
	return config.CreateMerchant(merchant)
}

// GetMerchantByID retrieves a merchant by ID
func (m *MerchantServiceImpl) GetMerchantByID(id uint) (models.Merchants, error) {
	return config.GetMerchantByID(id)
}

// GetAllMerchants retrieves all merchants
func (m *MerchantServiceImpl) GetAllMerchants() ([]models.Merchants, error) {
	return config.GetAllMerchants()
}

func NewMerchantService() MerchantService {
	return &MerchantServiceImpl{}
}
