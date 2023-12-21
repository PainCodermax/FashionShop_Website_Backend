package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Init() {
	// Tạo HTTP client
	client := &http.Client{}

	// Tạo yêu cầu HTTP
	req, err := http.NewRequest("GET", "https://online-gateway.ghn.vn/shiip/public-api/master-data/province", nil)
	if err != nil {
		fmt.Println("Lỗi khi tạo yêu cầu:", err)
		return
	}

	// Thêm header 'token' vào yêu cầu
	req.Header.Add("token", "5c4242e4-9bf5-11ee-96dc-de6f804954c9")

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
	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Lỗi khi decode JSON:", err)
		return
	}

	// In ra dữ liệu hoặc lưu vào một map
	provinceMap := make(map[string]interface{})
	for _, province := range data {
		provinceName, ok := province["ProvinceName"].(string)
		if !ok {
			fmt.Println("Lỗi: Không thể lấy tên tỉnh/thành phố.")
			continue
		}
		provinceMap[provinceName] = province
	}

	// In ra dữ liệu trong map
	for provinceName, provinceData := range provinceMap {
		fmt.Printf("Tỉnh/Thành phố: %s - Dữ liệu: %v\n", provinceName, provinceData)
	}
}
