package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PainCodermax/FashionShop_Website_Backend/models"
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


