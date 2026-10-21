<?php
function handleMissing($level, $message, $file, $line) {
    echo $level, '|', $file === __FILE__ ? 'same|' : 'wrong|', $line, '|';
    return true;
}
set_error_handler('handleMissing');
echo include __DIR__ . '/support/error_handler_absent.inc' ? 'loaded|' : 'missing|';
try {
    require __DIR__ . '/support/error_handler_absent.inc';
} catch (Error $e) {
    echo 'caught|';
}
echo "done\n";
