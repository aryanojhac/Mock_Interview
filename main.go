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

// Defining test values as required
var testUsers = []User{
	{Name: "Aryan", Password: "Aryan@30", Phone: "1234567890"},
	{Name: "Atharv", Password: "Atharv%03", Phone: "0986543211"},
	{Name: "Shravani", Password: "shanu@15", Phone: "1223341455"},
}

var validate *validator.Validate

func validateUser() {
	validate = validator.New()
	// Validating the user Name here
	validate.RegisterValidation("userName", func(f1 validator.FieldLevel) bool {
		name := f1.Field().String()
		match, _ := regexp.MatchString(`^[A-Za-z]+$`, name)
		return match
	})

	// Validating password here
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		match, _ := regexp.MatchString(`^[A-Za-z\d@$!%*?&]*$`, password)
		// At least 1 uppercase
		upper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		// At least 1 digit
		digit := regexp.MustCompile(`[0-9]`).MatchString(password)
		// At least 1 special char
		special := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

		return match &&upper && digit && special
	})
}

// Checking the input from postman (user)
func checkCredentials(name, password, phone string) bool {
	for _, u := range testUsers {
		if u.Name == name && u.Password == password && u.Phone == phone {
			return true
		}
	}
	return false
}

func signIn(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "binding error"})
		return
	}
	// Validating Structure using go-play...
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": "the input is not valid"})
		return
	}

	if !checkCredentials(user.Name, user.Password, user.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Name or Password or Phone Number"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Name": user.Name, "message": "Signed In Successfully"})
}

func main() {
	router := gin.Default()

	validateUser()

	router.POST("/signIn", signIn)

	router.Run(":8080")
}
