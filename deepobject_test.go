package runtime

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InnerObject struct {
	Name string
	ID   int
}

// These are all possible field types, mandatory and optional.
type AllFields struct {
	// Primitive types
	I   int             `json:"i"`
	Oi  *int            `json:"oi,omitempty"`
	F   float32         `json:"f"`
	Of  *float32        `json:"of,omitempty"`
	B   bool            `json:"b"`
	Ob  *bool           `json:"ob,omitempty"`
	As  []string        `json:"as"`
	Oas *[]string       `json:"oas,omitempty"`
	O   InnerObject     `json:"o"`
	Oo  *InnerObject    `json:"oo,omitempty"`
	M   map[string]int  `json:"m"`
	Om  *map[string]int `json:"om,omitempty"`

	// Complex types
	Bi  MockBinder  `json:"bi"`
	Obi *MockBinder `json:"obi,omitempty"`
	Da  types.Date  `json:"da"`
	Oda *types.Date `json:"oda,omitempty"`
	Ti  time.Time   `json:"ti"`
	Oti *time.Time  `json:"oti,omitempty"`
	U   types.UUID  `json:"u"`
	Ou  *types.UUID `json:"ou,omitempty"`

	// Nullable
	NiSet   nullable.Nullable[int]         `json:"ni_set,omitempty"`
	NiNull  nullable.Nullable[int]         `json:"ni_null,omitempty"`
	NiUnset nullable.Nullable[int]         `json:"ni_unset,omitempty"`
	No      nullable.Nullable[InnerObject] `json:"no,omitempty"`
	Nu      nullable.Nullable[uuid.UUID]   `json:"nu,omitempty"`
}

func TestDeepObject(t *testing.T) {
	oi := int(5)
	of := float32(3.7)
	ob := true
	oas := []string{"foo", "bar"}
	oo := InnerObject{
		Name: "Marcin Romaszewicz",
		ID:   123,
	}
	om := map[string]int{
		"additional": 1,
	}

	bi := MockBinder{Time: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)}
	da := types.Date{Time: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC)}
	ti := time.Now().UTC()
	u := uuid.New()

	srcObj := AllFields{
		// Primitive types
		I:   12,
		Oi:  &oi,
		F:   4.2,
		Of:  &of,
		B:   true,
		Ob:  &ob,
		As:  []string{"hello", "world"},
		Oas: &oas,
		O: InnerObject{
			Name: "Joe Schmoe",
			ID:   456,
		},
		Oo: &oo,
		M:  om,
		Om: &om,

		// Complex types
		Bi:  bi,
		Obi: &bi,
		Da:  da,
		Oda: &da,
		Ti:  ti,
		Oti: &ti,
		U:   u,
		Ou:  &u,

		// Nullable
		NiSet:  nullable.NewNullableWithValue(5),
		NiNull: nullable.NewNullNullable[int](),
		No: nullable.NewNullableWithValue(InnerObject{
			Name: "John Smith",
			ID:   456,
		}),
		Nu: nullable.NewNullableWithValue(uuid.New()),
	}

	marshaled, err := MarshalDeepObject(srcObj, "p")
	require.NoError(t, err)
	t.Log(marshaled)

	params := make(url.Values)
	marshaledParts := strings.Split(marshaled, "&")
	for _, p := range marshaledParts {
		parts := strings.Split(p, "=")
		require.Equal(t, 2, len(parts))
		params.Set(parts[0], parts[1])
	}

	var dstObj AllFields
	err = UnmarshalDeepObject(&dstObj, "p", params)
	require.NoError(t, err)
	assert.EqualValues(t, srcObj, dstObj)
}
