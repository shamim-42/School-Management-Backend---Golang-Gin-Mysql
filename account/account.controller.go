package account

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/saiyedulbas/second/dbconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	JwtSigningKey = "lskjfo58493mnhgffghwokdcn##@#$%^&vbny654edfbngnfre9r8y"
)

var (
	db *gorm.DB = dbconn.SetupConnection()
)

func Signin(c *gin.Context) {
	payload := SigninPayload{}
	jsonParsedBody := c.BindJSON(&payload)
	_ = jsonParsedBody

	email := payload.Email
	password := payload.Password

	users := []User{}
	result := db.Where("email=?", email).Find(&users)

	if result.RowsAffected > 0 {
		storedPassword := users[0].Password
		bytePassword := []byte(storedPassword)
		bytePassword2 := []byte(password)
		err := bcrypt.CompareHashAndPassword(bytePassword, bytePassword2)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "Login Error!",
				"message": "Please try with valid email and password",
			})
		} else {
			// Toekn creation
			createdToken, err := CreateJwtToken([]byte(JwtSigningKey), users[0].ID, users[0].Nickname, users[0].Email, users[0].Mobile)
			if err != nil {
				c.JSON(502, gin.H{
					"status":  "Something went wrong!",
					"message": err,
				})
			} else {
				fmt.Print(users[0])
				c.JSON(200, gin.H{
					"status":  "Success",
					"message": "Loggedin Successfully!",
					"data":    users,
					"jwt":     createdToken,
				})
			}
		}
	} else {
		c.JSON(200, gin.H{
			"status":  "Login error!",
			"message": "Email not found in the database",
		})
	}
}

func Signup(c *gin.Context) {
	payload := SignupPayload{}
	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// converting plain password into hashed one.
	password := []byte(payload.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	stringedHashedPassword := string(hashedPassword)
	payload.Password = stringedHashedPassword

	users := []User{}
	result := db.Where("Email = ?", payload.Email).Find(&users)

	if result.RowsAffected > 0 {
		c.JSON(200, gin.H{
			"status":  "Signup error!",
			"message": "Email already exists",
		})

	} else {
		user := User{
			Nickname: payload.Nickname,
			Email:    payload.Email,
			Mobile:   payload.Mobile,
			Password: stringedHashedPassword,
		}
		instance := db.Create(&user)

		if instance.RowsAffected > 0 {
			c.JSON(201, gin.H{
				"status":  "Signup successfull",
				"message": "User created successfully!",
			})
		} else {
			c.JSON(502, gin.H{
				"status":  "Signup error!",
				"message": "Something went wrong. Please contact administrator!",
			})
		}
	}

	_ = result
	_ = err
}

func CreateRole(c *gin.Context) {
	payload := Role{}
	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	roles := []Role{}
	result := db.Where("Role=?", payload.Role).Find(&roles)
	if result.RowsAffected > 0 {
		c.JSON(200, gin.H{
			"status_code": 200,
			"message":     "Role already exist",
		})
	} else {
		instance := db.Create(&payload)
		fmt.Println("New role created")
		fmt.Println(instance)
		fmt.Println("\n")
		if instance.RowsAffected > 0 {
			c.JSON(http.StatusCreated, gin.H{
				"status_code": 201,
				"message":     "Role created successfully",
				// "data":        JSON(instance),
				// "data":        instance,
			})
		} else {
			c.JSON(500, gin.H{
				"status_code": 500,
				"message":     "Internal Server Error",
			})
		}
	}
}

func ListRole(c *gin.Context) {
	roles := []Role{}
	query := db.Preload("Role").Find(&roles)
	_ = query
	c.JSON(http.StatusAccepted, gin.H{
		"status_cde": http.StatusAccepted,
		"message":    "All role list",
		"data":       &roles,
	})
}

func GetRole(c *gin.Context) {
	role := Role{}
	id := c.Param("id")
	query := db.Preload("Role").First(&role, id)
	_ = query
	c.JSON(http.StatusAccepted, gin.H{
		"status_cde": http.StatusAccepted,
		"message":    "Details of the Role",
		"data":       &role,
	})
}

func UpdateRole(c *gin.Context) {
	role := Role{}
	if err := c.BindJSON(&role); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	id := c.Param("id")
	query := db.Model(&Role{}).Where("id=?", id).Update("role", role.Role)
	_ = query
	c.JSON(http.StatusAccepted, gin.H{
		"status_cde": http.StatusAccepted,
		"message":    "Details of the Role",
		"data":       &role,
	})
}

func UpdateUser(c *gin.Context) {
	user := &UpdateUserPayload{}
	if err := c.BindJSON(&user); err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusAccepted, gin.H{
			"status_code": http.StatusBadRequest,
			"message":     "Something went error!",
			"error":       err,
		})
	}
	id := c.Param("id")
	db.Model(&User{}).Where("id=?", id).Updates(&user)

	c.JSON(http.StatusAccepted, gin.H{
		"status_code": http.StatusAccepted,
		"message":     "User updated successfully!",
	})
}

func CreateJwtToken(mySigningKey []byte, id uint32, name string, email string, mobile string) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Add some claims
	claims := make(jwt.MapClaims)
	claims["user_id"] = id
	claims["name"] = name
	claims["email"] = email
	claims["mobile"] = mobile
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	token.Claims = claims

	//Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}

func ParseJwtToken(myToken string, key string) bool {
	token, err := jwt.Parse(
		myToken,
		func(myToken *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
	)
	var validationMessage bool
	if err == nil && token.Valid {
		validationMessage = true
	} else {
		validationMessage = false
	}
	return validationMessage
}

// // This is a custom API that had been developed for onetime use.
// // Converting image urls from S3 to CloudFront. We can remove this function from here anytime
// func ConvertUrl(c *gin.Context) {
// 	for i := 2; i < 71; i++ {
// 		photo := Photos{}
// 		db.Where("id=?", i).Find(&photo)
// 		var url string = photo.Url
// 		var new_url string = "https://d2u47s18r13gjl.cloudfront.net" + url[52:]
// 		// var data Photos;
// 		// data.Photographer = photo.Photographer
// 		// data.Device = photo.Device
// 		// data.Location = photo.Location
// 		// data.Url = new_url
// 		// data.Caption = photo.Caption
// 		// data.Handler = photo.Handler>>
// 		photo.Url = new_url
// 		result := db.Model(&photo).Where("id=?", i).Updates(&photo)
// 		if result.Error != nil {
// 			fmt.Print("Bismillah")
// 		}
// 		if result.RowsAffected > 0 {
// 			c.JSON(200, gin.H{"data": photo})
// 		}
// 	}
// }
