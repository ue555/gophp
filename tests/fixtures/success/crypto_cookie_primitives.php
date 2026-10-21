<?php
echo hash_hmac('sha256','message','secret'),"\n";
echo bin2hex(hash_hmac('sha1','message','secret',true)),"\n";
echo json_encode([hash_equals('abc','abc'),hash_equals('abc','abd'),hash_equals('a','aa')]),"\n";
try{hash_equals(1,'1');}catch(Throwable $e){echo $e->getMessage(),"\n";}
try{hash_hmac('xxh3','a','b');}catch(Throwable $e){echo $e->getMessage(),"\n";}
foreach(['aes-128-cbc','aes-256-cbc','aes-128-gcm','aes-256-gcm'] as $cipher){
 $iv=str_repeat('i',openssl_cipher_iv_length($cipher));$key=str_repeat('k',32);$tag='unchanged';
 $encrypted=openssl_encrypt('日本 cookie',$cipher,$key,0,$iv,$tag,'aad');
 echo $cipher,':',$encrypted,':',bin2hex($tag??''),"\n";
 echo openssl_decrypt($encrypted,$cipher,$key,0,$iv,$tag??'','aad'),"\n";
}
$iv=str_repeat('i',16);
echo json_encode(openssl_decrypt('bad','aes-256-cbc','key',0,$iv)),"\n";
echo json_encode(openssl_encrypt('short','aes-128-cbc','key',OPENSSL_ZERO_PADDING,$iv)),"\n";
$tag='old';$i=0;$tags=[];
$v=openssl_encrypt('x','aes-128-gcm','key',OPENSSL_RAW_DATA,str_repeat('i',12),$tags[$i++]);
echo $i,':',strlen($tags[0]),':',openssl_decrypt($v,'aes-128-gcm','key',OPENSSL_RAW_DATA,str_repeat('i',12),$tags[0]),"\n";
