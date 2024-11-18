Hi reader, this is my package to convert local or web image to ASCII character image

Supports formats: png, jpg, jpeg, webp 

```
go get github.com/fandasy/ASCIIimage
```

Functions
---
- GetFromFile (string, float64, int, int, string) (*image.RGBA, error)
```
GetFromFile takes
 path to the image,
 compression percentage (0.0 - 1.0),
 maximum width (1 = 10px),
 maximum height (1 = 10px),
 chars that will be used to generate (dark - light)

 Possible output errors:
 ErrFileNotFound,
 ErrIncorrectFormat
 and other error
```

- GetFromWebsite (ctx, string, float64, int, int, string) (*image.RGBA, error)
```
GetFromWebsite takes
 context,
 image url,
 compression percentage (0.0 - 1.0),
 maximum width (1 = 10px),
 maximum height (1 = 10px),
 chars that will be used to generate (dark - light)

 Possible output errors:
 ErrIncorrectUrl,
 ErrPageNotFound,
 ErrIncorrectFormat
 and other error
```

#### Remark
```
 --------------------------------------------------
 There are default values for these parameters:

 reductionPercentage = 0.0
 maxWidth  = 5000 -> 50000px
 maxHeight = 5000 -> 50000px
 chars     = "@%#*+=:~-. "

 --------------------------------------------------
 Default values can be activated by specifying:

 reductionPercentage < 0 || reductionPercentage > 1
 maxWidth  <= 0
 maxHeight <= 0
 chars     == ""
 --------------------------------------------------
```

Example
---
```GO
asciiImage, err := asciiimage.GetFromWebsite(
	context.TODO(),
	"https://ir.ozone.ru/s3/multimedia-7/c1000/6755179327.jpg",
	0,
	0,
	0,
	"",
)
if err != nil {
	log.Println(err)
	return
}

file, err := os.Create("example.jpg")
if err != nil {
	log.Println("failed to create file: ", err)
	return
}

if err := jpeg.Encode(file, asciiImage, nil); err != nil {
	log.Println("failed to encode image: ", err)
	return
}
```

Note
---
Additionally, I would like to thank the nfnt developer for the package github.com/nfnt/resize
