<?php
set_error_handler(function ($level, $message) {
    echo 'handler|';
    return false;
});
@trigger_error('quiet', E_USER_WARNING);
restore_error_handler();
@trigger_error('also quiet', E_USER_WARNING);
echo 'done';
