<?php
function secondWarning($level) { echo 'second|'; return true; }
function firstWarning($level) {
    echo 'first|';
    set_error_handler('secondWarning');
    return true;
}
set_error_handler('firstWarning');
trigger_error('one');
trigger_error('two');
restore_error_handler();
echo "done\n";
