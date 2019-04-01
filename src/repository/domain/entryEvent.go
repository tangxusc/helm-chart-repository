package domain

import "mime/multipart"

type FileUploaded struct {
	File        *multipart.File
	ChartName   string
	FileName    string
	NewFileName string
}

type ChartCreated struct {
	*ChartVersion
	FileName string
}

type ChartDeleted struct {
	ChartName string
	Version   string
}

type ChartUpdated struct {
	ChartTotal int
}
