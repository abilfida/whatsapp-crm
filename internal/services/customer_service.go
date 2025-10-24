package services

import (
	"whatsapp-crm/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerService struct {
	db *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{db: db}
}

// GetOrCreateCustomer gets existing customer or creates new one
func (cs *CustomerService) GetOrCreateCustomer(whatsappID, name string) (*models.Customer, error) {
	var customer models.Customer

	// Try to find existing customer
	if err := cs.db.Preload("Contact").First(&customer, "whatsapp_id = ?", whatsappID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new customer
			customer = models.Customer{
				WhatsAppID: whatsappID,
				Name:       name,
			}

			if err := cs.db.Create(&customer).Error; err != nil {
				return nil, err
			}

			// Create contact record
			contact := models.Contact{
				CustomerID: customer.ID,
				WhatsAppID: whatsappID,
				DisplayName: name,
				Status:     models.ContactStatusValid,
			}

			if err := cs.db.Create(&contact).Error; err != nil {
				return nil, err
			}

			customer.Contact = &contact
		} else {
			return nil, err
		}
	}

	return &customer, nil
}

// GetCustomers returns paginated list of customers
func (cs *CustomerService) GetCustomers(page, limit int, search string) ([]models.Customer, int64, error) {
	offset := (page - 1) * limit

	query := cs.db.Model(&models.Customer{}).Preload("Contact")

	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get customers
	var customers []models.Customer
	if err := query.Offset(offset).Limit(limit).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

// GetCustomerByID returns customer by ID
func (cs *CustomerService) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := cs.db.Preload("Contact").Preload("Conversations").First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// UpdateCustomer updates customer information
func (cs *CustomerService) UpdateCustomer(id uuid.UUID, updates map[string]interface{}) (*models.Customer, error) {
	var customer models.Customer
	if err := cs.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := cs.db.Model(&customer).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload with preloaded data
	if err := cs.db.Preload("Contact").First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

// DeleteCustomer soft deletes customer
func (cs *CustomerService) DeleteCustomer(id uuid.UUID) error {
	return cs.db.Delete(&models.Customer{}, "id = ?", id).Error
}