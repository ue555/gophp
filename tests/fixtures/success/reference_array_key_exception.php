<?php
$a = [];
$v = 7;
set_error_handler(function ($level, $message) { throw new Exception('stop'); });
try { $a[1.5] =& $v; } catch (Exception $e) { echo 'caught|'; }
restore_error_handler();
echo count($a), '|', $v;
