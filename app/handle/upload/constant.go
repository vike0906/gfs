package upload

const (
	errMessage       = "need param %s, but not found"
	fileDownloadUrl  = "http://host:port/download/"
	fileCorrupted    = "file corrupted, please upload again"
	chunkCountError  = "count of chunks is error, upload failed"
	fileHashNotExist = "file hash should exist, but not"
)

const (
	paramUploadToken    = "uploadToken"
	paramFileBinary     = "fileBinary"
	paramFileName       = "fileName"
	paramFileHash       = "fileHash"
	paramChunkCount     = "chunkCount"
	paramChunkInfoArray = "chunkInfoArray"
	//paramParentFileHash = "parentFileHash"
	paramChunkHash   = "chunkHash"
	paramChunkIndex  = "chunkIndex"
	paramChunkStart  = "chunkStart"
	paramChunkEnd    = "chunkEnd"
	paramChunkBinary = "chunkBinary"
)

const (
	exist  = 1
	unExit = 2
)
