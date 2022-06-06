package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct for storing data
type brand struct {
	BId         string `json:"bid" validate:"required,alphanum,min=4,max=10"`
	ScId        string `json:"scid" validate:"required,alphanum,min=4,max=10"`
	CId         string `json:"cid" validate:"required,alphanum,min=4,max=10"`
	Bname       string `json:"bname" validate:"required,min=3,max=20"`
	Bdesc       string `json:"bdesc" validate:"required,min=5,max=100"`
	Bcreatedby  string `json:"bcreatedby" validate:"required,min=3,max=20"`
	Bmodifiedby string `json:"bmodifiedby" validate:"required,min=3,max=20"`
	Bstatus     bool   `json:"bstatus"`
}

var brandCollection = db().Database("ProductApp").Collection("Brand") // get collection "users" from db() which returns *mongo.Client

// Create Brand

func CreateBrand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // for adding Content-type
	validate := validator.New()
	var brand brand
	err := json.NewDecoder(r.Body).Decode(&brand) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(400)
	}
	errv := validate.Struct(brand) //create struct validation
	if errv != nil {
		fmt.Println(errv)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Validation error")
	} else {
		ccount, errc := categoryCollection.CountDocuments(context.TODO(), bson.D{{"cid", brand.CId}})
		scount, errs := subCategoryCollection.CountDocuments(context.TODO(), bson.D{{"scid", brand.ScId}})
		bcount, errb := brandCollection.CountDocuments(context.TODO(), bson.D{{"bid", brand.BId}})

		if bcount == 0 && scount == 1 && ccount == 0 {
			insertResult, err := brandCollection.InsertOne(context.TODO(), brand)
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(500)
			} else {
				fmt.Println("Inserted a single document: ", insertResult)
				json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document
				w.WriteHeader(200)
			}
		} else if scount == 0 {
			fmt.Print(errs)
			json.NewEncoder(w).Encode("Invalid Subcategory")
		} else if ccount == 0 {
			fmt.Print(errc)
			json.NewEncoder(w).Encode("Invalid Category")
		} else {
			fmt.Print(errb)
			json.NewEncoder(w).Encode("Invalid Brand")
		}

	}
}

// Get Brand

func GetBrand(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(params)
	// var body brand
	// e := json.NewDecoder(r.Body).Decode(&body)
	// if e != nil {
	// 	fmt.Print(e)
	// }
	if is_alphanumeric {
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
	} else {
		json.NewEncoder(w).Encode("Please enter the correct Brand ID!")
		w.WriteHeader(204)
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
		BId         string `json:"bid" validate:"required,alphanum,min=4,max=10"`  // value that has to be matched
		ScId        string `json:"scid" validate:"required,alphanum,min=4,max=10"` // value that has to be modified
		CId         string `json:"cid" validate:"required,alphanum,min=4,max=10"`  // value that has to be modified
		Bname       string `json:"bname" validate:"required,min=3,max=20"`         // value that has to be modified
		Bdesc       string `json:"bdesc" validate:"required,min=5,max=100"`        // value that has to be modified
		Bmodifiedby string `json:"bmodifiedby" validate:"required,min=3,max=20"`   // value that has to be modified
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
	validate := validator.New()

	type updateBody struct {
		BId     string `json:"bid" validate:"required,alphanum,min=4,max=10"` //value that has to be matched
		Bstatus bool   `json:"bstatus"`                                       // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
		w.WriteHeader(400)
	}
	errv := validate.Struct(body) // update status struct validation
	if errv != nil {
		fmt.Println(errv)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Validation error")
	} else {
		filter := bson.D{{"bid", body.BId}} // converting value to BSON type
		after := options.After              // for returning updated document
		returnOpt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		update := bson.D{{"$set", bson.D{{"bstatus", body.Bstatus}}}}
		updateResult := brandCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)
		var result primitive.M
		_ = updateResult.Decode(&result)
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(200)
	}

}

//Delete Brand

func DeleteBrand(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(params)
	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	if is_alphanumeric {
		opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
		res, err := brandCollection.DeleteOne(context.TODO(), bson.D{{"bid", params}}, opts)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(400)
		}
		fmt.Printf("deleted %v documents\n", res.DeletedCount)
		json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted
		w.WriteHeader(200)
	} else {
		json.NewEncoder(w).Encode("Please enter the correct Brand ID!")
		w.WriteHeader(204)
	}
}
