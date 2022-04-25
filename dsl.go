package es

// M Map.(abbr) JSON Object 类型别名
// 多用于Search().Source(es.M{})目的是为了提高代码可读性.
type M map[string]any

// A Array.(abbr) JSON Array 类型别名
// 多用于Search().Source(es.M{"...": es.A{"",""})目的是为了提高代码可读性.
type A []any
