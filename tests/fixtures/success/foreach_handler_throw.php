<?php
set_error_handler(function ($level, $message) {
    throw new Exception('stop');
});
$after = 0;
try {
    foreach (4 as $item) {
        $after = 1;
    }
    $after = 2;
} catch (Exception $e) {
    echo $e->getMessage(), '|';
}
restore_error_handler();
echo $after;
