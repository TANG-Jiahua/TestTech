package Model


type User struct {
	Id string
	Password string
	IsActive bool
	Balance string
	Age int
	Name string
	Gender string
	Company string
	Email string
	Phone string
	Address string
	About string
	Registered string
	Latitude float64
	Longitude float64
	Tags []string
	Friends []Friend
	Data string

}

type Friend struct {
	Id int
	Name string
}