package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
)

func Init() {
	ProvinceMap := make(map[int]string, 0)
	token := "5c4242e4-9bf5-11ee-96dc-de6f804954c9"
	// Tạo HTTP client
	client := &http.Client{}

	// params := url.Values{}
	// params.Set("paramKey", "paramValue")

	// // Thêm tham số vào URL
	// reqURL := "https://online-gateway.ghn.vn/shiip/public-api/master-data/province"
	// reqURL += "?" + params.Encode()

	req, err := http.NewRequest("GET", "https://online-gateway.ghn.vn/shiip/public-api/master-data/province", nil)
	if err != nil {
		fmt.Println("Lỗi khi tạo yêu cầu:", err)
		return
	}
	req.Header.Add("token", token)

	// Thực hiện yêu cầu HTTP
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Lỗi khi gọi API:", err)
		return
	}
	defer resp.Body.Close()
	// Đảm bảo response có status code 200 OK
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Lỗi: Không thể lấy dữ liệu. Mã trạng thái:", resp.StatusCode)
		return
	}

	// Decode JSON response vào một map[string]interface{}
	var address models.ProvinceResponse
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		fmt.Println("Lỗi khi decode JSON:", err)
		return
	}

	// In ra dữ liệu hoặc lưu vào một map
	for _, province := range address.Data {
		ProvinceMap[province.ProvinceID] = province.ProvinceName
	}
}

func CheckShipingFee(province, district, ward string) int {
	url := "https://online-gateway.ghn.vn/shiip/public-api/v2/shipping-order/fee"

	districtInt, _ := utils.ParseStringToIn64(district)
	payload := map[string]interface{}{
		"service_type_id":  2,
		"from_district_id": 1442,
		"to_district_id":   districtInt,
		"to_ward_code":     ward,
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
		return -1
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "5c4242e4-9bf5-11ee-96dc-de6f804954c9")
	req.Header.Set("ShopId", "4771536")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return -1
	}
	defer resp.Body.Close()

	// Đảm bảo response có status code 200 OK
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Lỗi: Không thể lấy dữ liệu. Mã trạng thái:", resp.StatusCode)
		return -1
	}

	// Decode JSON response vào một map[string]interface{}
	var shipment models.ShipmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&shipment); err != nil {
		fmt.Println("Lỗi khi decode JSON:", err)
		return -1
	}

	// In ra dữ liệu hoặc lưu vào một map
	return shipment.Data.Total
}
