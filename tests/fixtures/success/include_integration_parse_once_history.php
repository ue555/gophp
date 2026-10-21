<?php
echo 'before|';
try {
    require_once __DIR__ . '/../support/include_integration_bad.inc';
} catch (ParseError $e) {
    echo 'first|';
}
echo include_once __DIR__ . '/../support/include_integration_bad.inc', '|';
echo require_once __DIR__ . '/../support/include_integration_bad.inc', '|';
echo "after\n";
