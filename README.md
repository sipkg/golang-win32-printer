# golang-win32-printer
使用golang 语言封装win32 print API 支持打印图片,字符串，文件
# 开发目的
用于在windows 系统下 进行系统打印功能
最初想法是想为web 打印增加打印插件 以及获取系统信息,目前web打印功能太弱.

已封装win32 函数
1. DeleteDC
2. CreateDC
3. TextOut
4. StartDoc
5. EndDoc
6. StartPage
7. EndPage
8. OpenPrinter
9. ClosePrinter
10. StartDocPrinter
11. CloseDocPrinter
12. StartPagePrinter
13. ClosePagePrinter
14. ResetDC
15. SetPixel
16. GetPixel
17. GetDeviceCaps
18. StretchDIBits
19. MoveTo
20. LineTo
21. EnumPrinter
22. GetDefaultPrinter
23. SetDefaultPrinter

## 包结构
-  golang-win32-printer
     - image 封装 BGR 格式图像,支持24位BPP
     - printer win32 api 逻辑封装
     - win32 系统调用API封装 

## 目前打印流程
BGRImage 封装
通过 画图函数 把图片 文字 线条 矩形 写入BGRImage 再通过 StretchDIBits 复制到打印机HDC

## 代码示例
``` 
        printName := "Microsoft Print to PDF"
	dc, err := CreateDC(printName)
	fmt.Print(err)
	StartDCPrinter(dc, "gdiDoc")
	StartPage(dc)
	file, err := os.Open("C:\\Users\\wangjun\\Desktop\\USA.png")
	fmt.Print(err)
	image, err := png.Decode(file)
	fmt.Print(err)
	bgr := bgr2.NewBGRImage(image.Bounds())
	draw.Draw(bgr, image.Bounds(), image, image2.Point{0, 0}, draw.Src)
	src := bgr2.ReverseDIB(bgr.Pix, image.Bounds().Dx(), image.Bounds().Dy(), 24)
	DrawDIImage(dc, 0, uint32(image.Bounds().Dy())*10, uint32(image.Bounds().Dx())*10, uint32(image.Bounds().Dy())*10, 0, 0, int32(image.Bounds().Dx()), int32(image.Bounds().Dy()), src)
	EndPage(dc)
	EndDoc(dc)
	DeleteDC(dc)
```
功能正在完善中....欢迎参与改进
