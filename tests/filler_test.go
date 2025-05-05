package tests

import (
	"service-availability/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	h3 "github.com/uber/h3-go/v4"
)

func TestFillPolygonWithHexes(t *testing.T) {
	// Test cases for different polygon boundaries
	tests := []struct {
		name      string
		boundary  []h3.LatLng
		expected  []string
		expectErr bool
	}{
		{
			name: "polygon sylhet city",
			boundary: []h3.LatLng{
				{Lat: 24.914885929941853, Lng: 91.81924598937991},
				{Lat: 24.926477924225544, Lng: 91.83419538623046},
				{Lat: 24.926581810186146, Lng: 91.8560237751465},
				{Lat: 24.916333156927557, Lng: 91.86961241796875},
				{Lat: 24.91071212084026, Lng: 91.9029199897461},
				{Lat: 24.901058906895923, Lng: 91.91733954541016},
				{Lat: 24.87801614912562, Lng: 91.90322656616209},
				{Lat: 24.874470871319335, Lng: 91.8602034969482},
				{Lat: 24.8843165818294, Lng: 91.83030691479493},
			},
			
			expected: []string{
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
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function with the polygon boundary
			hexes, err := utils.FillPolygonWithHexes(tt.boundary)

			if tt.expectErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				// Check if the hexes match the expected result
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expected, hexes, "The returned hexes do not match the expected hexes")
			}
		})
	}
}
