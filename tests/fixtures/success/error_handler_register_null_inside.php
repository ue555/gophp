<?php
function original($level, $message) {
    echo $message, "|";
    echo set_error_handler(null) === null ? "null|" : "wrong|";
    return true;
}
set_error_handler('original');
trigger_error('first', E_USER_WARNING);
trigger_error('second', E_USER_WARNING);
echo "done\n";
