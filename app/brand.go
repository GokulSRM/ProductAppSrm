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
type brand struct {
	Bid         string `json:"bid"`
	ScId        string `json:"scid"`
	CId         string `json:"cid"`
	Bname       string `json:"bname"`
	Bdesc       string `json:"bdesc"`
	Bcreatedby  string `json:"bcreatedby"`
	Bmodifiedby string `json:"bmodifiedby"`
	Bstatus     bool   `json:"bstatus"`
}

var brandCollection = db().Database("ProductApp").Collection("Brand") // get collection "users" from db() which returns *mongo.Client

// Create Brand

func CreateBrand(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var brand brand
	err := json.NewDecoder(r.Body).Decode(&brand) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(400)
	}
	insertResult, err := brandCollection.InsertOne(context.TODO(), brand)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
	} else {
		fmt.Println("Inserted a single document: ", insertResult)
		json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document
		w.WriteHeader(200)
	}

}

// Get Brand

func GetBrand(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	// var body brand
	// e := json.NewDecoder(r.Body).Decode(&body)
	// if e != nil {

	// 	fmt.Print(e)
	// }
	var result primitive.M //  an unordered representation of a BSON document which is a Map
	filter := bson.M{"bid": params}
	err := brandCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {

		fmt.Println(err)
		json.NewEncoder(w).Encode("No data found!")

	} else {
		json.NewEncoder(w).Encode(result) // returns a Map containing document
		w.WriteHeader(200)
	}

}

// Get All Brand

func GetAllBrand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                    //slice for multiple documents
	cur, err := brandCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)

		w.WriteHeader(400)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(500)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(200)

}

//Update Brand

func UpdateBrand(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		Bid         string `json:"bid"`         //value that has to be matched
		ScId        string `json:"scid"`        // value that has to be modified
		CId         string `json:"cid"`         // value that has to be modified
		Bname       string `json:"bname"`       // value that has to be modified
		Bdesc       string `json:"bdesc"`       // value that has to be modified
		Bmodifiedby string `json:"bmodifiedby"` // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	filter := bson.D{{"bid", body.Bid}} // converting value to BSON type
	after := options.After              // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"scid", body.ScId}, {"cid", body.CId}, {"bname", body.Bname}, {"bdesc", body.Bdesc}, {"bmodifiedby", body.Bmodifiedby}}}}
	updateResult := brandCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
	w.WriteHeader(200)

}

// Update Brand Status

func UpdateBrandStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		BId     string `json:"bid"`     //value that has to be matched
		Bstatus bool   `json:"bstatus"` // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	filter := bson.D{{"bid", body.BId}} // converting value to BSON type
	after := options.After              // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"bstatus", body.Bstatus}}}}
	updateResult := categoryCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
	w.WriteHeader(200)

}

//Delete Brand

func DeleteBrand(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res, err := brandCollection.DeleteOne(context.TODO(), bson.D{{"bid", params}}, opts)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(400)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted
	w.WriteHeader(200)
}
