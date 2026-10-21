<?php
function seeLevel($level, $message, $file, $line) {
    echo $level, ':', $message, '|';
    return true;
}
set_error_handler('seeLevel');
echo trigger_error('n') ? 'true|' : 'bad|';
trigger_error('w', E_USER_WARNING);
trigger_error('d', E_USER_DEPRECATED);
echo "done\n";
