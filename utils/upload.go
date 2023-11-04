package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func Credentials() (*cloudinary.Cloudinary, context.Context) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, _ := cloudinary.New()
	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx
}

func UploadImage(cld *cloudinary.Cloudinary, data string, idx int, ctx context.Context) string {

	// imageFilePath := "image.jpeg"
	// imageFile, err := ioutil.ReadFile(imageFilePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Encode the image data to Base64
	// base64Image := base64.StdEncoding.EncodeToString(imageFile)

	// // Create a Base64 URL
	// base64URL := "data:image/jpeg;base64," + url.PathEscape(base64Image)

	// Upload the image.
	// Set the asset's public ID and allow overwriting the asset with new versions
	now := time.Now().UTC().String()
	resp, err := cld.Upload.Upload(ctx, data, uploader.UploadParams{
		PublicID:       fmt.Sprintf("product-%s-%d", now, idx),
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true)})
	if err != nil {
		fmt.Println("error")
	}
	return resp.SecureURL
}

// func getAssetInfo(cld *cloudinary.Cloudinary, ctx context.Context) {
// 	// Get and use details of the image
// 	// ==============================
// 	resp, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "quickstart_butterfly"})
// 	if err != nil {
// 		fmt.Println("error")
// 	}
// 	fmt.Println("****3. Get and use details of the image****\nDetailed response:\n", resp, "\n")

// 	// Assign tags to the uploaded image based on its width. Save the response to the update in the variable 'update_resp'.
// 	if resp.Width > 900 {
// 		update_resp, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
// 			PublicID: "quickstart_butterfly",
// 			Tags:     []string{"large"}})
// 		if err != nil {
// 			fmt.Println("error")
// 		} else {
// 			// Log the new tag to the console.
// 			fmt.Println("New tag: ", update_resp.Tags, "\n")
// 		}
// 	} else {
// 		update_resp, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
// 			PublicID: "quickstart_butterfly",
// 			Tags:     []string{"small"}})
// 		if err != nil {
// 			fmt.Println("error")
// 		} else {
// 			// Log the new tag to the console.
// 			fmt.Println("New tag: ", update_resp.Tags, "\n")
// 		}
// 	}

// }

// func transformImage(cld *cloudinary.Cloudinary, ctx context.Context) {
// 	// Instantiate an object for the asset with public ID "my_image"
// 	qs_img, err := cld.Image("quickstart_butterfly")
// 	if err != nil {
// 		fmt.Println("error")
// 	}

// 	// Add the transformation
// 	qs_img.Transformation = "r_max/e_sepia"

// 	// Generate and log the delivery URL
// 	new_url, err := qs_img.String()
// 	if err != nil {
// 		fmt.Println("error")
// 	} else {
// 		fmt.Println("****4. Transform the image****\nTransfrmation URL: ", new_url, "\n")
// 	}
// }
