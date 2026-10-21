<?php
echo 'before|';
try {
    include_once __DIR__ . '/../support/include_integration_bad.inc';
} catch (ParseError $e) {
    echo 'parse|';
}
echo include_once __DIR__ . '/../support/include_integration_bad.inc', '|';
echo include __DIR__ . '/../support/include_integration_return.inc', '|';
echo 'after', "\n";
