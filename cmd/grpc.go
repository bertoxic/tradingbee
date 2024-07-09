package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	config "github.com/bertoxic/tradingbee/configs"
	"github.com/bertoxic/tradingbee/internal/handlers"

	"github.com/bertoxic/tradingbee/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/gin-gonic/gin"
	// config "github.com/bertoxic/med/services/authentication/configs"
	// "github.com/bertoxic/med/services/authentication/grpc"
	// handler "github.com/bertoxic/med/services/authentication/internal/handlers"
	// "github.com/bertoxic/med/services/authentication/internal/models"
	// "github.com/dgrijalva/jwt-go"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// gp "google.golang.org/grpc"
	// "google.golang.org/grpc/keepalive"
)

type UserServer struct {
}

var db *mongo.Client
var conf *config.AppConfig
var dataMed *mongo.Database

type DataB struct {
	app *config.AppConfig
}

func NewDataB(appconfig *config.AppConfig) {
	db = appconfig.Client
	conf = appconfig
	dataMed = db.Database(appconfig.Config.DBNAME)
}

func (s *UserServer) RegisterUser(ctx *gin.Context) (*models.JsonResponse, error) {
//	var JsonResponse models.JsonResponse
	log.Println("in reg user now")
	user := &models.UserDetails{}
	ctx.ShouldBind(user)

	log.Println(user)
	patientCollection := dataMed.Collection("user")
	ctxs, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	emailExist, err := emailExists(user.Email, dataMed)
	if err != nil {
		error := models.ErrorJson{
			Code:    500,
			Message: "something went wrong during signup",
		}
		errstr, err := json.Marshal(error)
		if err != nil {
			log.Println(err)
		}

		return &models.JsonResponse{
			Success: false,
			Message: "registation unsuccessful" + err.Error(),
			Error:   &models.ErrorJson{Message: string(errstr)},
		}, err
	}
	if emailExist {
		var errorJson = models.ErrorJson{
			Code:    http.StatusBadRequest,
			Message: errors.New("this email has already been registered please signin").Error(),
		}
		jsondat, err := json.MarshalIndent(errorJson, "", " ")
		if err != nil {
			log.Println(err)
		}
		_ = string(jsondat)

		return &models.JsonResponse{
			Success: false,
			Message: "registation unsuccessful",
			Error:   &errorJson,
		}, nil
	}
	log.Println("trying too register")
	bsoncmd := bson.M{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.PassWord,
		"role":       user.UserType,
	}
	result, err := patientCollection.InsertOne(ctxs, bsoncmd)
	if err != nil {
		return &models.JsonResponse{
			Success: false,
			Message: "registation successful",
			Error:   &models.ErrorJson{Code: 400, Message: err.Error()},
		}, nil
	}
	insertedId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to assert type of insertedIdzzzzzzzz")
	}
	insertedIdStr := insertedId.Hex()
	data := map[string]interface{}{
		"id": insertedIdStr,
	}
	jsondat, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return &models.JsonResponse{
			Success: false,
			Message: "registation successful",
			Data:    string(jsondat),
			Error:   &models.ErrorJson{Code: 400, Message: err.Error()},
		}, nil
	}
	return &models.JsonResponse{
		Success: true,
		Message: "registation successful",
		Data:    string(jsondat),
	}, nil
}

func (us *UserServer) LoginUser(ctx *gin.Context) (*models.JsonResponse, error) {
	log.Println("in login user now")

	user := models.UserDetails{}
	ctx.ShouldBind(user)
	//Jsondata := &models.JsonResponse{}

	emailExist, err := emailExists(user.Email, dataMed)
	if err != nil {
		error := models.ErrorJson{
			Code:    500,
			Message: "something went wrong during signup",
		}
		errstr, err := json.Marshal(error)
		if err != nil {
			log.Println(err)
		}
		return &models.JsonResponse{
			Success: false,
			Message: "login unsuccessful" + err.Error(),
			Error:   &models.ErrorJson{Code: 400, Message: string(errstr)},
		}, err
	}
	if !emailExist {
		var errorJson = models.ErrorJson{
			Code:    http.StatusBadRequest,
			Message: errors.New("this email is not registered, please register").Error(),
		}
		jsondat, err := json.MarshalIndent(errorJson, "", " ")
		if err != nil {
			log.Println(err)
		}
		error := string(jsondat)

		return &models.JsonResponse{
			Success: false,
			Message: "login unsuccessful",
			Error:   &models.ErrorJson{Code: 400, Message: error},
		}, nil
	}

	patientCollection := dataMed.Collection("patients")
	ctxs, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	bsonCmd := bson.M{
		"email":    user.Email,
		"password": user.PassWord,
	}

	userdetail := models.UserDetails{}
	cursor := patientCollection.FindOne(ctxs, bsonCmd)
	cursor.Decode(&userdetail)
	if user.Email != userdetail.Email && user.PassWord != userdetail.PassWord {
		var error = models.ErrorJson{
			Code:    http.StatusBadRequest,
			Message: "email or password is incorrect",
		}
		var jsondata = models.JsonResponse{
			Success: false,
			Message: "",
			Error:   &error,
		}
		jsondat, err := json.Marshal(jsondata)
		if err != nil {
			log.Println(err)
		}
		data := string(jsondat)

		return &models.JsonResponse{
			Success: false,
			Message: "login unsuccessful",
			Data:    data,
		}, nil
	}
	user_person := models.User{}
	if userdetail.UserType == "patient" {

		cursor.Decode(&user_person)
	}
	token, err := handler.GenerateToken(userdetail)

	if err != nil {

		return &models.JsonResponse{
			Success: false,
			Message: "unable to generate token",
			Error: &models.ErrorJson{Code: 400, Message: err.Error()},
		}, nil
	}
	// datax := struct{token handler.Tokens; user models.UserDetails;}{
	// 	token: token,
	// 	user: user ,
	// }
	data := map[string]interface{}{
		"token": token,
		"user":  user_person,
	}
	var jsondata = models.JsonResponse{
		Success: true,
		Message: "",
		Data:    data,
	}

	jsondat, err := json.Marshal(jsondata)
	if err != nil {
		log.Println(err)
	}
	dat := string(jsondat)
	return &models.JsonResponse{
		Success: true,
		Message: "login successful",
		Data:    dat,
	}, nil

}

