<?php
function throwingHandler($level, $message) { throw new Exception($message); }
set_error_handler('throwingHandler');
try {
    echo $missing, 'unreachable';
} catch (Exception $e) {
    echo $e->getMessage(), '|';
}
restore_error_handler();
echo "done\n";
