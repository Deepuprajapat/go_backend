package request

type UploadFileRequest struct {
	FileName    string `json:"file_name" form:"file_name"`
	AltKeywords string `json:"alt_keywords" form:"alt_keywords"`
	FilePath    string `json:"file_path" form:"file_path"` //Webcard/images  , // sitePlan/images
}
