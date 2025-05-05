package utils

import (
	"fmt"
	"math/rand"
)

func GenerateOTP() string {
	otp := rand.Intn(999999)        // Generate a random 6-digit OTP
	return fmt.Sprintf("%06d", otp) // Format it as a string with leading zeros
}
