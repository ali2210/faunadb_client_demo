package wrapper

import (
	"log"

	"github.com/fauna/faunadb-go/v4/faunadb"
)

// Fauna Wrapper attributes
const DATABASE_NAME string = "equity_index"
const COLLECTION_NAME string = "equity_market"
const Index_Name string = "today"
const Document_Name string = "btc_today"

// Fauna Wrapper (Single document) interface. This plugin allow you to create a single document in faunadb.
type FaunaWrapper interface {
	WrapperConnection(api, endpoint string) *faunadb.FaunaClient
	FaunaClientObject(client *faunadb.FaunaClient) (faunadb.Value, error)
	SuperDB(client *faunadb.FaunaClient) (*faunadb.FaunaClient, error)

	WCollection(client *faunadb.FaunaClient) error
	WIndex(client *faunadb.FaunaClient) error
	WSDocs(client *faunadb.FaunaClient, ref ...string) error
	WSData(client *faunadb.FaunaClient, value ...Index) error
}

// Fauna Data store
type Index struct {
	Name      string  `json:"name"`
	Price     string  `json:"price"`
	MarketCap string  `json:"marketCap"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Symbol    string  `json:"symbol"`
}

// Receiver Function "New Wrapper" @ return Fauna Wrapper
func NewWrapper() FaunaWrapper { return &Index{} }

// Wrapper connection @ args api and endpoint ("string") ; return @ faunadb Client
// This function allow to established connection application and fauna database
func (i *Index) WrapperConnection(api, endpoint string) *faunadb.FaunaClient {
	return faunadb.NewFaunaClient(api, faunadb.Endpoint(endpoint))
}

// Fauna Client Object args  @ client ; return @ faunadb Value.
// This function valid that database name.
func (i *Index) FaunaClientObject(client *faunadb.FaunaClient) (faunadb.Value, error) {

	db, err := client.Query(faunadb.If(faunadb.Exists(faunadb.Database(DATABASE_NAME)), true, faunadb.CreateDatabase(faunadb.Obj{
		"name": DATABASE_NAME,
	})))
	if err != nil && db != faunadb.BooleanV(true) {
		return db, err
	}

	return db, nil

}

// SuperDB args client ; returns @ fauna client and error
// Super DB create allow to application server permission to create new documents inside database.
func (i *Index) SuperDB(client *faunadb.FaunaClient) (*faunadb.FaunaClient, error) {

	key, err := client.Query(faunadb.CreateKey(faunadb.Obj{
		"database": faunadb.Database(DATABASE_NAME),
		"role":     "server",
	}))
	if err != nil {
		return &faunadb.FaunaClient{}, err
	}

	const secret_key string = "secret"
	var valid_key string = " "

	err = key.At(faunadb.ObjKey(secret_key)).Get(&valid_key)
	if err != nil {
		return &faunadb.FaunaClient{}, err
	}

	return client.NewSessionClient(valid_key), nil
}

// WCollection @ args client ; returns @ error
// This function allow to create a collection inside the database.
func (i *Index) WCollection(client *faunadb.FaunaClient) error {

	value, err := client.Query(faunadb.If(faunadb.Exists(faunadb.Collection(COLLECTION_NAME)), true, faunadb.CreateCollection(faunadb.Obj{
		"name": COLLECTION_NAME,
	})))
	if err != nil {
		return err
	}

	if value != faunadb.BooleanV(true) {
		return i.WCollection(client)
	}

	return nil
}

// WIndex @ args client ; returns @ error
// This function allow to create index. In Index new documents will be created.
//  database ===* collection *===* index *===** documents **===** data

func (i *Index) WIndex(client *faunadb.FaunaClient) error {

	index_value, err := client.Query(faunadb.If(faunadb.Exists(faunadb.Index(Index_Name)), true, faunadb.CreateIndex(faunadb.Obj{
		"name":   Index_Name,
		"source": faunadb.Collection(COLLECTION_NAME),
		"unique": true,
	})))

	if err != nil {
		return err
	}

	if index_value != faunadb.BooleanV(true) {
		return i.WIndex(client)
	}

	return nil
}

// WSDoc @args client, reference []string ; returns @ error
func (i *Index) WSDocs(client *faunadb.FaunaClient, ref ...string) error {

	_, err := client.Query(faunadb.If(faunadb.Exists(
		faunadb.Documents(faunadb.Collection(COLLECTION_NAME))), true, faunadb.Map(
		faunadb.Paginate(faunadb.Documents(
			faunadb.Collection(COLLECTION_NAME)), faunadb.Size(1)),
		faunadb.Lambda(Document_Name+ref[0], faunadb.Get(faunadb.Var(Document_Name+ref[0]))))))
	if err != nil {
		return err
	}

	return nil
}

// WSData @ args client, value []Index ; returns @ error
func (i *Index) WSData(client *faunadb.FaunaClient, value ...Index) error {

	values := faunadb.Create(faunadb.Ref(faunadb.Collection(COLLECTION_NAME), 0), faunadb.Obj{
		"data": value[0],
	})

	equity := make([]faunadb.Expr, 1)
	equity = append(equity, values)

	_, err := client.Query(faunadb.Do(equity))
	if err != nil {
		log.Fatalln("Error creating index_value", err)
		return err
	}

	return nil
}

type MultiFaunaOperations interface {
	WMDocs(client *faunadb.FaunaClient, ref ...string) error
	WMData(client *faunadb.FaunaClient, num int64, value ...Index) error
	WMIndex(client *faunadb.FaunaClient, doc ...string) error
}

func NewMultiReader() MultiFaunaOperations { return &Index{} }

func (r *Index) WMDocs(client *faunadb.FaunaClient, ref ...string) error {

	_, err := client.Query(faunadb.If(faunadb.Exists(
		faunadb.Documents(faunadb.Collection(COLLECTION_NAME))), true, faunadb.Map(
		faunadb.Paginate(faunadb.Documents(
			faunadb.Collection(COLLECTION_NAME)), faunadb.Size(len(ref)-1)),
		faunadb.Lambda(Document_Name+ref[0], faunadb.Get(faunadb.Var(Document_Name+ref[0]))))))
	if err != nil {
		return err
	}

	return nil
}

func (r *Index) WMData(client *faunadb.FaunaClient, num int64, value ...Index) error {

	log.Println("Num :", num)
	values := faunadb.Create(faunadb.Ref(faunadb.Collection(COLLECTION_NAME), num), faunadb.Obj{
		"data": value[0],
	})

	// // equity = append(equity, values)
	_, err := client.Query(faunadb.Do(values))
	if err != nil {

		log.Fatalln("Error creating index_value:", err)
		return err
	}

	return nil

}

func (r *Index) WMIndex(client *faunadb.FaunaClient, doc ...string) error {

	var index_value faunadb.Value
	var err error

	index_value, err = client.Query(faunadb.If(faunadb.Exists(faunadb.Index(doc[0])), true, faunadb.CreateIndex(faunadb.Obj{
		"name":   doc[0],
		"source": faunadb.Collection(COLLECTION_NAME),
		"unique": true,
	})))

	if err != nil {
		return err
	}

	if index_value != faunadb.BooleanV(true) {
		return r.WMIndex(client, doc[0])
	}

	return nil
}
