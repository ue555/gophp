<?php
$old = set_include_path(include_path: __DIR__ . '/../support/include_path/first');
echo ($old === '.:/usr/local/lib/php' ? 'old' : 'bad'), '|';
echo include 'choice.inc', '|';
ini_restore(option: 'include_path');
echo (get_include_path() === $old ? 'restored' : 'bad'), "\n";
