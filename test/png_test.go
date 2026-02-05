package ruyi

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

// TestPngToJpg  png 转 jpg
// @param t
func TestPngToJpg(t *testing.T) {

	t.Run("png to jpg case 1", func(t *testing.T) {
		err := os.MkdirAll("testdata/output", 0755)
		if err != nil {
			t.Fatalf("testdata/output dir create failed, err = %v", err)
		}

		pngFile, err := os.Open("testdata/monkey.png")
		if err != nil {
			t.Fatalf("pngFile open failed, err=%v", err)
		}
		defer pngFile.Close()

		// 2. 解码 PNG
		// 注意：png 包会自动处理透明度。JPG 不支持透明，通常会变成黑色背景。
		img, err := png.Decode(pngFile)
		if err != nil {
			t.Fatalf("png Decode failed, err=%v", err)
		}

		// 3. 创建 JPG 目标文件
		out, err := os.Create("testdata/output/monkey_2.jpg")
		if err != nil {
			t.Fatalf("jpg file create failed, err=%v", err)
		}
		defer out.Close()

		// 4. 以 JPG 格式编码并保存
		// 设置质量参数
		opt := jpeg.Options{
			Quality: 100,
		}

		err = jpeg.Encode(out, img, &opt)
		if err != nil {
			t.Fatalf("jpg Encode failed, err=%v", err)
		}

	})

	// png支持透明，转jpg 不支持透明，手动指定背景色为白色
	t.Run("with_white_background", func(t *testing.T) {
		err := os.MkdirAll("testdata/output", 0755)
		if err != nil {
			t.Fatalf("testdata/output dir create failed, err = %v", err)
		}

		pngFile, err := os.Open("testdata/shop.png")
		if err != nil {
			t.Fatalf("pngFile open failed, err=%v", err)
		}
		defer pngFile.Close()

		// 2. 解码 PNG
		// 注意：png 包会自动处理透明度。JPG 不支持透明，通常会变成黑色背景。
		img, err := png.Decode(pngFile)
		if err != nil {
			t.Fatalf("png Decode failed, err=%v", err)
		}

		// 创建一个和原图一样大的新画布
		newImg := image.NewRGBA(img.Bounds())

		// 画布涂满白色
		draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

		// 把原图覆盖上去（Over 模式，保留非透明部分的颜色）
		draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)

		// 3. 创建 JPG 目标文件
		out, err := os.Create("testdata/output/shop_with_white_bg.jpg")
		if err != nil {
			t.Fatalf("jpg file create failed, err=%v", err)
		}
		defer out.Close()

		// 4. 以 JPG 格式编码并保存
		// 设置质量参数
		opt := jpeg.Options{
			Quality: 100,
		}

		err = jpeg.Encode(out, newImg, &opt)
		if err != nil {
			t.Fatalf("jpg Encode failed, err=%v", err)
		}

	})

	// png支持透明，转jpg 不支持透明，默认背景色是黑色
	t.Run("with_default_background", func(t *testing.T) {
		err := os.MkdirAll("testdata/output", 0755)
		if err != nil {
			t.Fatalf("testdata/output dir create failed, err = %v", err)
		}

		pngFile, err := os.Open("testdata/shop.png")
		if err != nil {
			t.Fatalf("pngFile open failed, err=%v", err)
		}
		defer pngFile.Close()

		// 2. 解码 PNG
		// 注意：png 包会自动处理透明度。JPG 不支持透明，通常会变成黑色背景。
		img, err := png.Decode(pngFile)
		if err != nil {
			t.Fatalf("png Decode failed, err=%v", err)
		}

		// 3. 创建 JPG 目标文件
		out, err := os.Create("testdata/output/shop_with_default_bg.jpg")
		if err != nil {
			t.Fatalf("jpg file create failed, err=%v", err)
		}
		defer out.Close()

		// 4. 以 JPG 格式编码并保存
		// 设置质量参数
		opt := jpeg.Options{
			Quality: 100,
		}

		err = jpeg.Encode(out, img, &opt)
		if err != nil {
			t.Fatalf("jpg Encode failed, err=%v", err)
		}

	})

}
