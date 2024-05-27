package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func StartServer(wallet *Application) {
	wallet.FiberApp.Post("/payment", wallet.payment)
	wallet.FiberApp.Post("/deduct", wallet.deduct)
	if err := wallet.FiberApp.Listen(":3000"); err != nil {
		wallet.ErrorLog.Fatalf("Failed to start server: %v", err)
	} else {
		wallet.InfoLog.Println("Started Fiber API Server")
	}
}

func (wallet *Application) deduct(c *fiber.Ctx) error {
	var PassengerID struct {
		ID string `json:"id"`
	}
	if err := c.BodyParser(&PassengerID); err != nil {
		wallet.ErrorLog.Printf("Cannot parse JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var passenger Passenger
	result := wallet.Config.DB.First(&passenger, "id = ?", PassengerID.ID)

	if result.Error != nil {
		// If passenger not found, create a new one
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No Money, Top Up FIRST",
			})
		} else {
			// Other database error
			wallet.ErrorLog.Printf("Database error: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
	} else {
		// If passenger exists, update balance
		newBalance := passenger.Balance - 15
		if newBalance < 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Low Balance, please top up",
			})
		}
		if err := wallet.Config.DB.Model(&passenger).Update("balance", newBalance).Error; err != nil {
			wallet.ErrorLog.Printf("Failed to update balance: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update balance",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment 15rs Successful",
	})
}

func (wallet *Application) payment(c *fiber.Ctx) error {
	// Parse request body
	var paymentInfo struct {
		PaymentID string `json:"id"`
		TopUp     int    `json:"amount"`
	}
	if err := c.BodyParser(&paymentInfo); err != nil {
		wallet.ErrorLog.Printf("Cannot parse JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if passenger exists
	var passenger Passenger
	result := wallet.Config.DB.First(&passenger, "id = ?", paymentInfo.PaymentID)
	if result.Error != nil {
		// If passenger not found, create a new one
		if result.Error == gorm.ErrRecordNotFound {
			newPassenger := Passenger{
				ID:      paymentInfo.PaymentID,
				Balance: paymentInfo.TopUp,
			}
			if err := wallet.Config.DB.Create(&newPassenger).Error; err != nil {
				wallet.ErrorLog.Printf("Failed to create passenger: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create passenger",
				})
			}
		} else {
			// Other database error
			wallet.ErrorLog.Printf("Database error: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
	} else {
		// If passenger exists, update balance
		newBalance := passenger.Balance + paymentInfo.TopUp
		if err := wallet.Config.DB.Model(&passenger).Update("balance", newBalance).Error; err != nil {
			wallet.ErrorLog.Printf("Failed to update balance: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update balance",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment Successful",
	})
}
