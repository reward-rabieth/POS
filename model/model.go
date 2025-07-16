// Package model has domain models used throughout the application.
package model

// Thing with a name.
type Item struct {
	ID          string  `json:"id"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CreaatedAt  string  `json:"created"`
	UpdatedAt   string  `json:"updated"`
}

type Order struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Items     []Item  `json:"items"`
	TotalCost float64 `json:"totalcost"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	Updated   string  `json:"updated_at"`
	Expand    struct {
		User User `json:"user_id"`
	} `json:"expand"`
}

type LoginRequest struct {
	Identity string `json:"identity"` // Could be username or email
	Password string `json:"password"`
}

type User struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Username        string `json:"username"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	EmailVisibility bool   `json:"emailVisibility"`
	Verified        bool   `json:"verified"`
	Avatar          string `json:"avatar"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
