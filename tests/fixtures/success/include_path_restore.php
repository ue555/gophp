<?php
$initial = get_include_path();
$first = __DIR__ . '/../support/include_path/first';
echo (set_include_path($first) === $initial ? 'old' : 'wrong'), '|';
echo (get_include_path() === $first ? 'set' : 'wrong'), '|';
set_include_path('');
ini_restore('include_path');
echo (get_include_path() === $initial ? 'restored' : 'wrong'), "\n";
