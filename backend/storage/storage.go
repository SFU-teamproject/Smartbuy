package storage

import (
	"github.com/sfu-teamproject/smartbuy/backend/models"
)

type Storage interface {
	GetSmartphone(ID int) (models.Smartphone, error)
	GetSmartphones() ([]models.Smartphone, error)
	GetSmartphonesByIDs(IDs []int) ([]models.Smartphone, error)

	GetUser(ID int) (models.User, error)
	GetUsers() ([]models.User, error)
	GetUserByName(name string) (models.User, error)
	CreateUser(user models.User) (models.User, error)

	GetReview(ID int) (models.Review, error)
	GetReviews(smartphoneID int) ([]models.Review, error)
	CreateReview(review models.Review) (models.Review, error)
	UpdateReview(review models.Review) (models.Review, error)
	DeleteReview(ID int) (models.Review, error)

	GetCarts() ([]models.Cart, error)
	GetCart(ID int) (models.Cart, error)
	GetCartByUserID(userID int) (models.Cart, error)

	GetCartItem(ID int) (models.CartItem, error)
	GetCartItems(cartID int) ([]models.CartItem, error)
	AddToCart(cartItem models.CartItem) (models.CartItem, error)
	SetQuantity(cartItem models.CartItem) (models.CartItem, error)
	DeleteFromCart(cartID, itemID int) (models.CartItem, error)
}
