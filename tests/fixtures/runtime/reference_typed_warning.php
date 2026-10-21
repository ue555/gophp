<?php
class TypedWarning { public int $value = 1; }
$box = new TypedWarning;
$alias =& $box->value;
set_error_handler(function ($level, $message) { throw new Exception('stop'); });
try { $alias = 2.5; } catch (Exception $e) { echo 'caught|'; }
restore_error_handler();
echo $box->value;
