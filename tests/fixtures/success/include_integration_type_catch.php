<?php
try {
    require __DIR__ . '/../support/include_integration_type.inc';
} catch (TypeError $e) {
    echo 'typed|', $e->getLine(), '|';
}
echo 'continued', "\n";
