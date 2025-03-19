package services

type City struct {
	ID         int
	Name       string
	Timezone   string
	Country    string
	TimeOffset int
}

type RegistryService struct{}

func NewRegistryService() *RegistryService {
	return &RegistryService{}
}

func (r *RegistryService) GetCityByID(cityID int) (City, bool) {
	cities := map[int]City{
		1: {ID: 1, Name: "Dhaka", Timezone: "Asia/Dhaka", Country: "Bangladesh", TimeOffset: 21600},
		2: {ID: 2, Name: "Chittagong", Timezone: "Asia/Dhaka", Country: "Bangladesh", TimeOffset: 21600},
		3: {ID: 3, Name: "Sylhet", Timezone: "Asia/Dhaka", Country: "Bangladesh", TimeOffset: 21600},
		4: {ID: 4, Name: "Kathmandu", Timezone: "Asia/Kathmandu", Country: "Nepal", TimeOffset: 20700},
		5: {ID: 5, Name: "Khulna", Timezone: "Asia/Dhaka", Country: "Bangladesh", TimeOffset: 21600},
		6: {ID: 6, Name: "Chitwan", Timezone: "Asia/Kathmandu", Country: "Nepal", TimeOffset: 20700},
		7: {ID: 7, Name: "Pokhara", Timezone: "Asia/Kathmandu", Country: "Nepal", TimeOffset: 20700},
	}

	city, exists := cities[cityID]
	return city, exists
}
