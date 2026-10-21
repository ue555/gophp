<?php
$visits = 0;
$first = __DIR__ . '/../support/include_path/first';
set_include_path($first);
echo include_once 'once.inc', '|';
echo include_once $first . '/./once.inc', '|';
echo include_once $first . '/link.inc', '|';
set_include_path(__DIR__ . '/../support/include_path/second:' . $first);
echo require_once 'once.inc', '|', $visits, "\n";
