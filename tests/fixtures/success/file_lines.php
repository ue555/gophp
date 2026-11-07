<?php
// file(): lines keep their newline unless FILE_IGNORE_NEW_LINES; empty lines
// are dropped only with FILE_SKIP_EMPTY_LINES; a missing file warns and
// returns false. Symfony TypeInfo reads class files with file() to collect
// `use` statements.
$path = __DIR__ . '/../support/file_lines.txt';
$lines = file($path);
echo count($lines), '|', json_encode($lines), "\n";
echo json_encode(file($path, FILE_IGNORE_NEW_LINES)), "\n";
echo json_encode(file($path, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES)), "\n";
echo json_encode(file($path, FILE_SKIP_EMPTY_LINES)), "\n";
echo json_encode(file(__DIR__ . '/../support/empty.txt')), "\n";
echo var_export(@file(__DIR__ . '/../support/missing.txt'), true), "\n";
