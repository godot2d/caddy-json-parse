package jsonparse

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFetchers(t *testing.T) {
	tests := []struct {
		json     string
		key      string
		expected interface{}
	}{
		{
			json:     `{"ref":"ok"}`,
			key:      "ref",
			expected: "ok",
		},
		{
			json:     `[7,8,9,0]`,
			key:      "2",
			expected: float64(9),
		},
		{
			json:     `["what", "is", "this"]`,
			key:      "2",
			expected: "this",
		},
		{
			json:     `{"ref": [5,8,9]}`,
			key:      "ref.1",
			expected: float64(8),
		},
		{
			json:     `{"ref": {"joe": [1,2, {"sum": 100 } ]}}`,
			key:      "ref.joe.2.sum",
			expected: float64(100),
		},
		{
			json:     `{"ref": {"joe": [1,2, {"sum": {"100": {"dave" : "lee"}}  } ]}}`,
			key:      "ref.joe.2.sum.100.dave",
			expected: "lee",
		},
		{
			json:     `{"serverName":"example.com"}`,
			key:      "serverName",
			expected: "example.com",
		},
		{
			json:     `{"config": {"serverName": "api.example.com"}}`,
			key:      "config.serverName",
			expected: "api.example.com",
		},
	}

	for i, tt := range tests {
		var v interface{}
		err := json.Unmarshal([]byte(tt.json), &v)
		if err != nil {
			t.Fatal(err)
			continue
		}
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			val := fetchValue(v, tt.key)
			if val != tt.expected {
				t.Errorf("want: %v, got: %v", tt.expected, val)
			}
		})
	}

}

func TestAttachServerNameExtraction(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		key      string
		expected interface{}
	}{
		{
			name:     "extract ServerName from attach field",
			json:     `{"game": "com.arrow.defense3d", "orderId": "order-id", "uid": "test", "amount": "100", "tradeState": "SUCCESS", "timestamp": 123456789, "platform": "ios", "sku": "whatever", "attach": "{ \"roleId\": 100016, \"gameServerId\": 1, \"shopId\": 1, \"merchandiseId\": 1001 , \"ServerName\" : \"dev\"}"}`,
			key:      "attach.ServerName",
			expected: "dev",
		},
		{
			name:     "extract ServerName from attach field with different value",
			json:     `{"attach": "{ \"ServerName\": \"production\" }"}`,
			key:      "attach.ServerName",
			expected: "production",
		},
		{
			name:     "no attach field",
			json:     `{"game": "com.arrow.defense3d"}`,
			key:      "attach.ServerName",
			expected: nil,
		},
		{
			name:     "attach field without ServerName",
			json:     `{"attach": "{ \"roleId\": 100016 }"}`,
			key:      "attach.ServerName",
			expected: nil,
		},
		{
			name:     "invalid JSON in attach field",
			json:     `{"attach": "invalid json"}`,
			key:      "attach.ServerName",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v interface{}
			err := json.Unmarshal([]byte(tt.json), &v)
			if err != nil {
				t.Fatal(err)
			}

			// Simulate the replacer function logic
			var serverName interface{}
			if attachVal := fetchValue(v, "attach"); attachVal != nil {
				if attachStr, ok := attachVal.(string); ok {
					var attachData map[string]interface{}
					if err := json.Unmarshal([]byte(attachStr), &attachData); err == nil {
						if sn, exists := attachData["ServerName"]; exists {
							serverName = sn
						}
					}
				}
			}

			if serverName != tt.expected {
				t.Errorf("want: %v, got: %v", tt.expected, serverName)
			}
		})
	}
}
