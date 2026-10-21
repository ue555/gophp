<?php
set_error_handler(function ($level, $message) { return true; });
set_include_path('');
$missing = include 'choice.inc';
echo (get_include_path() === '' ? 'empty' : 'bad'), '|', ($missing ? 'found' : 'missing'), '|';
set_include_path(':' . __DIR__ . '/../support/include_path/first');
echo include 'choice.inc', '|';
ini_restore('include_path');
echo (get_include_path() === '.:/usr/local/lib/php' ? 'restored' : 'bad'), "\n";
restore_error_handler();
