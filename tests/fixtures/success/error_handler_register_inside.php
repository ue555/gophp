<?php
function replacement($level, $message) {
    echo "new:", $message, "|";
    return true;
}
function original($level, $message) {
    echo set_error_handler('replacement') === null ? "null|" : "wrong|";
    trigger_error('nested', E_USER_WARNING);
    return true;
}
set_error_handler('original');
trigger_error('outer', E_USER_WARNING);
trigger_error('later', E_USER_WARNING);
echo "done\n";
