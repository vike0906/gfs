package upload

const (
	errMessage      = "need param %s, but not found"
	fileDownloadUrl = "http://host:port/download/"
	fileCorrupted   = "file corrupted, please upload again"
)

const (
	paramUploadToken    = "uploadToken"
	paramFileBinary     = "fileBinary"
	paramFileName       = "fileName"
	paramFileHash       = "fileHash"
	paramChunkCount     = "chunkCount"
	paramChunkInfoArray = "chunkInfoArray"
	paramParentFileHash = "parentFileHash"
	paramChunkIndex     = "chunkIndex"
	paramChunkStart     = "chunkStart"
	paramChunkEnd       = "chunkEnd"
)

const (
	exist  = 1
	unExit = 2
)