func emailExists(email string, dataB *mongo.Database) (bool, error) {
	user := models.UserDetails{}
	patientCollection := dataB.Collection("patients")
	ctxs, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"email": email}
	cursor := patientCollection.FindOne(ctxs, filter)

	err := cursor.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GenerateToken(userDetails models.UserDetails) (models.Tokens, error) {
	var tokens models.Tokens
	claims := &models.SignedDetails{
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Email:     userDetails.Email,
		UserType:  userDetails.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 25).Unix(),
		},
	}
	refreshClaims := &models.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 196).Unix(),
		},
	}
	secret_key := []byte("bert")
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.StandardClaims).SignedString(secret_key)
	if err != nil {
		return tokens, err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshClaims.StandardClaims).SignedString(secret_key)
	if err != nil {
		return tokens, err
	}
	tokens = models.Tokens{
		Token:        token,
		RefreshToken: refreshtoken,
	}

	return tokens, nil
}

func ValidateToken(signedToken string) (*models.SignedDetails, string) {
	var userclaims models.SignedDetails

	token, err := jwt.ParseWithClaims(signedToken, &userclaims.StandardClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte("bert"), nil
	})
	if err != nil {

		return nil, "cannot parse token"
	}
	claims, ok := token.Claims.(*models.SignedDetails)
	if !ok {
		return nil, "invalid token"
	}
	if claims.StandardClaims.ExpiresAt < time.Now().Local().Unix() {
		return nil, "expired token"
	}
	return claims, ""
}

// func grpcListen() {
// 	port := os.Getenv("GRPC_PORT")
// 	if port == "" {
// 		port = "10000" // Default to 10000 if GRPC_PORT is not set
// 	}

// 	lis, err := net.Listen("tcp4", "0.0.0.0:"+port)
// 	if err != nil {
// 		log.Fatalf("failed to listen on port %v: %v", port, err)
// 	}
// 	s := gp.NewServer(
// 		gp.KeepaliveParams(keepalive.ServerParameters{
// 			MaxConnectionIdle:     15 * time.Second,
// 			MaxConnectionAge:      30 * time.Second,
// 			MaxConnectionAgeGrace: 5 * time.Second,
// 			Time:                  5 * time.Second,
// 			Timeout:               1 * time.Second,
// 		}),
// 		gp.MaxRecvMsgSize(1024*1024*50), // 10MB
// 		gp.MaxSendMsgSize(1024*1024*50),
// 	)

// 	log.Printf("running on port..>> %s", port)
// 	grpc.RegisterUserAuthServiceServer(s, &UserServer{})
// 	log.Printf("server listening at %v on port %v", lis.Addr(), port)

// 	if err := s.Serve(lis); err != nil {
// 		log.Printf("failed to serve: %v", err)
// 	}
// }

// func grpcListen(){
// 	port := os.Getenv("GRPC_PORT")
//     if port == "" {
//         port = "10000" // Default to 10000 if PORT is not set
//     }

//     lis, err := net.Listen("tcp", "0.0.0.0:"+port)
//     if err != nil {
//         log.Fatalf("failed to listen to: %v, error is : %v", port,err)
//     }
// 	//lis, err := net.Listen("tcp", ":5001")
// 	// if err != nil {
// 	// 	log.Printf("did not listen failed dto isten: %v", err)
// 	// }
// 	s := gp.NewServer()
// 	log.Printf("running on port..>> %s",port)
// 	grpc.RegisterUserAuthServiceServer(s, &UserServer{})
// 	log.Printf("server listening at %v on port %v", lis.Addr(),port)
// 	if err := s.Serve(lis); err != nil {
// 		log.Printf("failed to serve: %v", err)
// 	}

// }
