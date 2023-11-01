package bare_metal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertToListMap_Configurations(t *testing.T) {
	configurations := []Configuration{
		{
			Name:         "test-conf-1",
			Cores:        10,
			Memory:       "256 GB",
			Storage:      "1 TB",
			Disks:        1,
			Nics:         1,
			Processors:   1,
			MachineCount: 10,
			Price: Price{
				Amount: 1.25,
				Period: "HOURLY",
			},
		},
	}

	list := ConvertToListMap(configurations)

	assert.Equal(t, len(configurations), len(list))

	conf := configurations[0]
	convertedConf := list[0]
	assert.Equal(t, conf.Name, convertedConf["name"])
	assert.Equal(t, conf.Cores, convertedConf["cores"])
	assert.Equal(t, conf.Memory, convertedConf["memory"])
	assert.Equal(t, conf.Storage, convertedConf["storage"])
	assert.Equal(t, conf.Disks, convertedConf["disks"])
	assert.Equal(t, conf.Nics, convertedConf["nics"])
	assert.Equal(t, conf.Processors, convertedConf["processors"])
	assert.Equal(t, conf.MachineCount, convertedConf["machine_count"])
	assert.Equal(t, "$1.25/HOURLY", convertedConf["price"])
}

func TestConvertToListMap_Locations(t *testing.T) {
	locations := []Location{
		{
			ID:     "TEST_ID",
			Name:   "TEST SITE 1",
			Status: "Status",
			Region: "NA",
		},
	}

	list := ConvertToListMap(locations)

	assert.Equal(t, len(locations), len(list))

	location := locations[0]
	convertedLocation := list[0]
	assert.Equal(t, location.ID, convertedLocation["id"])
	assert.Equal(t, location.Name, convertedLocation["name"])
	assert.Equal(t, location.Status, convertedLocation["status"])
	assert.Equal(t, location.Region, convertedLocation["region"])
}

func TestConvertToListMap_NetworkSize(t *testing.T) {
	networkSizes := []NetworkSize{
		{
			ID:           "IN28",
			Name:         "Lumen Network IP Block /28",
			CIDR:         "/28",
			AvailableIPs: 11,
			Price: Price{
				Amount: 0.08,
				Period: "HOURLY",
			},
		},
	}

	list := ConvertToListMap(networkSizes)

	assert.Equal(t, len(networkSizes), len(list))

	networkSize := networkSizes[0]
	convertedNS := list[0]
	assert.Equal(t, networkSize.ID, convertedNS["id"])
	assert.Equal(t, networkSize.Name, convertedNS["name"])
	assert.Equal(t, networkSize.CIDR, convertedNS["cidr"])
	assert.Equal(t, networkSize.AvailableIPs, convertedNS["available_ips"])
	assert.Equal(t, "$0.08/HOURLY", convertedNS["price"])
}

func TestConvertToListMap_OsImage(t *testing.T) {
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

	list := ConvertToListMap(osImages)

	assert.Equal(t, len(osImages), len(list))

	osImage := osImages[0]
	convertedImage := list[0]
	assert.Equal(t, osImage.Name, convertedImage["name"])
	assert.Equal(t, osImage.Ready, convertedImage["ready"])
	assert.Equal(t, "$45.00/MONTHLY", convertedImage["price"])
}
