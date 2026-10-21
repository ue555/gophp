<?php
$first = __DIR__ . '/../support/include_path/first';
$second = __DIR__ . '/../support/include_path/second';
$old = set_include_path($second . ':' . $first);
echo include 'choice.inc', '|';
set_include_path($first . ':' . $second);
echo include 'choice.inc', '|';
echo ($old === get_include_path() ? 'same' : 'changed'), "\n";
