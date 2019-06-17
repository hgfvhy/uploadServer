<?php

$url = "http://127.0.0.1:9080/upload";
//$url = "http://test.com/upload.php";
//$path =realpath("157k.jpg") ;
$path ="157k.jpg" ;
$path1 ="158k.jpg" ;
//var_dump($path);
//exit();
if (class_exists("CURLFile")) {
    var_dump(1);
    $path = new CURLFile($path);
    $path1 = new CURLFile($path1);
} else {
    var_dump(2);
    $path = "@" . $path;
    $path1 = "@" . $path1;
}
$ff = json_encode([$path,$path]);
$post_data = array(
    "file" =>$path
);
$post_data1 = array(
    "file" =>$path1

);
$dat = array_merge($post_data,$post_data1);
$curl = curl_init(); // 启动一个CURL会话
curl_setopt($curl, CURLOPT_URL, $url); // 要访问的地址
curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, 0); // 对认证证书来源的检查
curl_setopt($curl, CURLOPT_SSL_VERIFYHOST, 0); // 从证书中检查SSL加密算法是否存在
curl_setopt($curl, CURLOPT_USERAGENT, $_SERVER['HTTP_USER_AGENT']); // 模拟用户使用的浏览器
curl_setopt($curl, CURLOPT_FOLLOWLOCATION, 1); // 使用自动跳转
curl_setopt($curl, CURLOPT_AUTOREFERER, 1); // 自动设置Referer
curl_setopt($curl, CURLOPT_POST, 1); // 发送一个常规的Post请求
curl_setopt($curl, CURLOPT_POSTFIELDS, $post_data); // Post提交的数据包
curl_setopt($curl, CURLOPT_POSTFIELDS, $dat); // Post提交的数据包
curl_setopt($curl, CURLOPT_TIMEOUT, 30); // 设置超时限制防止死循环
curl_setopt($curl, CURLOPT_HEADER, 0); // 显示返回的Header区域内容
curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1); // 获取的信息以文件流的形式返回
$result = curl_exec($curl); // 执行操作
if (curl_errno($curl)) {
    echo 'Errno' . curl_error($curl);//捕抓异常
}
curl_close($curl); // 关闭CURL会话
var_dump($result);
?>

<!--<!doctype html>-->
<!--<html lang="en">-->
<!--<head>-->
<!--    <meta charset="UTF-8">-->
<!--    <meta name="viewport"-->
<!--          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">-->
<!--    <meta http-equiv="X-UA-Compatible" content="ie=edge">-->
<!--    <title>Document</title>-->
<!--    <script src="https://code.jquery.com/jquery-3.1.1.min.js"></script>-->
<!--</head>-->
<!--<body>-->
<!--<form action="http://127.0.0.1:9080/upload" enctype="multipart/form-data" method="post" id="form">-->
<!--    <input type="file" name="file">-->
<!--    <input type="file" name="file">-->
<!--    <input type="button" value="submit" onclick="sub()">-->
<!--</form>-->
<!--<script type="text/javascript">-->
<!--    function sub() {-->
<!--        var formData = new FormData($("#form")[0]);-->
<!--        $.ajax({-->
<!--            url: "http://127.0.0.1:9080/upload",-->
<!--            type: "POST",-->
<!--            data: formData,-->
<!--            dataType: "json",-->
<!--            contentType: false,-->
<!--            processData: false,-->
<!--            mimeType: "multipart/form-data",-->
<!--            error: function (XMLHttpRequest, textStatus, errorThrown) {-->
<!---->
<!--            },-->
<!--            success: function (response) {-->
<!---->
<!--            }-->
<!--        });-->
<!--    }-->
<!---->
<!--</script>-->
<!--</body>-->
<!--</html>-->
