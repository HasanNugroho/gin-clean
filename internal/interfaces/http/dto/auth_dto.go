package dto

type (
	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	AuthResponse struct {
		Token        string      `json:"token"`
		RefreshToken string      `json:"refresh_token"`
		Data         interface{} `json:"data"`
	}

	RenewalTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)
