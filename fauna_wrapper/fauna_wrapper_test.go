package wrapper

import (
	"testing"

	"github.com/fauna/faunadb-go/v4/faunadb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Wrapper(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "Fauna DB wrapper test")
}

var _ = Describe("Fauna Wrapper  Test suit", func() {
	var wrapperObject FaunaWrapper
	var connect, client *faunadb.FaunaClient
	var _dbinstance faunadb.Value
	var err error

	const api string = "fnAEulZtxTACTKkUAO11O4w-K6vxQWCk3bsmqLLA"
	const addr string = "https://db.fauna.com:443"

	BeforeEach(func() {

		When("Fauna Database wrapper connection established", func() {
			wrapperObject = NewWrapper()
			connect = wrapperObject.WrapperConnection(api, addr)
			Expect(connect).To(Succeed())
		})

		AfterEach(func() {
			When("After established connection", func() {
				Context(" Secondary fauna database with admin role ", func() {
					It(" Database exists", func() {

						_dbinstance, err = wrapperObject.FaunaClientObject(client)
						Expect(err).Error().ShouldNot(BeNil())

						err = _dbinstance.Get("equity_index")
						Expect(err).Error().ShouldNot(BeNil())

						client, err = wrapperObject.SuperDB(connect)
						Expect(err).Error().ShouldNot(BeNil())
						Expect(client).Should(BeEmpty())

					})
				})
				Context("Fauna wrapper new collection ", func() {
					It("should create collection", func() {

						err = wrapperObject.WCollection(client)
						Expect(err).ShouldNot(BeNil())

					})
				})
				Context("Fauna wrapper new document", func() {
					It("should create new document", func() {

						err = wrapperObject.WIndex(client)
						Expect(err).ShouldNot(BeNil())
					})
				})
				Context("Fauna create document", func() {

					It("should new document created", func() {

						err = wrapperObject.WSDocs(client, "*")
						Expect(err).ShouldNot(BeNil())
					})
				})
				Context("Fauna insert document details", func() {
					It("should insert data", func() {

						err = wrapperObject.WSData(client, Index{
							Name:      "Bitcoin",
							Price:     "21,377.30 USD",
							MarketCap: "$402.35B",
							High:      21394.11,
							Low:       209078,
							Symbol:    "Btc",
						})

						Expect(err).ShouldNot(BeNil())
					})
				})
			})

		})
	})
})
