package ingresses

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpsertPaths(t *testing.T) {

	updateItem := Path{
		UID:      "server1",
		Address:  "test.org",
		Endpoint: Endpoint{Name: "server1", Port: 3333},
	}

	addItem := Path{
		UID:      "server2",
		Address:  "test.org",
		Endpoint: Endpoint{Name: "server2", Port: 7777},
	}

	list := []Path{
		{
			UID:      "server0",
			Address:  "test.org",
			Endpoint: Endpoint{Name: "server0", Port: 3333},
		},
		{
			UID:      "server1",
			Address:  "test.org",
			Endpoint: Endpoint{Name: "server1", Port: 1234},
		},
	}

	expected := []Path{
		{
			UID:      "server0",
			Address:  "test.org",
			Endpoint: Endpoint{Name: "server0", Port: 3333},
		},
		updateItem,
		addItem,
	}

	result := UpsertPaths(list, []Path{updateItem, addItem})

	assert.Equal(t, expected, result)
}

func TestUpsertPathsWithSameHost(t *testing.T) {
	list := []Path{
		{
			UID:      "gazebosim.org",
			Address:  "/example/0",
			Endpoint: Endpoint{Name: "server0", Port: 3333},
		},
		{
			UID:      "gazebosim.org",
			Address:  "/example/1",
			Endpoint: Endpoint{Name: "server1", Port: 1234},
		},
	}

	item := Path{
		UID:      "gazebosim.org",
		Address:  "/example/2",
		Endpoint: Endpoint{Name: "server2", Port: 3333},
	}

	result := UpsertPaths(list, []Path{item})
	assert.Len(t, result, 3)
}

func TestRemovePaths(t *testing.T) {
	removeItem := Path{
		UID:      "server1",
		Address:  "test.org",
		Endpoint: Endpoint{Name: "server1", Port: 1234},
	}

	list := []Path{
		{
			UID:      "server0",
			Address:  "test.org",
			Endpoint: Endpoint{Name: "server0", Port: 3333},
		},
		{
			Address:  "test.org",
			Endpoint: Endpoint{Name: "server1", Port: 1234},
		},
	}

	expected := []Path{
		{
			UID:      "server0",
			Address:  "test.org",
			Endpoint: Endpoint{Name: "server0", Port: 3333},
		},
	}

	result := RemovePaths(list, []Path{removeItem})

	assert.Equal(t, expected, result)
}
