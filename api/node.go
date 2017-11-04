package api

type Node {
    id int
    name string
    parentId int `json:"parent_id"`
}
