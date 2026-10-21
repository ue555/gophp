<?php
function takeUnion(int|string $value): int|string { echo 'entered'; return $value; }
set_error_handler(function ($level, $message) { throw new Exception('halt'); });
try { takeUnion(3.5); echo 'after'; } catch (Exception $e) { echo $e->getMessage(); }
echo "\n";
