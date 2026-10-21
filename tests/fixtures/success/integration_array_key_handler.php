<?php
function integratedKeyStop($level, $message) { throw new Exception('key'); }
$values = [];
set_error_handler('integratedKeyStop');
try { $values[3.5] = 7; } catch (Exception $e) { echo $e->getMessage(), '|'; }
restore_error_handler();
echo count($values), '|';
$values[3] = 8;
echo $values[3];
