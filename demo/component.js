async function smallUplaod(server, uploadToken, file) {
    return new Promise(async resolve=>{
        if (file.size > maxUploadDirectSize) {
            resolve("文件大小超过32M，请选择分片上传");
            return
        }
        processInfoSwitch(true);
        processUpdate(1);
        fileHash = await fileMd5(file);//计算hash
        let response = await uploadDirect(server, uploadToken, file, fileHash, file.name);
        message(response);
        if (response.code == 0 && response.content.hash == fileHash) {
            processUpdate(100);
            resolve(0);
        } else {
            resolve("上传失败！");
        }
    }

    ); 
}

async function bigFileUpload(server, uploadToken, file) {
    return new Promise(async resolve=>{
    processInfoSwitch(true);
    //初始化
    processUpdate(1);
    fileHash = await fileMd5(file);//计算hash
    chunkCount = Math.ceil(file.size / chunkSize);
    let initResponse = await uploadInit(server, uploadToken, fileHash, file.name, chunkCount);
    if (initResponse.code != 0) {
        resolve("上传初始化失败");
        return
    } else if (initResponse.content.isExist == 1) {
        resolve(1);
        return
    } else {
        existedChunkArray = initResponse.content.chunkInfoArray;
    }
    isInit = true;
    let re = await uploadAfterInit(file);
    if(re!=1){
        resolve(re);
    }else{
        resolve(re);
    }
    });
}

async function uploadAfterInit(file){
    return new Promise(async resolve=>{
        //upload chunks
        for (; currentChunkIndex < chunkCount;) {
            if(bigFileUploadSwitch==false){
                resolve("上传取消");
                return
            }
            let startIndex = currentChunkIndex * chunkSize;
            let endIndex = ((startIndex + chunkSize) >= file.size) ? file.size : startIndex + chunkSize;
            let chunkBinary = file.slice(startIndex, endIndex);
            let chunkHash = await fileMd5(chunkBinary);
            message(chunkHash);
            if (existedChunkArray.indexOf(chunkHash) < 0) {
                let chunkResponse = await uploadChunk(server, uploadToken, fileHash, chunkHash, currentChunkIndex, startIndex, endIndex, chunkBinary);
                if (chunkResponse.code != 0) {
                    resolve("块文件上传失败");
                    return
                } else if (chunkResponse.content.chunkHash != chunkHash) {
                    resolve("块文件回传hash错误，上传失败");
                    return
                }
            } else {
                message(chunkHash + "文件已存在");
            }
            currentChunkIndex++;
            let p = Math.round(currentChunkIndex / chunkCount * 99);
            processUpdate(p);
        }
        //合并
        let mergeResponse = await uploadMerge(server, uploadToken, fileHash)
        if(bigFileUploadSwitch==false){
            resolve("上传取消");
            return
        }
        message(mergeResponse);
        if (mergeResponse.code == 0) {
            processUpdate(100);
            resolve(1);
        } else {
            resolve("文件合并失败");
        }
    });
    
}

