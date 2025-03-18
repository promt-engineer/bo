package repositories

type Filters struct {
	Where map[string]interface{}
	Not   map[string]interface{}
}
