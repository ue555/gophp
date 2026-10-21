<?php
class ModuloBox { public float $value = 5.5; }
function stopModulo($level, $message) { throw new Exception('stop'); }
$number = 5.5;
$array = [5.5];
$box = new ModuloBox;
set_error_handler('stopModulo');
try { $number %= 2; } catch (Exception $e) { echo 'caught|'; }
try { $array[0] %= 2; } catch (Exception $e) { echo 'caught|'; }
try { $box->value %= 2; } catch (Exception $e) { echo 'caught|'; }
restore_error_handler();
echo gettype($number), ':', $number, '|';
echo gettype($array[0]), ':', $array[0], '|';
echo gettype($box->value), ':', $box->value;
