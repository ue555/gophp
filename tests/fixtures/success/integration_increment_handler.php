<?php
function integratedIncrementStop($level, $message) { throw new Exception('increment'); }
$value = '!';
set_error_handler('integratedIncrementStop');
try { $value++; } catch (Exception $e) { echo $e->getMessage(), '|'; }
restore_error_handler();
echo $value;
