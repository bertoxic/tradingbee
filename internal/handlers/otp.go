package handler

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/bertoxic/tradingbee/internal/models"
	"github.com/bertoxic/tradingbee/internal/transport/httputil"
	"github.com/gin-gonic/gin"
)

const (
	maxOTPLength    = 8
	otpExpiryDur   = 5 * time.Minute
	maxOTPAttempts = 3
	encodingBase   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

func SendOTP() {}

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", err
	}

    otpNum := new(big.Int).SetBytes(buffer).String()
	log.Println(otpNum)
	if len(otpNum) < length {
		otpNum = fmt.Sprintf("%0*s", length, otpNum) 
	}
	log.Println(otpNum)
	otp := otpNum[:length]

	return otp,nil
}

func VerifyOTP(ctx *gin.Context) {
	type oTp struct {
		otpstring string
	}
	var otp oTp
	ctx.ShouldBind(otp)

}
func GenerateOTPResponse(ctx *gin.Context){
	authDetails := &models.AuthPayload{}
   ctx.ShouldBind(authDetails)

   otp , err := GenerateOTP(8)
   if err != nil {
	   return
   }
   type OTP struct {
	   Otp string `json:"otp"`
   }
   var otpx OTP
   otpx.Otp = otp
   httputil.WriteJson(ctx, true,200,&models.JsonResponse{
	   Success: true,
	   Message: "otp generated successfully",
	   Data: otpx,

   })
}


func verifyOTP(generatedOTP, userOTP string) bool {
	attempts:=0
	validUntil:= time.Now().Add(otpExpiryDur)

	for attempts < maxOTPAttempts && time.Now().Before(validUntil) {
		if userOTP == generatedOTP {
			return true
		}
		attempts++  
	}

	return false
}
