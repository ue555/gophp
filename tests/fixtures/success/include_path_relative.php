<?php
$old = set_include_path('tests/fixtures/support/include_path/first:tests/fixtures/support/include_path/second');
echo include 'choice.inc', '|';
echo include 'nested.inc', '|';
echo ($old === '.:/usr/local/lib/php' ? 'default' : 'other'), "\n";
