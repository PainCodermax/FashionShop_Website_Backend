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
	var data []models.Province
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Lỗi khi decode JSON:", err)
		return
	}

	// In ra dữ liệu hoặc lưu vào một map
	for _, province := range data {
		ProvinceMap[province.ProvinceID] = province.ProvinceName
		fmt.Print(province.ProvinceName)
	}

}
