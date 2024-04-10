package jwtx

import (
	"fmt"
	"testing"
)

func TestJWT(t *testing.T) {
	claims := UserClaims{
		UserId:   1,
		UserName: "test",
	}

	token, err := GenToken(claims)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(token)
	fmt.Println(ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsIlVzZXJOYW1lIjoidGVzdCIsIlJvbGUiOiJhZG1pbiIsIlJvbGVJZCI6MSwiZXhwIjoxNjk5MzQ0MjIzLCJuYmYiOjE2OTkzNDA2MjN9.kRHxEQIlsIAbxXfe76-CvKpgVbSMF9LHxjdgZEBJUGg\n"))
}
