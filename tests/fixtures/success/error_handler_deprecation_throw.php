<?php
function acceptsInt(int $value) { echo 'unreachable'; }
function stopOnDeprecated($level) { throw new Exception('deprecated'); }
set_error_handler('stopOnDeprecated', E_DEPRECATED);
try {
    acceptsInt(3.5);
} catch (Exception $e) {
    echo $e->getMessage(), "\n";
}
