<?php
set_error_handler(function ($level, $message) {
    throw new Exception('stop');
});
$after = 0;
try {
    strlen(null);
    $after = 1;
} catch (Exception $e) {
    echo $e->getMessage(), '|';
}
restore_error_handler();
echo $after;
