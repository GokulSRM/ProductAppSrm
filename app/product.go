package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct for storing data
type product struct {
	PId           string  `json:"pid"`
	CategoryId    string  `json:"categoryid"`
	SubCategoryId string  `json:"subcategoryid"`
	BrandId       string  `json:"brandid"`
	VarientId     string  `json:"varientid"`
	Pname         string  `json:"pname"`
	Pdesc         string  `json:"pdesc"`
	Pqty          int     `json:"pqty"`
	Pmrp          float32 `json:"pmrp"`
	Pprice        float32 `json:"pprice"`
}

var productCollection = db().Database("ProductApp").Collection("Product") // get collection "users" from db() which returns *mongo.Client

// Create Profile or Signup

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var prod product
	err := json.NewDecoder(r.Body).Decode(&prod) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := productCollection.InsertOne(context.TODO(), prod)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
	json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document

}

// Get Profile of a particular User by Name

func GetProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body product
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	var result primitive.M //  an unordered representation of a BSON document which is a Map
	err := productCollection.FindOne(context.TODO(), bson.D{{"pid", body.PId}}).Decode(&result)
	if err != nil {

		fmt.Println(err)

	}
	json.NewEncoder(w).Encode(result) // returns a Map containing document

}

//Update Profile of User

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		PId           int     `json:"pid"`           //value that has to be matched
		CategoryId    int     `json:"categoryid"`    // value that has to be modified
		SubCategoryId int     `json:"subcategoryid"` // value that has to be modified
		BrandId       int     `json:"brandid"`       // value that has to be modified
		VarientId     string  `json:"varientid"`     // value that has to be modified
		Pname         string  `json:"pname"`         // value that has to be modified
		Pdesc         string  `json:"pdesc"`         // value that has to be modified
		Pqty          int     `json:"pqty"`          // value that has to be modified
		Pmrp          float32 `json:"pmrp"`          // value that has to be modified
		Pprice        float32 `json:"pprice"`        // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"pid", body.PId}} // converting value to BSON type
	after := options.After              // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"categoryid", body.CategoryId}, {"varientid", body.VarientId}, {"subcategoryid", body.SubCategoryId}, {"brandid", body.BrandId}, {"pname", body.Pname}, {"pdesc", body.Pdesc}, {"pqty", body.Pqty}, {"pmrp", body.Pmrp}, {"pprice", body.Pprice}}}}
	updateResult := productCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	// update1 := bson.D{{"$set", bson.D{{"age", body.Age}}}}
	// updateResult1 := productCollection.FindOneAndUpdate(context.TODO(), filter, update1, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
}

//Delete Profile of User

func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res, err := productCollection.DeleteOne(context.TODO(), bson.D{{"pid", params}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted

}

func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                      //slice for multiple documents
	cur, err := productCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	json.NewEncoder(w).Encode(results)
}
