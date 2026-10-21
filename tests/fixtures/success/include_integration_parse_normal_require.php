<?php
echo 'before|';
try {
    require __DIR__ . '/../support/include_integration_bad.inc';
} catch (ParseError $e) {
    echo 'first|';
}
try {
    require_once __DIR__ . '/../support/include_integration_bad.inc';
} catch (ParseError $e) {
    echo 'second|';
}
echo require_once __DIR__ . '/../support/include_integration_bad.inc', '|';
echo "after\n";