//计算目标文件MD5
async function fileMd5(file) {
    return new Promise((resolve, reject) => {
        let blobSlice = File.prototype.slice || File.prototype.mozSlice || File.prototype.webkitSlice;
        let chunkSize = 2097152; // 每次读取2MB
        let startIndex = 0;
        let endIndex = 0;
        let spark = new SparkMD5.ArrayBuffer();
        function frOnload(e) {
            spark.append(e.target.result);
            if (endIndex < file.size) {
                readChunk();
            } else {
                resolve(spark.end());
            }
        };
        function frOnerror() {
            reject("calculate md5 for binary fialed");
        };
        let fileReader = new FileReader();
        fileReader.onload = frOnload;
        fileReader.onerror = frOnerror;
        function readChunk() {
            endIndex = ((startIndex + chunkSize) >= file.size) ? file.size : startIndex + chunkSize;
            fileReader.readAsArrayBuffer(blobSlice.call(file, startIndex, endIndex));
            startIndex = startIndex + chunkSize;
        }
        readChunk();
    });
}
//进度监控
function onprogress() {

}
//速度监控
function speed() {

}
//文件直传
async function uploadDirect(server, uploadToken, fileBinary, fileHash, fileName) {
    return new Promise(resolve => {
        let url = server + "/upload";
        let form = new FormData(); // FormData
        form.append("uploadToken", uploadToken);
        form.append("fileHash", fileHash);
        form.append("fileBinary", fileBinary); // 文件对象
        form.append("fileName", fileName);
        xhr = new XMLHttpRequest();  // XMLHttpRequest
        xhr.open("post", url, true);
        xhr.onerror = function () {
            resolve("request failed");
        };
        xhr.onload = function (e) {
            let data = JSON.parse(e.target.responseText);
            resolve(data);
        }
        xhr.upload.onprogress = function (evt) {
            if (evt.lengthComputable) {
                let p = Math.round(evt.loaded / evt.total * 99);
                processUpdate(p);
            }
        }
        xhr.send(form);
    });
}
//大文件初始化
async function uploadInit(server, uploadToken, fileHash, fileName, chunkCount) {
    return new Promise(resolve => {
        let url = server + "/init";
        let form = new FormData(); // FormData
        form.append("uploadToken", uploadToken);
        form.append("fileName", fileName);
        form.append("fileHash", fileHash);
        form.append("chunkCount", chunkCount);
        xhr = new XMLHttpRequest();  // XMLHttpRequest
        xhr.open("post", url, true);
        xhr.onerror = function () {
            resolve("request failed");
        };
        xhr.onload = function (e) {
            let data = JSON.parse(e.target.responseText);
            resolve(data);
        }
        xhr.send(form);
    });
}
//上传文件块
async function uploadChunk(server, uploadToken, fileHash, chunkHash, chunkIndex, chunkStart, chunkEnd, chunkBinary) {
    return new Promise(resolve => {
        let url = server + "/chunk";
        let form = new FormData(); // FormData
        form.append("uploadToken", uploadToken);

        form.append("fileHash", fileHash);
        form.append("chunkHash", chunkHash);
        form.append("chunkIndex", chunkIndex);
        form.append("chunkStart", chunkStart);
        form.append("chunkEnd", chunkEnd);
        form.append("chunkBinary", chunkBinary);
        xhr = new XMLHttpRequest();  // XMLHttpRequest
        xhr.open("post", url, true);
        xhr.onerror = function () {
            resolve("request failed");
        };
        xhr.onload = function (e) {
            let data = JSON.parse(e.target.responseText);
            resolve(data);
        }
        xhr.send(form);
    });
}

//请求合并
async function uploadMerge(server, uploadToken, fileHash) {
    return new Promise(resolve => {
        let url = server + "/merge";
        let form = new FormData(); // FormData
        form.append("uploadToken", uploadToken);
        form.append("fileHash", fileHash);
        xhr = new XMLHttpRequest();  // XMLHttpRequest
        xhr.open("post", url, true);
        xhr.onerror = function () {
            resolve("request failed");
        };
        xhr.onload = function (e) {
            let data = JSON.parse(e.target.responseText);
            resolve(data);
        }
        xhr.send(form);
    });
}

//显示进度条
function processInfoSwitch(x) {
    if (x == true) {
        $("#progressInfo").replaceWith("<div class='progress' id='progressInfo' style='height: 2rem;margin: 0 auto;'><div class='progress-bar progress-bar-striped bg-info progress-bar-animated' id='progressbar' role='progressbar'  aria-valuenow='0' aria-valuemin='0' aria-valuemax='100'></div></div>")
    } else {
        $("#progressInfo").replaceWith("<hr id='progressInfo'>")
    }
}

//更新进度条
function processUpdate(p) {
    let process = p + '%'
    let progressbar = $("#progressbar");
    progressbar.css('width', process);
    progressbar.attr("aria-valuenow", String(p));
    progressbar.text(process);
}

//进度条状态
function processStatus(x) {
    if (x == true) {
        $('#progressbar').addClass("progress-bar-animated")
    } else {
        $('#progressbar').removeClass("progress-bar-animated")
    }
}
//各全局变量恢复初始状态
function reset() {
    isInit = false;
    bigFileUploadSwitch = true;//分片上传开关
    chunkSize = 2097152; // 分块尺寸2MB
    chunkCount = 0;//分块数量，最大支持65535
    currentChunkIndex = 0;//当前块
    existedChunkArray = new Array();//已存在块
    fileHash = '';//文件hash值
    processInfoSwitch(false);
}
//打印操作日志

function message(info){
    if(infoArray.length<5){
        //直接添加
        infoArray[infoArray.length] = info;
        //打印信息
        printInfo();
    }else{
        for(let x=0;x<4;x++){
            infoArray[x] = infoArray[x+1];
        }
        infoArray[4] = info;
        //打印信息
        printInfo();
    }
}
function printInfo(){
    for(let x=0;x<infoArray.length;x++){
        
        let infoIndex = '#info'+x;
        let info = "<div class='alert alert-info' id='info"+x+"' role='alert'>"+infoArray[x]+"</div>";
        $(infoIndex).replaceWith(info);
    }
}
