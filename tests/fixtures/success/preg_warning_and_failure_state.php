<?php
set_error_handler(function($n,$s){echo $s,"\n";return true;});
echo json_encode(preg_match('/(/','x')),':',preg_last_error(),"\n";
try {preg_match([], 'x');}catch(Throwable $e){echo $e->getMessage(),"\n";}
echo preg_last_error(),"\n";
echo json_encode(preg_match('/a/q','x')),':',preg_last_error(),"\n";
echo json_encode(preg_match('/a/','x')),':',preg_last_error(),"\n";
echo json_encode(preg_replace('/x/u','z',"\xff")),':',preg_last_error(),"\n";
echo json_encode(preg_split('//u',"\xff")),':',preg_last_error(),"\n";
echo preg_match('/(*LIMIT_MATCH=1)(a+)+b/',str_repeat('a',30).'c'),':',preg_last_error(),"\n";
echo preg_match('/x/','x'),':',preg_last_error(),"\n";
