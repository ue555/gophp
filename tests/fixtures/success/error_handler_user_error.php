<?php
function handleFatal($level, $message) { echo $level, ':', $message, '|'; return true; }
set_error_handler('handleFatal');
echo trigger_error('handled', E_USER_ERROR) ? "continued\n" : "failed\n";
