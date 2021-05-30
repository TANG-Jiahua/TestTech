package Service

import (
	"TestTech/Model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type UserModel struct {
	DB * mongo.Collection
}
var MyUser *UserModel

func AddUser(c *gin.Context) {
	x,_:=ioutil.ReadAll(c.Request.Body)
	fmt.Printf("%s \n",string(x))
	users := make([]Model.User,0)
	err := json.Unmarshal(x, &users)
	if err != nil {
		fmt.Printf("%s",err)
	}

	for _,u := range users{
		var exist,_ =UserExist(u.Id)
		fmt.Println("\n",exist)
		if !exist {
			var pw,_ = GeneratePassWd(u.Password)
			u.Password = pw
			//fmt.Printf("tags: ",u.Tags)
			MyUser.DB.InsertOne(context.Background(),u)
			c.String(http.StatusOK,"\nSuccessfully add: "+u.Id)
			fmt.Println(u)
		}else {
			c.String(http.StatusOK,"\nError message : Model exists!")
		}

		f,err:=os.Create("DataFichier/"+u.Id+".txt")
		ioutil.WriteFile("DataFichier/"+u.Id+".txt",[]byte(u.Data),0644)
		if err!=nil{
			fmt.Println(err)
		}
		defer f.Close()


	}
}

func Login(c *gin.Context){
	json := make(map[string]interface{})
	err := c.BindJSON(&json)
	if err != nil {
		return
	}
	fmt.Printf("%v",&json)
	fmt.Printf(json["id"].(string))
	var exist,result = UserExist(json["id"].(string))

	if json["id"] == nil{
		c.String(http.StatusOK, "Error message : no id!")
	}else {
		if exist {
			pswd := json["password"].(string)
			if ComparePasswords(result.Password, pswd) {
				c.String(http.StatusOK, "Successful!")
			} else {
				c.String(http.StatusOK, "Error message : password error!")
			}
		} else {
			c.String(http.StatusOK, "Error message : Model Id not found!")
		}
	}

}

func DeleteUser(c *gin.Context){
	id:=c.Param("id")

	var exist,_ = UserExist(id)
	if exist{
		deleteResult, err := MyUser.DB.DeleteOne(context.Background(), bson.M{"id": id})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("collection.DeleteOne:", deleteResult)
		err2 := os.Remove("DataFichier/"+id+".txt")
		if err!=nil {
			fmt.Println(err2)
		}
		c.String(http.StatusOK,"Delete successfully!")

	}else{
		c.String(http.StatusOK,"Error message : Model id not found!")
	}
}

func GetUserList(c *gin.Context)  {
	cur, err := MyUser.DB.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	var all []*Model.User
	err = cur.All(context.Background(), &all)
	if err != nil {
		log.Fatal(err)
	}
	cur.Close(context.Background())

	log.Println("collection.Find curl.All: ", all)
	c.String(http.StatusOK,"Voici la liste d'id\n")
	for _, one := range all {
		c.String(http.StatusOK,one.Id+"\n")
	}

	//jsonBytes, err := json.Marshal(all)
	//c.String(http.StatusOK,string(jsonBytes))
}

func GetUserById(c *gin.Context){
	id:=c.Param("id")
	log.Println("id:",id)
	var exist,result = UserExist(id)
	if exist{
		jsonBytes, err := json.Marshal(result)
		log.Println(string(jsonBytes))
		if err!= nil{
			fmt.Println(err)
		}
		c.String(http.StatusOK,string(jsonBytes))
	}else {
		c.String(http.StatusOK,"Error message: Model id not found!")
	}
}


func UpdateUser (c *gin.Context){
	id:=c.Param("id")
	con,_:=ioutil.ReadAll(c.Request.Body)
	//log.Println("context:",con)
	inputUser := new(Model.User)
	json.Unmarshal(con,&inputUser)
	log.Println("Input:",inputUser)

	result := new(Model.User)
	err := MyUser.DB.FindOne(context.Background(), bson.D{{"id", id}}).Decode(&result)
	log.Println("result:",result)
	if err != nil {
		fmt.Printf("%s",err)
		if err==mongo.ErrNoDocuments{
			c.String(http.StatusOK,"Error message : User id not found!")
		}
	}else{
		copyTo(inputUser, result)
		update := bson.M{"$set":result}
		updateResult,err := MyUser.DB.UpdateOne(context.Background(),bson.M{"id":result.Id},update)
		if err!= nil{
			log.Fatal(err)
		}
		log.Println("collection.updateOne:",updateResult)
		c.String(http.StatusOK,"Update successfully!")
	}


}
