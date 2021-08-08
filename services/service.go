package services


type Service interface {
	Run() (*interface{}, error)
}
