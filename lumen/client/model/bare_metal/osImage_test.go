package bare_metal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertOSImagesToListMap(t *testing.T) {
	osImages := []OsImage{
		{
			Name:  "Ubuntu 20.04",
			Ready: true,
			Price: Price{
				Amount: 45.00,
				Period: "MONTHLY",
			},
		},
	}

	list := ConvertOSImagesToListMap(osImages)

	assert.Equal(t, len(osImages), len(list))

	osImage := osImages[0]
	convertedImage := list[0]
	assert.Equal(t, osImage.Name, convertedImage["name"])
	assert.Equal(t, "$45.00/MONTHLY", convertedImage["price"])
}

func TestConvertOSImagesToListMap_TieredPrice(t *testing.T) {
	tier1 := 1
	tier2 := 2
	osImages := []OsImage{
		{
			Name:  "RHEL 7.9",
			Ready: true,
			Prices: []TierPrice{
				{
					Price: Price{
						Amount: 45.00,
						Period: "MONTHLY",
					},
					Tier: &tier1,
				},
				{
					Price: Price{
						Amount: 90.00,
						Period: "MONTHLY",
					},
					Tier: &tier2,
				},
			},
		},
	}

	list := ConvertOSImagesToListMap(osImages)

	assert.Equal(t, 2, len(list))

	tierOneOSImage := osImages[0]
	convertedTierOneOSImage := list[0]
	assert.Equal(t, tierOneOSImage.Name, convertedTierOneOSImage["name"])
	assert.Equal(t, "$45.00/MONTHLY", convertedTierOneOSImage["price"])
	convertedTierTwoOSImage := list[1]
	assert.Equal(t, tierOneOSImage.Name, convertedTierTwoOSImage["name"])
	assert.Equal(t, "$90.00/MONTHLY", convertedTierTwoOSImage["price"])
}
