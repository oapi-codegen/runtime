package runtime

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InnerArrayObject struct {
	Names []string `json:"names"`
}

type InnerObject struct {
	Name string
	ID   int
}

type InnerObject2 struct {
	Foo string
	Is  bool
}

type InnerObject3 struct {
	Foo   string
	Count *int `json:"count,omitempty"`
}

// These are all possible field types, mandatory and optional.
type AllFields struct {
	I    int              `json:"i"`
	Oi   *int             `json:"oi,omitempty"`
	Ab   *[]bool          `json:"ab,omitempty"`
	F    float32          `json:"f"`
	Of   *float32         `json:"of,omitempty"`
	B    bool             `json:"b"`
	Ob   *bool            `json:"ob,omitempty"`
	As   []string         `json:"as"`
	Oas  *[]string        `json:"oas,omitempty"`
	O    InnerObject      `json:"o"`
	Ao   []InnerObject2   `json:"ao"`
	Aop  *[]InnerObject3  `json:"aop"`
	Onas InnerArrayObject `json:"onas"`
	Oo   *InnerObject     `json:"oo,omitempty"`
	D    MockBinder       `json:"d"`
	Od   *MockBinder      `json:"od,omitempty"`
	M    map[string]int   `json:"m"`
	Om   *map[string]int  `json:"om,omitempty"`
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
	d := MockBinder{Time: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)}

	two := 2

	srcObj := AllFields{
		I:   12,
		Oi:  &oi,
		F:   4.2,
		Of:  &of,
		B:   true,
		Ob:  &ob,
		Ab:  &[]bool{true},
		As:  []string{"hello", "world"},
		Oas: &oas,
		O: InnerObject{
			Name: "Joe Schmoe",
			ID:   456,
		},
		Ao: []InnerObject2{
			{Foo: "bar", Is: true},
			{Foo: "baz"},
		},
		Aop: &[]InnerObject3{
			{Foo: "a"},
			{Foo: "b", Count: &two},
		},
		Onas: InnerArrayObject{
			Names: []string{"Bill", "Frank"},
		},
		Oo: &oo,
		D:  d,
		Od: &d,
		M:  om,
		Om: &om,
	}

	marshaled, err := MarshalDeepObject(srcObj, "p")
	require.NoError(t, err)
	require.EqualValues(t, "p[ab][0]=true&p[ao][0][Foo]=bar&p[ao][0][Is]=true&p[ao][1][Foo]=baz&p[ao][1][Is]=false&p[aop][0][Foo]=a&p[aop][1][Foo]=b&p[aop][1][count]=2&p[as][0]=hello&p[as][1]=world&p[b]=true&p[d]=2020-02-01&p[f]=4.2&p[i]=12&p[m][additional]=1&p[o][ID]=456&p[o][Name]=Joe Schmoe&p[oas][0]=foo&p[oas][1]=bar&p[ob]=true&p[od]=2020-02-01&p[of]=3.7&p[oi]=5&p[om][additional]=1&p[onas][names][0]=Bill&p[onas][names][1]=Frank&p[oo][ID]=123&p[oo][Name]=Marcin Romaszewicz", marshaled)

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

// Item represents an item object for testing array of objects
type Item struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func TestDeepObject_ArrayOfObjects(t *testing.T) {
	// Test case for:
	// name: items
	// style: deepObject
	// required: false
	// explode: true
	// schema:
	//   type: array
	//   minItems: 1
	//   items:
	//     type: object

	srcArray := []Item{
		{Name: "first", Value: "value1"},
		{Name: "second", Value: "value2"},
	}

	// Marshal the array to deepObject format
	marshaled, err := MarshalDeepObject(srcArray, "items")
	require.NoError(t, err)
	t.Log("Marshaled:", marshaled)

	// Expected format for array of objects with explode:true should be:
	// items[0][name]=first&items[0][value]=value1&items[1][name]=second&items[1][value]=value2

	// Parse the marshaled string into url.Values
	params := make(url.Values)
	marshaledParts := strings.Split(marshaled, "&")
	for _, p := range marshaledParts {
		parts := strings.Split(p, "=")
		require.Equal(t, 2, len(parts))
		params.Set(parts[0], parts[1])
	}

	// Unmarshal back to the destination array
	var dstArray []Item
	err = UnmarshalDeepObject(&dstArray, "items", params)
	require.NoError(t, err)

	// Verify the result matches the source
	assert.EqualValues(t, srcArray, dstArray)
	assert.Len(t, dstArray, 2)
	assert.Equal(t, "first", dstArray[0].Name)
	assert.Equal(t, "value1", dstArray[0].Value)
	assert.Equal(t, "second", dstArray[1].Name)
	assert.Equal(t, "value2", dstArray[1].Value)
}
