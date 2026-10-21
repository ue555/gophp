<?php
echo 'before|';
try {
    include __DIR__ . '/../support/include_integration_bad.inc';
} catch (ParseError $e) {
    echo 'first|';
}
try {
    include_once __DIR__ . '/../support/include_integration_bad.inc';
} catch (ParseError $e) {
    echo 'second|';
}
echo include_once __DIR__ . '/../support/include_integration_bad.inc', '|';
echo "after\n";
