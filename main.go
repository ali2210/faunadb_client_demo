package main

import (
	"log"

	wrapper "github.com/ali2210/faunadb_client_demo/fauna_wrapper"
	"github.com/fauna/faunadb-go/v4/faunadb"
)

// fauna credentials

const API_KEY string = "fnAEulZtxTACTKkUAO11O4w-K6vxQWCk3bsmqLLA"
const CONNECTION_ADDR string = "https://db.fauna.com:443"

func main() {

	wrapper_Client := wrapper.NewWrapper()

	_to_Connect := wrapper_Client.WrapperConnection(API_KEY, CONNECTION_ADDR)

	value, err := wrapper_Client.FaunaClientObject(_to_Connect)

	if err != nil && value != faunadb.BooleanV(true) {
		log.Fatalln(" Error :", err.Error())
		return
	}

	client, err := wrapper_Client.SuperDB(_to_Connect)
	if err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = wrapper_Client.WCollection(client); err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = wrapper_Client.WIndex(client); err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = wrapper_Client.WSDocs(client, "security_exchange"); err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = wrapper_Client.WSData(client, wrapper.Index{
		Name:      "Bitcoin",
		Price:     "21,377.30 USD",
		MarketCap: "$402.35B",
		High:      21394.11,
		Low:       209078,
		Symbol:    "Btc",
	}); err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

	// Multitple documents created , however documents hold replicate data
	multi_wrapper := wrapper.NewMultiReader()

	comodity := []string{"Gold", "Bitcoin", "Silver"}
	securities := []wrapper.Index{
		{
			Name:      "Gold",
			Price:     "1,731.70 USD",
			MarketCap: "$11.617 T",
			High:      2065.89,
			Low:       1681.43,
			Symbol:    "GOLD",
		},
		{
			Name:      "Bitcoin",
			Price:     "21,377.30 USD",
			MarketCap: "$402.35B",
			High:      21394.11,
			Low:       209078,
			Symbol:    "Btc",
		},
		{
			Name:      "Silver",
			Price:     "614.83 USD",
			High:      665.35,
			Low:       610.00,
			Symbol:    "SILVER",
			MarketCap: "$1.069 T",
		},
	}

	if err = multi_wrapper.WMIndex(client, comodity[0]); err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = multi_wrapper.WMDocs(client, comodity[0]); err != nil {

		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = multi_wrapper.WMData(client, 0, securities[0]); err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = multi_wrapper.WMIndex(client, comodity[1]); err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = multi_wrapper.WMDocs(client, comodity[1]); err != nil {

		log.Fatalln(" Error :", err.Error())
		return
	}

	if err = multi_wrapper.WMData(client, 1, securities[1]); err != nil {
		log.Fatalln(" Error :", err.Error())
		return
	}

}
