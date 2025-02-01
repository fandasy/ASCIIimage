Hi reader, this is my package to convert local or web image to ASCII character image

Supports formats: png, jpg, jpeg, webp 

```
go get github.com/fandasy/ASCIIimage
```

Functions
---
- GetFromFile (ctx context.Context, path string, opts Options) (*image.RGBA, error)
```
GetFromFile reads an image from a file and converts it to an ASCII art image.

Possible output errors:
ErrFileNotFound
ErrIncorrectFormat
```

- GetFromWebsite (ctx context.Context, url string, opts Options) (*image.RGBA, error)
```
GetFromWebsite downloads an image from a URL and converts it to an ASCII art image.

Possible output errors:
ErrIncorrectUrl
ErrPageNotFound
ErrIncorrectFormat
```

Struct

```
type Options struct {
    Compress  uint8   // 0-99
    MaxWidth  uint    // 1 = 10px
    MaxHeight uint    //
    Chars     string  // dark to light, e.g., '@%#*+=:~-. '
}
```

#### Remark
```
--------------------------------------------------
 There are default values for these parameters:

 Compress  = 0
 MaxWidth  = 10000 -> 100000px
 MaxHeight = 10000 -> 100000px
 Chars     = "@%#*+=:~-. "

--------------------------------------------------
 Default values can be activated by specifying:

 Compress  >= 100
 MaxWidth  == 0
 MaxHeight == 0
 Chars     == ""
--------------------------------------------------
```

Example
---
```Go
asciiImage, err := asciiimage.GetFromWebsite(
    context.TODO(),
    "https://your_domain.com/file.jpg",
    asciiimage.Options{
        Compress:  0,
        MaxWidth:  0,  // <- default value is activated
        MaxHeight: 0,  // <- |
        Chars:     "", // <- /
    })
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
