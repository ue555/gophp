<?php
class HandlerDynamicBox { public int $value = 1; }
function stopWarning($level, $message) { throw new Exception('handled'); }
set_error_handler('stopWarning');
$box = new HandlerDynamicBox;
try { $box->{$missing} = 7; echo 'assigned'; } catch (Exception $e) { echo $e->getMessage(); }
restore_error_handler();
echo '|', $box->value, "\n";
