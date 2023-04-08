package model

type Search struct {
	Query string `form:"q"`
}

type Suggest struct {
	Query string `form:"q"`
}
