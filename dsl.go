package es

// O Object.(abbr) JSON Object 类型别名
// 多用于Search().Source(es.O{})目的是为了提高代码可读性.
type O map[string]any

// A Array.(abbr) JSON Array 类型别名
// 多用于Search().Source(es.O{"...": es.A{"",""})目的是为了提高代码可读性.
type A []any
