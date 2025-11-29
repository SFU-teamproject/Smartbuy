package mockstorage

import (
	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetSmartphone(ID int) (models.Smartphone, error) {
	args := m.Called(ID)
	return args.Get(0).(models.Smartphone), args.Error(1)
}

func (m *MockStorage) GetSmartphones() ([]models.Smartphone, error) {
	args := m.Called()
	return args.Get(0).([]models.Smartphone), args.Error(1)
}

func (m *MockStorage) GetSmartphonesByIDs(IDs []int) ([]models.Smartphone, error) {
	args := m.Called(IDs)
	return args.Get(0).([]models.Smartphone), args.Error(1)
}

func (m *MockStorage) GetUser(ID int) (models.User, error) {
	args := m.Called(ID)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockStorage) GetUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockStorage) GetUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockStorage) CreateUser(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockStorage) UpdateUser(userID int, updates map[string]any) (models.User, error) {
	args := m.Called(updates)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockStorage) CreateTmpPassword(token models.TmpPassword) (models.TmpPassword, error) {
	args := m.Called(token)
	return args.Get(0).(models.TmpPassword), args.Error(1)
}

func (m *MockStorage) GetTmpPassword(email string) (models.TmpPassword, error) {
	args := m.Called(email)
	return args.Get(0).(models.TmpPassword), args.Error(1)
}

func (m *MockStorage) DeleteTmpPassword(email string) (models.TmpPassword, error) {
	args := m.Called(email)
	return args.Get(0).(models.TmpPassword), args.Error(1)
}

func (m *MockStorage) GetReview(ID int) (models.Review, error) {
	args := m.Called(ID)
	return args.Get(0).(models.Review), args.Error(1)
}

func (m *MockStorage) GetReviews(smartphoneID int) ([]models.Review, error) {
	args := m.Called(smartphoneID)
	return args.Get(0).([]models.Review), args.Error(1)
}

func (m *MockStorage) CreateReview(review models.Review) (models.Review, error) {
	args := m.Called(review)
	return args.Get(0).(models.Review), args.Error(1)
}

func (m *MockStorage) UpdateReview(review models.Review) (models.Review, error) {
	args := m.Called(review)
	return args.Get(0).(models.Review), args.Error(1)
}

func (m *MockStorage) DeleteReview(ID int) (models.Review, error) {
	args := m.Called(ID)
	return args.Get(0).(models.Review), args.Error(1)
}

func (m *MockStorage) GetCarts() ([]models.Cart, error) {
	args := m.Called()
	return args.Get(0).([]models.Cart), args.Error(1)
}

func (m *MockStorage) GetCart(ID int) (models.Cart, error) {
	args := m.Called(ID)
	return args.Get(0).(models.Cart), args.Error(1)
}

func (m *MockStorage) GetCartByUserID(userID int) (models.Cart, error) {
	args := m.Called(userID)
	return args.Get(0).(models.Cart), args.Error(1)
}

func (m *MockStorage) GetCartItem(ID int) (models.CartItem, error) {
	args := m.Called(ID)
	return args.Get(0).(models.CartItem), args.Error(1)
}

func (m *MockStorage) GetCartItems(cartID int) ([]models.CartItem, error) {
	args := m.Called(cartID)
	return args.Get(0).([]models.CartItem), args.Error(1)
}

func (m *MockStorage) AddToCart(cartItem models.CartItem) (models.CartItem, error) {
	args := m.Called(cartItem)
	return args.Get(0).(models.CartItem), args.Error(1)
}

func (m *MockStorage) SetQuantity(cartItem models.CartItem) (models.CartItem, error) {
	args := m.Called(cartItem)
	return args.Get(0).(models.CartItem), args.Error(1)
}

func (m *MockStorage) DeleteFromCart(cartID, itemID int) (models.CartItem, error) {
	args := m.Called(cartID, itemID)
	return args.Get(0).(models.CartItem), args.Error(1)
}
