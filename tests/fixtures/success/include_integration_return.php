<?php
function readChild() {
    echo 'caller|';
    $value = include __DIR__ . '/../support/include_integration_return.inc';
    echo 'resumed|';
    return $value + 1;
}
echo readChild(), "\n";
