package main

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Defining structure of the User
type User struct {
	Name     string `json:"name" binding:"required,min=2,max=100" validate:"userName"`
	Password string `json:"password" binding:"required,min=6,max=100" validate:"password"`
	Phone    string `json:"phone" binding:"required,min=10,max=13"`
}

// Test values
var testUsers = []User{
	{Name: "Aryan", Password: "Aryan@30", Phone: "1234567890"},
	{Name: "Atharv", Password: "Atharv%03", Phone: "0986543211"},
	{Name: "Shravani", Password: "shanu@15", Phone: "1223341455"},
}

var validate *validator.Validate

// Username validation: only letters
func userNameValidator(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if name == "" {
		return false
	}
	match, _ := regexp.MatchString(`^[A-Za-z]+$`, name)
	return match
}

// Password validation: allowed chars, at least 1 uppercase, 1 digit, 1 special char
func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return false
	}
	match, _ := regexp.MatchString(`^[A-Za-z\d@$!%*?&]*$`, password)
	upper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	digit := regexp.MustCompile(`[0-9]`).MatchString(password)
	special := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

	return match && upper && digit && special
}

// Check hardcoded credentials
func checkCredentials(name, password, phone string) bool {
	for _, u := range testUsers {
		if u.Name == name && u.Password == password && u.Phone == phone {
			return true
		}
	}
	return false
}

// SignIn handler
func signIn(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "binding error"})
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": "the input is not valid"})
		return
	}

	if !checkCredentials(user.Name, user.Password, user.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong Name or Password or Phone Number"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Name": user.Name, "message": "Signed In Successfully"})
}

func main() {
	router := gin.Default()

	// Registering custom rules
	validate = validator.New()
	validate.RegisterValidation("userName", userNameValidator)
	validate.RegisterValidation("password", passwordValidator)

	router.POST("/signIn", signIn)

	router.Run(":8080")
}
