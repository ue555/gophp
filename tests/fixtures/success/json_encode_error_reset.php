<?php
json_decode('{bad');echo json_last_error(),"\n";
echo json_encode(['ok'=>1]),':',json_last_error(),"\n";
json_decode('{bad');json_encode(['ok'=>1],JSON_THROW_ON_ERROR);echo json_last_error(),"\n";
json_encode(INF);echo json_last_error(),"\n";
json_encode('ok');echo json_last_error(),"\n";
$resource=fopen(__FILE__,'r');json_encode($resource);echo json_last_error(),"\n";fclose($resource);
