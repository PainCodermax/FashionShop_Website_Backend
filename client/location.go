package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

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

// func createRequest(method, url string, payload interface{}, token, shopId, endpoint string) (*http.Request, error) {
// 	var reqBody *bytes.Buffer
// 	if payload != nil {
// 		jsonPayload, err := json.Marshal(payload)
// 		if err != nil {
// 			return nil, fmt.Errorf("Error marshaling payload: %v", err)
// 		}
// 		reqBody = bytes.NewBuffer(jsonPayload)
// 	} else {
// 		reqBody = new(bytes.Buffer)
// 	}
// 	url = url + endpoint
// 	req, err := http.NewRequest(method, url, reqBody)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error creating request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Token", token)
// 	req.Header.Set("ShopId", shopId)

// 	return req, nil
// }

// func doRequest(client *http.Client, req *http.Request) (*http.Response, error) {
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		defer resp.Body.Close()
// 		return nil, fmt.Errorf("Error: Unable to fetch data. Status code: %d", resp.StatusCode)
// 	}
// 	return resp, nil
// }

// func GetAddressString(province, district, ward string) string {
// 	url := "https://online-gateway.ghn.vn/shiip/public-api/master-data/"

// 	provinceInt, _ := utils.ParseStringToIn64(province)
// 	districtInt, _ := utils.ParseStringToIn64(district)
// 	// wardInt, _ := utils.ParseStringToIn64(ward)

// 	districtPayload := map[string]interface{}{
// 		"province_id": provinceInt,
// 	}

// 	wardPayload := map[string]interface{}{
// 		"district_id": districtInt,
// 	}

// 	client := &http.Client{}

// 	reqProvince, err := createRequest("GET", url, nil, "5c4242e4-9bf5-11ee-96dc-de6f804954c9", "4771536", "province")
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return ""
// 	}

// 	reqDistrict, err := createRequest("GET", url, districtPayload, "5c4242e4-9bf5-11ee-96dc-de6f804954c9", "4771536", "district")
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return ""
// 	}

// 	reqWard, err := createRequest("GET", url, wardPayload, "5c4242e4-9bf5-11ee-96dc-de6f804954c9", "4771536", "ward")
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return ""
// 	}

// 	// Process province request
// 	provinceRes := models.ProvinceResponse{}
// 	resp, err := doRequest(client, reqProvince)
// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		return ""
// 	}
// 	defer resp.Body.Close()
// 	if err := json.NewDecoder(resp.Body).Decode(&provinceRes); err != nil {
// 		fmt.Println("Error decoding JSON:", err)
// 		return ""
// 	}

// 	// Process district request
// 	districtRes := models.DistrictResponse{}
// 	resp, err = doRequest(client, reqDistrict)
// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		return ""
// 	}
// 	defer resp.Body.Close()
// 	if err := json.NewDecoder(resp.Body).Decode(&districtRes); err != nil {
// 		fmt.Println("Error decoding JSON:", err)
// 		return ""
// 	}

// 	// Process ward request
// 	wardRes := models.WardResponse{}
// 	resp, err = doRequest(client, reqWard)
// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		return ""
// 	}
// 	defer resp.Body.Close()
// 	if err := json.NewDecoder(resp.Body).Decode(&wardRes); err != nil {
// 		fmt.Println("Error decoding JSON:", err)
// 		return ""
// 	}

// 	for _, p := range provinceRes.Data {
// 		if p.ProvinceID == int(provinceInt) {
// 			province = p.NameExtension[2]
// 			break
// 		}
// 	}

// 	for _, d := range districtRes.Data {
// 		if d.DistrictID == int(districtInt) {
// 			district = d.DistrictName
// 			break
// 		}
// 	}

// 	for _, w := range wardRes.Data {
// 		if w.WardID == ward {
// 			ward = w.WardName
// 			break
// 		}
// 	}

// 	return fmt.Sprintf("%s, %s, %s", ward, district, province)
// }

func createRequest(method, url, token, shopId, endpoint string, payload interface{}) (*http.Request, error) {
	var reqBody *bytes.Buffer
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling payload: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonPayload)
	} else {
		reqBody = new(bytes.Buffer)
	}
	fullUrl := url + endpoint
	req, err := http.NewRequest(method, fullUrl, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	req.Header.Set("ShopId", shopId)

	return req, nil
}

type RequestResult struct {
	Resp  *http.Response
	Err   error
	Which string
}

func doRequest(client *http.Client, req *http.Request, which string, wg *sync.WaitGroup, ch chan<- RequestResult) {
	defer wg.Done()
	resp, err := client.Do(req)
	ch <- RequestResult{Resp: resp, Err: err, Which: which}
}

func GetAddressString(province, district, ward string) string {
	url := "https://online-gateway.ghn.vn/shiip/public-api/master-data/"
	token := "5c4242e4-9bf5-11ee-96dc-de6f804954c9"
	shopId := "4771536"

	provinceInt, _ := utils.ParseStringToIn64(province)
	districtInt, _ := utils.ParseStringToIn64(district)
	// wardInt, _ := utils.ParseStringToIn64(ward)

	districtPayload := map[string]interface{}{
		"province_id": provinceInt,
	}

	wardPayload := map[string]interface{}{
		"district_id": districtInt,
	}

	client := &http.Client{}
	var wg sync.WaitGroup
	ch := make(chan RequestResult, 3)

	// Create and start requests concurrently
	reqProvince, err := createRequest("GET", url, token, shopId, "province", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	wg.Add(1)
	go doRequest(client, reqProvince, "province", &wg, ch)

	reqDistrict, err := createRequest("GET", url, token, shopId, "district", districtPayload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	wg.Add(1)
	go doRequest(client, reqDistrict, "district", &wg, ch)

	reqWard, err := createRequest("GET", url, token, shopId, "ward", wardPayload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	wg.Add(1)
	go doRequest(client, reqWard, "ward", &wg, ch)

	// Wait for all requests to complete
	go func() {
		wg.Wait()
		close(ch)
	}()

	var provinceRes models.ProvinceResponse
	var districtRes models.DistrictResponse
	var wardRes models.WardResponse

	// Process responses as they come in
	for result := range ch {
		if result.Err != nil {
			fmt.Println("Error making request:", result.Err)
			continue
		}
		defer result.Resp.Body.Close()

		if result.Resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: Unable to fetch data for %s. Status code: %d\n", result.Which, result.Resp.StatusCode)
			continue
		}

		switch result.Which {
		case "province":
			if err := json.NewDecoder(result.Resp.Body).Decode(&provinceRes); err != nil {
				fmt.Println("Error decoding JSON:", err)
				return ""
			}
		case "district":
			if err := json.NewDecoder(result.Resp.Body).Decode(&districtRes); err != nil {
				fmt.Println("Error decoding JSON:", err)
				return ""
			}
		case "ward":
			if err := json.NewDecoder(result.Resp.Body).Decode(&wardRes); err != nil {
				fmt.Println("Error decoding JSON:", err)
				return ""
			}
		}
	}

	for _, p := range provinceRes.Data {
		if p.ProvinceID == int(provinceInt) {
			province = p.NameExtension[2]
			break
		}
	}

	for _, d := range districtRes.Data {
		if d.DistrictID == int(districtInt) {
			district = d.DistrictName
			break
		}
	}

	for _, w := range wardRes.Data {
		if w.WardID == ward {
			ward = w.WardName
			break
		}
	}

	return fmt.Sprintf("%s, %s, %s", ward, district, province)
}
