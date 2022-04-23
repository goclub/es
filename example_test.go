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

var example es.Example

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

func init() {
	ExampleExample_NewClient()
}

func ExampleExample_NewClient() {
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

func TestExample(t *testing.T) {
	suite.Run(t, new(TestExampleSuite))
}

type TestExampleSuite struct {
	suite.Suite
}

func (suite TestExampleSuite) TestExampleIndex() {
	ExampleExample_Index()
}

func (suite TestExampleSuite) TestExampleSearch() {
	ExampleExample_Search()
}
func (suite TestExampleSuite) TestExampleSearchFromSize() {
	ExampleExample_SearchFromSize()
}
func (suite TestExampleSuite) TestExampleSearchMatch() {
	ExampleExample_SearchMatch()
}
func (suite TestExampleSuite) TestExampleSearchMatchPhrase() {
	ExampleExample_SearchMatchPhrase()
}
func (suite TestExampleSuite) TestExampleSearchBool() {
	ExampleExample_SearchBool()
}
func (suite TestExampleSuite) TestExampleSearchBoolFilter() {
	ExampleExample_SearchBoolFilter()
}
func (suite TestExampleSuite) TestExampleSearchGroupBy() {
	ExampleExample_SearchGroupBy()
}

// ExampleExample_Index 索引文档
func ExampleExample_Index() {
	/*
		POST /test/_doc
		{
			"name": "nimo",
			"age" : 18
		}
	*/
	var err error
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

// ExampleExample_Search 搜索全部数据
func ExampleExample_Search() {
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

// ExampleExample_SearchFromSize 分页搜索
func ExampleExample_SearchFromSize() {
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

// ExampleExample_SearchMatch 搜索匹配分词
// 确定一下 bank 是否配置了mappings address,如果 mappings 类型是 keyword 则无法分词搜索
func ExampleExample_SearchMatch() {
	/*
		GET /bank/_search
		{
		  "query": { "match": { "address": "mill lane" } }
		}
	*/
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

// ExampleExample_SearchMatchPhrase 搜索匹配短语
func ExampleExample_SearchMatchPhrase() {
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

// ExampleExample_SearchBool 搜索 bool 条件
func ExampleExample_SearchBool() {
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

// ExampleExample_SearchBool 搜索 bool 条件
func ExampleExample_SearchBoolFilter() {
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
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	searchResult, err := example.Client.Search().
		Index("bank").
		Source(es.O{
			"query": es.O{
				"bool": es.O{
					"must": es.O{"match_all": es.O{}},
					"filter": es.O{
						"range": es.O{
							"balance": es.O{
								"gte": 20000,
								"lte": 30000,
							},
						},
					},
				},
			},
		}).
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

// ExampleExample_SearchGroupBy 搜索-聚合分析-分组
func ExampleExample_SearchGroupBy() {
	/*
		GET /bank/_search
		{
		  "size": 0,
		  "aggs": {
		    "group_by_state": {
		      "terms": {
		        "field": "state"
		      }
		    }
		  }
		}
	*/
	var err error
	/* only test use */ defer func() {
		if err != nil {
			xerr.PrintStack(err)
		}
	}()
	ctx := context.Background()
	result, err := example.Client.Search().Source(es.O{
		"size": 0,
		"aggs": es.O{
			"group_by_state": es.O{
				"terms": es.O{
					"field": "state",
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		return
	}
	xjson.PrintIndent("group_by_state", result.Aggregations["group_by_state"])
	return
}
