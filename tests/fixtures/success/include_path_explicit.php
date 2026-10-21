<?php
set_include_path(__DIR__ . '/../support/include_path/second');
echo include __DIR__ . '/../support/include_path/first/choice.inc', '|';
echo include './tests/fixtures/support/include_path/first/choice.inc', '|';
echo include '../' . basename(dirname(__DIR__, 3)) . '/tests/fixtures/support/include_path/first/choice.inc', "\n";
