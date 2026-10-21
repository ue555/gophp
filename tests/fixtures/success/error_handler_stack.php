<?php
function firstHandler($level, $message) { echo "first:$level:$message|"; return true; }
function secondHandler($level, $message) { echo "second:$level:$message|"; return true; }
echo set_error_handler('firstHandler') === null ? 'null|' : 'bad|';
echo set_error_handler('secondHandler') === 'firstHandler' ? 'first|' : 'bad|';
trigger_error('one', E_USER_NOTICE);
restore_error_handler();
trigger_error('two', E_USER_WARNING);
set_error_handler(null);
restore_error_handler();
trigger_error('three', E_USER_NOTICE);
restore_error_handler();
echo "done\n";
