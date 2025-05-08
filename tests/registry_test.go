package tests_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"service-availability/internal/models"
	"service-availability/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestFetchCitySettings_Success(t *testing.T) {
	// Mock response data
	mockResponse := models.CitySettings{
		CityID: 3,
			FoodOpenHours: []string{
				"0800,2359,08:00AM-23:59PM",
        		"0800,2359,08:00AM-23:59PM",
        		"0800,2359,08:00AM-23:59PM",
        		"0800,2359,08:00AM-23:59PM",
        		"0800,2359,08:00AM-23:59PM",
        		"0800,2359,08:00AM-23:59PM",
        		"0800,2359,08:00AM-23:59PM",
			},
			FoodHexagons:[]string{
				"893ce36153bffff", "893ce36046bffff", "893ce360237ffff", "893ce36044bffff", "893ce361527ffff", "893ce36334fffff",
				"883ce36333fffff", "893ce36000bffff", "893ce363553ffff", "893ce3614afffff", "893ce3630cbffff", "883ce36335fffff",
				"893ce360233ffff", "893ce361537ffff", "893ce36323bffff", "893ce363113ffff", "893ce36001bffff", "893ce3614a7ffff",
				"893ce360257ffff", "893ce361533ffff", "893ce363543ffff", "883ce36307fffff", "893ce3631d3ffff", "883ce36303fffff",
				"893ce3615d3ffff", "893ce360697ffff", "883ce36323fffff", "893ce363167ffff", "893ce360687ffff", "893ce360207ffff",
				"893ce36354fffff", "893ce3615abffff", "883ce36069fffff", "883ce36067fffff", "893ce36006fffff", "893ce3602dbffff",
				"893ce3614b3ffff", "893ce3602afffff", "893ce363217ffff", "883ce36007fffff", "893ce363093ffff", "893ce36355bffff",
				"893ce3600d3ffff", "893ce3606b7ffff", "883ce3606dfffff", "883ce36021fffff", "883ce36005fffff", "893ce3631d7ffff",
				"893ce361487ffff", "893ce3602cbffff", "883ce36315fffff", "893ce363213ffff", "893ce363573ffff", "893ce3632abffff",
				"893ce3614bbffff", "893ce36312bffff", "893ce363563ffff", "873ce3633ffffff", "893ce36025bffff", "883ce3606bfffff",
				"893ce36310fffff", "893ce36357bffff", "893ce36328fffff", "883ce36311fffff", "893ce3600cfffff", "893ce36020bffff",
				"893ce36319bffff", "883ce3633dfffff", "893ce36313bffff", "883ce3602bfffff", "893ce3600c7ffff", "893ce3636a7ffff",
				"893ce3603d3ffff", "893ce36342bffff", "893ce36309bffff", "893ce3615cbffff", "893ce36039bffff", "893ce360077ffff",
				"893ce363183ffff", "893ce363087ffff", "893ce363117ffff", "893ce363637ffff", "893ce36354bffff", "893ce360247ffff",
				"893ce363203ffff", "893ce36362fffff", "893ce360067ffff", "893ce3630b3ffff", "893ce36346bffff", "893ce3615dbffff",
				"893ce3630d3ffff", "893ce36322bffff", "883ce36337fffff", "893ce363107ffff", "893ce363177ffff", "873ce3606ffffff",
				"893ce36318bffff", "893ce3630a7ffff", "883ce36305fffff", "883ce3615dfffff", "883ce36065fffff", "893ce36336bffff",
				"893ce36305bffff", "893ce36020fffff", "893ce360253ffff", "893ce361507ffff", "893ce360273ffff", "893ce36044fffff",
				"883ce36029fffff", "883ce36151fffff", "893ce36322fffff", "893ce3602a7ffff", "893ce36023bffff", "883ce36063fffff",
				"883ce3615bfffff", "893ce3630dbffff", "893ce363097ffff", "883ce3631dfffff", "893ce360073ffff", "893ce360693ffff",
				"893ce363283ffff", "883ce36159fffff", "893ce36150fffff", "893ce361523ffff", "893ce36337bffff", "893ce3614b7ffff",
				"893ce361497ffff", "893ce360277ffff", "893ce3632b3ffff", "883ce36339fffff", "893ce363277ffff", "893ce363567ffff",
				"893ce3602cfffff", "893ce36304bffff", "893ce363287ffff", "893ce3606a7ffff", "893ce36000fffff", "893ce36320fffff",
				"893ce363627ffff", "893ce363547ffff", "893ce3602abffff", "873ce3630ffffff", "893ce360243ffff", "893ce3630b7ffff",
				"893ce36335bffff", "893ce360013ffff", "893ce36356fffff", "893ce36356bffff", "893ce36045bffff", "893ce3630d7ffff",
				"883ce3600dfffff", "883ce3630bfffff", "893ce36320bffff", "883ce3633bfffff", "893ce361483ffff", "893ce36318fffff",
				"893ce36346fffff", "893ce3615afffff", "893ce360393ffff", "893ce363123ffff", "893ce363053ffff", "893ce3631c7ffff",
				"893ce360003ffff", "883ce3632bfffff", "893ce36312fffff", "893ce363233ffff", "893ce363467ffff", "893ce3606b3ffff",
				"893ce363193ffff", "893ce363297ffff", "893ce361493ffff", "893ce36021bffff", "893ce360057ffff", "893ce3614a3ffff",
				"893ce363207ffff", "893ce363463ffff", "893ce3603dbffff", "883ce36317fffff", "893ce361437ffff", "893ce36334bffff",
				"893ce36006bffff", "883ce3602dfffff", "893ce3632bbffff", "893ce36342fffff", "893ce360047ffff",
			},
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/api/v1/settings/cities/3")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Override the URL in the actual service for test purposes
	originalBaseURL := services.RegistryBaseURL
	services.RegistryBaseURL = server.URL
	defer func() { services.RegistryBaseURL = originalBaseURL }()

	// Run the test
	resp, err := services.FetchCitySettings(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, resp.CityID)
	assert.Equal(t, "0800,2359,08:00AM-23:59PM", resp.FoodOpenHours[0])
}
