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

type subprod struct {
	CategoryId    string `json:"categoryid" validate:"required,alphanum,min=4,max=10"`
	SubCategoryId string `json:"subcategoryid" validate:"required,alphanum,min=4,max=10"`
	BrandId       string `json:"brandid" validate:"required,alphanum,min=4,max=10"`
	VarientId     string `json:"varientid" validate:"required,alphanum,min=4,max=10"`
}
type product struct {
	PId         string  `json:"pid" validate:"required,alphanum,min=4,max=10"`
	Pname       string  `json:"pname" validate:"required,min=3,max=20"`
	Pdesc       string  `json:"pdesc" validate:"required,min=5,max=100"`
	Pqty        int     `json:"pqty" validate:"required, numeric"`
	Pmrp        float32 `json:"pmrp" validate:"required, numeric"`
	Pprice      float32 `json:"pprice" validate:"required, numeric"`
	Pcreatedby  string  `json:"pcreatedby" validate:"required,min=3,max=20"`
	Pmodifiedby string  `json:"pmodifiedby" validate:"required,min=3,max=20"`
	Pstatus     bool    `json:"pstatus"`
	SubProd     subprod `json:"subprod" validate:"required"`
}

var productCollection = db().Database("ProductApp").Collection("Product") // get collection "users" from db() which returns *mongo.Client

// Create Product

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type
	validate := validator.New()
	var prod product
	err := json.NewDecoder(r.Body).Decode(&prod) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(400)
	}
	errv := validate.Struct(prod) //create struct validation
	if errv != nil {
		fmt.Println(errv)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Validation error")
	} else {
		pcount, errp := productCollection.CountDocuments(context.TODO(), bson.D{{"pid", prod.PId}})
		ccount, errc := categoryCollection.CountDocuments(context.TODO(), bson.D{{"categoryid", prod.SubProd.CategoryId}})
		scount, errs := subCategoryCollection.CountDocuments(context.TODO(), bson.D{{"subcategoryid", prod.SubProd.SubCategoryId}})
		bcount, errb := brandCollection.CountDocuments(context.TODO(), bson.D{{"brandid", prod.SubProd.BrandId}})
		vcount, errv := varientCollection.CountDocuments(context.TODO(), bson.D{{"varientid", prod.SubProd.VarientId}})

		if pcount == 0 && ccount == 1 && scount == 1 && bcount == 1 && vcount == 1 {
			insertResult, err := productCollection.InsertOne(context.TODO(), prod)
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(500)
			} else {
				fmt.Println("Inserted a single document: ", insertResult)
				json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document
				w.WriteHeader(200)
			}
		} else if ccount == 0 {
			fmt.Print(errc)
			json.NewEncoder(w).Encode("Invalid Category")
		} else if scount == 0 {
			fmt.Print(errs)
			json.NewEncoder(w).Encode("Invalid Subcategory")
		} else if bcount == 0 {
			fmt.Print(errb)
			json.NewEncoder(w).Encode("Invalid Brand")
		} else if vcount == 0 {
			fmt.Print(errv)
			json.NewEncoder(w).Encode("Invalid Varient")
		} else {
			fmt.Print(errp)
			json.NewEncoder(w).Encode("Invalid Product")
		}
	}

}

// Get Product

func GetProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(params)
	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }

	// var body product
	// e := json.NewDecoder(r.Body).Decode(&body)
	// if e != nil {

	// 	fmt.Print(e)
	// }
	if is_alphanumeric {
		var result primitive.M //  an unordered representation of a BSON document which is a Map
		filter := bson.M{"pid": params}
		err := productCollection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("No data found!")
		} else {
			json.NewEncoder(w).Encode(result) // returns a Map containing document
			w.WriteHeader(200)
		}
	} else {
		json.NewEncoder(w).Encode("Please enter the correct Product ID!")
		w.WriteHeader(204)
	}

}

// Get All Product

func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                      //slice for multiple documents
	cur, err := productCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
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

//Update Product

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	validate := validator.New()
	type subprod struct {
		CategoryId    string `json:"categoryid" validate:"required,alphanum,min=4,max=10"`
		SubCategoryId string `json:"subcategoryid" validate:"required,alphanum,min=4,max=10"`
		BrandId       string `json:"brandid" validate:"required,alphanum,min=4,max=10"`
		VarientId     string `json:"varientid" validate:"required,alphanum,min=4,max=10"`
	}

	type updateBody struct {
		PId         string  `json:"pid" validate:"required,alphanum,min=4,max=10"` //value that has to be matched
		Pname       string  `json:"pname" validate:"required,min=3,max=20"`        // value that has to be modified
		Pdesc       string  `json:"pdesc" validate:"required,min=5,max=100"`       // value that has to be modified
		Pqty        int     `json:"pqty" validate:"required, numeric"`             // value that has to be modified
		Pmrp        float32 `json:"pmrp" validate:"required, numeric"`             // value that has to be modified
		Pprice      float32 `json:"pprice" validate:"required, numeric"`           // value that has to be modified
		Pmodifiedby string  `json:"pmodifiedby" validate:"required,min=3,max=20"`  // value that has to be modified
		SubProd     subprod `json:"subprod" validate:"required"`                   // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	errv := validate.Struct(body) // update struct validation
	if errv != nil {
		fmt.Println(errv)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Validation error")
	} else {
		filter := bson.D{{"pid", body.PId}} // converting value to BSON type
		after := options.After              // for returning updated document
		returnOpt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		update := bson.D{{"$set", bson.D{{"pname", body.Pname}, {"pdesc", body.Pdesc}, {"pqty", body.Pqty}, {"pmrp", body.Pmrp}, {"pprice", body.Pprice}, {"pmodifiedby", body.Pmodifiedby}, {"subprod", bson.D{{"categoryid", body.SubProd.CategoryId}, {"varientid", body.SubProd.VarientId}, {"subcategoryid", body.SubProd.SubCategoryId}, {"brandid", body.SubProd.BrandId}}}}}}
		updateResult := productCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)
		// update1 := bson.D{{"$set", bson.D{{"age", body.Age}}}}
		// updateResult1 := productCollection.FindOneAndUpdate(context.TODO(), filter, update1, &returnOpt)
		var result primitive.M
		_ = updateResult.Decode(&result)

		json.NewEncoder(w).Encode(result)
		w.WriteHeader(200)
	}

}

// Update Product Status

func UpdateProductStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	validate := validator.New()

	type updateBody struct {
		PId     string `json:"pid" validate:"required,alphanum,min=4,max=10"` //value that has to be matched
		Pstatus bool   `json:"pstatus"`                                       // value that has to be modified
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
		filter := bson.D{{"pid", body.PId}} // converting value to BSON type
		after := options.After              // for returning updated document
		returnOpt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		update := bson.D{{"$set", bson.D{{"pstatus", body.Pstatus}}}}
		updateResult := productCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)
		var result primitive.M
		_ = updateResult.Decode(&result)
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(200)
	}

}

//Delete Product

func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(params)
	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	if is_alphanumeric {
		opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
		res, err := productCollection.DeleteOne(context.TODO(), bson.D{{"pid", params}}, opts)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(400)
		} else {
			fmt.Printf("deleted %v documents\n", res.DeletedCount)
			json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted
			w.WriteHeader(200)
		}
	} else {
		json.NewEncoder(w).Encode("Please enter the correct Product ID!")
		w.WriteHeader(204)
	}

}
