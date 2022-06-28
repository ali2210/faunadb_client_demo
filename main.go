package main

import (
	"log"

	"github.com/fauna/faunadb-go/v4/faunadb"
)

const API_KEY string = "fnAEqFiXToACS-TCQs8568lPVukXhBNKdDuWMYVX"
const CONNECTION_ADDR string = "https://db.fauna.com:443"
const DATABASE_NAME string = "equity_index"
const COLLECTION_NAME string = "equity_market"
const Index_Name string = "today"
const Document_Name string = "btc_today"

type Index struct {
	Name      string  `json:"name"`
	Price     string  `json:"price"`
	MarketCap string  `json:"marketCap"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
}

func main() {

	client := faunadb.NewFaunaClient(API_KEY, faunadb.Endpoint(CONNECTION_ADDR))

	ConnectFauna(client)
	newFClient, err := GetFaunadbDatabaseObject(client)
	if err != nil {
		log.Fatalln(" Database Client error: ", err)
		return
	}

	if ok := FaunaCollection(newFClient); !ok {
		log.Fatalln("Collection is not created:", ok)
		return
	}

	if err := FaunaIndex(newFClient); err != nil {
		log.Fatalln("Index is not created:", err)
		return
	}

	dip := make([]Index, 3)
	dip = append(dip, Index{Name: "BTC",
		Price:     "21,377.30 USD",
		MarketCap: "$402.35B",
		High:      21394.11,
		Low:       209078,
	})

	if err = FaunaHoldState(newFClient, dip); err != nil {
		log.Fatalln("Index created by fail to hold index value", err)
		return
	}

	if err = FaunaDoc(newFClient); err != nil {
		log.Fatalln("Documents is not created:", err)
		return
	}

	log.Println("Operation succeeded....")

}

func ConnectFauna(client *faunadb.FaunaClient) (*faunadb.Value, error) {

	db, err := client.Query(faunadb.If(faunadb.Exists(faunadb.Database(DATABASE_NAME)), true, faunadb.CreateDatabase(faunadb.Obj{
		"name": DATABASE_NAME,
	})))
	if err != nil {
		return &db, err
	}

	if db != faunadb.BooleanV(true) {
		return &db, err
	}

	return &db, nil
}

func GetFaunadbDatabaseObject(client *faunadb.FaunaClient) (*faunadb.FaunaClient, error) {

	key, err := client.Query(faunadb.CreateKey(faunadb.Obj{
		"database": faunadb.Database(DATABASE_NAME),
		"role":     "server",
	}))
	if err != nil {
		return &faunadb.FaunaClient{}, err
	}

	const passcode string = "secret"
	var valid_key string = ""

	err = key.At(faunadb.ObjKey(passcode)).Get(&valid_key)
	if err != nil {
		return &faunadb.FaunaClient{}, err
	}

	return client.NewSessionClient(valid_key), nil
}

func FaunaCollection(client *faunadb.FaunaClient) bool {

	value, err := client.Query(faunadb.If(faunadb.Exists(faunadb.Collection(COLLECTION_NAME)), true, faunadb.CreateCollection(faunadb.Obj{
		"name": COLLECTION_NAME,
	})))
	if err != nil {
		return false
	}

	if value != faunadb.BooleanV(true) {
		return false
	}

	return true

}

func FaunaIndex(client *faunadb.FaunaClient) (err error) {

	index_value, err := client.Query(faunadb.If(faunadb.Exists(faunadb.Index(Index_Name)), true, faunadb.CreateIndex(faunadb.Obj{
		"name":       Index_Name,
		"source":     faunadb.Collection(COLLECTION_NAME),
		"serialized": true,
		"unique":     true,
		"terms":      []string{"data", "name"},
	})))

	if err != nil {

		return err
	}

	if index_value != faunadb.BooleanV(true) {

		return err
	}

	return nil
}

func FaunaDoc(client *faunadb.FaunaClient) error {

	_, err := client.Query(faunadb.If(faunadb.Exists(
		faunadb.Documents(faunadb.Collection(COLLECTION_NAME))), true, faunadb.Map(
		faunadb.Paginate(faunadb.Documents(
			faunadb.Collection(COLLECTION_NAME)), faunadb.Size(1)),
		faunadb.Lambda(Document_Name+"Ref", faunadb.Get(faunadb.Var(Document_Name+"Ref"))))))
	if err != nil {
		return err
	}

	return nil
}

func FaunaHoldState(client *faunadb.FaunaClient, value []Index) error {

	values := faunadb.Create(faunadb.Ref(faunadb.Collection(COLLECTION_NAME), 0), faunadb.Obj{
		"data": value[len(value)-1],
	})

	equity := make([]faunadb.Expr, len(value))
	equity = append(equity, values)

	_, err := client.Query(faunadb.Do(equity))
	if err != nil {
		log.Fatalln("Error creating index_value", err)
		return err
	}

	return nil
}
