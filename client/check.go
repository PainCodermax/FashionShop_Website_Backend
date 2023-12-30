package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GOs() {
	// Khai báo các biến để lưu giá trị từ flag

	url := "https://online-gateway.ghn.vn/shiip/public-api/v2/shipping-order/fee"

	payload := map[string]interface{}{
		"service_type_id":  5,
		"from_district_id": 1442,
		"to_district_id":   1455,
		"to_ward_code":     "21412",
		"height":           20,
		"length":           30,
		"weight":           3000,
		"width":            40,
		"insurance_value":  0,
		"coupon":           nil,
		"items": []map[string]interface{}{
			{
				"name":     "TEST1",
				"quantity": 1,
				"height":   200,
				"weight":   1000,
				"length":   200,
				"width":    200,
			},
		},
	}

	reqBody, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "5c4242e4-9bf5-11ee-96dc-de6f804954c9")
	req.Header.Set("ShopId", "4771536")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("Response Body:", result)
}

