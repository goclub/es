package es_test

import (
	"context"
	"github.com/goclub/error"
	es "github.com/goclub/es"
	xjson "github.com/goclub/json"
	es7 "github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/suite"
	"log"
	"reflect"
	"testing"
)

var example Example

type Example struct {
	Client *es7.Client
}

func init() {
	/*
		请提前准备数据
		PUT /bank
		{
			"mappings": {
			"properties": {
				"address": {"type": "text"}
			}
		}
		}
		https://download.elastic.co/demos/kibana/gettingstarted/accounts.zip
		curl -H "Content-Type: application/json" -XPOST "localhost:9200/bank/_bulk?pretty&refresh" --data-binary "@accounts.json"
	*/
	client, err := es7.NewClient(
		es7.SetURL(ExampleConnectURL),
		es7.SetSniff(false),
		es7.SetHealthcheck(false),
	)
	if err != nil {
		xerr.PrintStack(err)
	}
	example.Client = client
}

type ExampleAccount struct {
	AccountNumber int    `json:"account_number"`
	Balance       int    `json:"balance"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Age           int    `json:"age"`
	Gender        string `json:"gender"`
	Address       string `json:"address"`
	Employer      string `json:"employer"`
	Email         string `json:"email"`
	City          string `json:"city"`
	State         string `json:"state"`
}

func TestExample(t *testing.T) {
	suite.Run(t, new(TestExampleSuite))
}

type TestExampleSuite struct {
	suite.Suite
}

func (suite TestExampleSuite) TestExampleIndex() {
	ExampleIndex()
}

func (suite TestExampleSuite) TestExampleSearch() {
	ExampleSearch()
}
func (suite TestExampleSuite) TestExampleSearchFromSize() {
	ExampleSearchFromSize()
}
func (suite TestExampleSuite) TestExampleSearchMatch() {
	ExampleSearchMatch()
}
func (suite TestExampleSuite) TestExampleSearchMatchPhrase() {
	ExampleSearchMatchPhrase()
}
func (suite TestExampleSuite) TestExampleSearchBool() {
	ExampleSearchBool()
}
func (suite TestExampleSuite) TestExampleSearchBoolFilter() {
	ExampleSearchBoolFilter()
}

/*
POST /test/_doc
{
	"name": "nimo",
	"age" : 18
}
*/
func ExampleIndex() {
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	user := struct {
		Name string `json:"name"`
		Age  uint8  `json:"age"`
	}{Name: "nimo", Age: 18}
	indexResp, err := example.Client.Index().
		Index("test").
		BodyJson(user).
		Do(ctx)
	if err != nil {
		return
	}
	xjson.Print("indexResp", indexResp)
	return
}

/*
GET /bank/_search
{
  "query": { "match_all": {} },
  "sort": [
    { "account_number": "asc" }
  ],
  "size": 10
}
*/
func ExampleSearch() {
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	searchResult, err := example.Client.Search().
		Index("bank").
		Sort("account_number", true).
		Size(3).
		Do(ctx)
	if err != nil {
		return
	}
	xjson.Print("searchResult", searchResult)
	xjson.Print("TotalHits", searchResult.TotalHits())
	var itemType ExampleAccount
	var eachAccountList []ExampleAccount
	log.Print("------------------")
	log.Print("使用 Each 遍历,json解析失败的会赋空值")
	for _, item := range searchResult.Each(reflect.TypeOf(itemType)) {
		if account, ok := item.(ExampleAccount); ok {
			eachAccountList = append(eachAccountList, account)
		}
	}
	xjson.Print("eachAccountList", eachAccountList)
	log.Print("------------------")
	log.Print("使用 es.SearchSlice7 获取 slice, json 解释失败会返回错误( goclub/json )")
	slice7AccountList, err := es.SearchSlice7(searchResult, ExampleAccount{})
	if err != nil {
		if decodeErr, as := es.AsDecodeSearchResultError(err); as {
			log.Print(string(decodeErr.SearchHit))
			log.Printf("%+#v", decodeErr.Element)
		}
		return
	}
	xjson.Print("slice7AccountList", slice7AccountList)
	return
}

/*
GET /bank/_search
{
  "query": { "match_all": {} },
  "sort": [
    { "account_number": "asc" }
  ],
  "from": 10,
  "size": 10
}
*/
func ExampleSearchFromSize() {
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	searchResult, err := example.Client.Search().
		Index("bank").
		Sort("account_number", true).
		From(10).
		Size(10).
		Do(ctx)
	if err != nil {
		return
	}
	formSizeList, err := es.SearchSlice7(searchResult, ExampleAccount{})
	if err != nil {
		return
	}
	xjson.Print("formSizeList", formSizeList)
	return
}
func ExampleSearchMatch() {
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	searchResult, err := example.Client.Search().
		Index("bank").
		Query(es7.NewMatchQuery(
			"address", "mill lane",
		)).
		Size(10).
		Do(ctx)
	if err != nil {
		return
	}
	xjson.Print("searchResult.TotalHits()", searchResult.TotalHits())
	matchList, err := es.SearchSlice7(searchResult, ExampleAccount{})
	if err != nil {
		return
	}
	xjson.Print("matchList", matchList)
	return
}

/*
GET /bank/_search
{
  "query": {
    "bool": {
      "must": [
        { "match": { "age": "40" } }
      ],
      "must_not": [
        { "match": { "state": "ID" } }
      ]
    }
  }
}
*/
func ExampleSearchMatchPhrase() {
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	searchResult, err := example.Client.Search().
		Index("bank").
		Query(es7.NewMatchPhraseQuery(
			"address", "mill lane",
		)).
		Size(10).
		Do(ctx)
	if err != nil {
		return
	}
	xjson.Print("searchResult.TotalHits()", searchResult.TotalHits())
	matchList, err := es.SearchSlice7(searchResult, ExampleAccount{})
	if err != nil {
		return
	}
	xjson.Print("matchList", matchList)
	return
}

/*
GET /bank/_search
{
  "query": {
    "bool": {
      "must": [
        { "match": { "age": "40" } }
      ],
      "must_not": [
        { "match": { "state": "ID" } }
      ]
    }
  }
}
*/
func ExampleSearchBool() {
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	searchResult, err := example.Client.Search().
		Index("bank").
		Query(es7.NewBoolQuery().Must(
			es7.NewMatchQuery("age", "40"),
		).MustNot(
			es7.NewMatchQuery("state", "ID"),
		),
		).
		Do(ctx)
	if err != nil {
		return
	}
	xjson.Print("searchResult.TotalHits()", searchResult.TotalHits())
	matchBoolList, err := es.SearchSlice7(searchResult, ExampleAccount{})
	if err != nil {
		return
	}
	xjson.Print("matchBoolList", matchBoolList)
	return
}

/*
GET /bank/_search
{
  "query": {
    "bool": {
      "must": { "match_all": {} },
      "filter": {
        "range": {
          "balance": {
            "gte": 20000,
            "lte": 30000
          }
        }
      }
    }
  }
}
*/
func ExampleSearchBoolFilter() {
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	searchResult, err := example.Client.Search().
		Index("bank").
		Query(es7.NewBoolQuery().Must(
			es7.NewMatchAllQuery(),
		).Filter(
			es7.NewRangeQuery("balance").Gte(20000).Lte(30000),
		),
		).
		Do(ctx)
	if err != nil {
		return
	}
	xjson.Print("searchResult.TotalHits()", searchResult.TotalHits())
	matchBoolFilterList, err := es.SearchSlice7(searchResult, ExampleAccount{})
	if err != nil {
		return
	}
	xjson.Print("matchBoolFilterList", matchBoolFilterList)
	return
}
