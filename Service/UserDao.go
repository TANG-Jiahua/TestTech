package Service

import (
	"TestTech/Model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"reflect"
)

func copyTo(src, dst interface{}) {
	sval := reflect.ValueOf(src).Elem()
	dval := reflect.ValueOf(dst).Elem()

	for i := 0; i < sval.NumField(); i++ {
		value := sval.Field(i)
		name := sval.Type().Field(i).Name


		dvalue := dval.FieldByName(name)
		if dvalue.IsValid() == false {
			continue
		}
		t := sval.Type().Field(i).Type.Kind()
		switch t {
		case reflect.String:
			if value.String() == "" {
				continue
			}
		case reflect.Int:
			if value.Int() == 0 {
				continue
			}
		case reflect.Float64:
			if value.Float() == 0{
				continue
			}
		case reflect.Array:
			if value.Len()==0{
				continue
			}
		//case reflect.Bool

		}
		dvalue.Set(value)
	}
}

func GeneratePassWd(src string) (pw string, err error) {
	data := []byte(src)
	hsh, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	pw = string(hsh)
	return
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}

func UserExist(id string) (bool, Model.User) {
	var result Model.User
	err := MyUser.DB.FindOne(context.Background(), bson.D{{"id", id}}).Decode(&result)
	//log.Println(result)
	if  err!= nil{
		log.Println(err)
		return false, Model.User{}
	}else {
		return true,result
	}

}