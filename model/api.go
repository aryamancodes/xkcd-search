package model

type Search struct {
	Query       string `form:"q" binding:"required"`
	Autocorrect bool   `form:"autocorrect"`
	Start       int    `form:"start"`
}

type Suggest struct {
	Query string `form:"q" binding:"required"`
}
