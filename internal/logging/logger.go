package logging

import (
	"github.com/rs/zerolog/log"
)

// LogLoginAttempt logs a user's login attempt.
func LogLoginAttempt(email string, success bool) {
	if success {
		log.Info().Str("email", email).Msg("Successful login attempt")
	} else {
		log.Warn().Str("email", email).Msg("Failed login attempt")
	}
}

// LogPasswordChange logs a password change event.
func LogPasswordChange(email string) {
	log.Info().Str("email", email).Msg("Password changed successfully")
}

// LogOtpVerification logs the result of an OTP verification.
func LogOtpVerification(email string, success bool) {
	if success {
		log.Info().Str("email", email).Msg("OTP verification successful")
	} else {
		log.Warn().Str("email", email).Msg("OTP verification failed")
	}
}
