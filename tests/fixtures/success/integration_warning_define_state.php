<?php
function stopDefine($level, $message) { throw new Exception('stop'); }

set_error_handler('stopDefine');
try { define(null, 9); } catch (Exception $e) { echo 'name|'; }
restore_error_handler();
echo defined('') ? 'created|' : 'clear|';

define('AUDIT_AFTER_THROW', 4);
echo constant('AUDIT_AFTER_THROW'), '|';

$runs = 0;
function afterWarning() { global $runs; ++$runs; return 1; }
set_error_handler('stopDefine');
try { strlen(null) + afterWarning(); } catch (Exception $e) { echo 'strlen:', $runs, '|'; }
restore_error_handler();

$runs = 0;
set_error_handler('stopDefine');
try { count([1], null) + afterWarning(); } catch (Exception $e) { echo 'count:', $runs, '|'; }
restore_error_handler();

$runs = 0;
set_error_handler('stopDefine');
try { ((int) new stdClass) + afterWarning(); } catch (Exception $e) { echo 'cast:', $runs, '|'; }
restore_error_handler();

$runs = 0;
set_error_handler('stopDefine');
try { foreach (123 as $item) { ++$runs; } } catch (Exception $e) { echo 'foreach:', $runs; }
restore_error_handler();
