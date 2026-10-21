<?php
class IntegratedBox { public int $number = 1; }
function integratedStop($level, $message) { throw new Exception('stopped'); }
$box = new IntegratedBox;
set_error_handler('integratedStop');
try { $box->extra = 7; } catch (Exception $e) { echo isset($box->extra) ? 'set|' : 'absent|'; }
try { $box->number = 3.5; } catch (Exception $e) { echo $box->number, '|'; }
restore_error_handler();
try { $box->number = 'bad'; } catch (TypeError $e) { echo $box->number, '|'; }
echo 'done';
