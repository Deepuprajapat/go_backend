package request

import "io"

type UploadFileRequest struct {
	File        io.Reader `json:"file" form:"file"`
	AltKeywords string    `json:"alt_keywords" form:"alt_keywords"`
	FilePath    string    `json:"file_path" form:"file_path"` //Webcard/images  , // sitePlan/images
}
